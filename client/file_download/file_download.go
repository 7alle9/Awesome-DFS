package file_download

import (
	pool "Awesome-DFS/client/connection_pool_manager"
	fp "Awesome-DFS/client/file_partition"
	down "Awesome-DFS/protobuf/download"
	part "Awesome-DFS/protobuf/partition"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var (
	file     *os.File
	fileUuid string
)

func writeData(offset int64, data []byte) error {
	_, err := file.WriteAt(data, offset)
	if err != nil {
		return err
	}

	return nil
}

func worker(jobs <-chan *part.Chunk) {
	for info := range jobs {
		chunkDesc := &down.ChunkDesc{
			FileName:  fileUuid,
			ChunkName: info.Name,
			Size:      info.Size,
		}

		info.ReplicaChain = append(info.ReplicaChain, info.SendTo)

		var chunkIsGood bool

		for _, storageNode := range info.ReplicaChain {
			conn, connId := pool.ConnectTo(storageNode)
			client := down.NewDownloadClient(conn)
			stream, err := client.Download(context.Background(), chunkDesc)

			if err != nil {
				log.Printf("could not reach storage node %s: %v", storageNode, err)
				log.Printf("trying next storage node")
				pool.ReleaseConn(connId)
				continue
			}

			chunkIsGood = true

			offset := info.Offset

			hasher := sha256.New()

			for {
				chunk, err := stream.Recv()

				if err == io.EOF {
					break
				}

				if err != nil {
					log.Printf("error while downloading chunk: %v", err)
					log.Printf("trying next storage node")
					pool.ReleaseConn(connId)
					chunkIsGood = false
					break
				}

				switch payload := chunk.Payload.(type) {
				case *down.Chunk_Data:
					err = writeData(offset, payload.Data.RawBytes)
					if err != nil {
						log.Printf("error while writing data: %v", err)
						log.Printf("trying next storage node")
						pool.ReleaseConn(connId)
						chunkIsGood = false
						break
					}

					hasher.Write(payload.Data.RawBytes)

					offset += int64(len(payload.Data.RawBytes))
				case *down.Chunk_IntegrityCheck:
					checksum := fmt.Sprintf("%x", hasher.Sum(nil))
					if checksum != payload.IntegrityCheck.Checksum {
						log.Printf("checksum mismatch: chunk is corrupted")
						log.Printf("trying next storage node")
						pool.ReleaseConn(connId)
						chunkIsGood = false
						break
					}
				}
			}

			if chunkIsGood {
				pool.ReleaseConn(connId)
				break
			}
		}

		if !chunkIsGood {
			log.Fatalf("chunk %s failed to download", info.Name)
		}
	}
}

func initJobs(partition *part.FilePartition) chan *part.Chunk {
	jobs := make(chan *part.Chunk, len(partition.Chunks))
	for _, chunk := range partition.Chunks {
		jobs <- chunk
	}
	close(jobs)

	return jobs
}

func Download(fileName string) error {
	log.Printf("Request to download file %s", fileName)

	log.Printf("Reconstructing file %s", fileName)
	filePartition, err := fp.ReconstructFile(fileName)
	if err != nil {
		log.Printf("failed to get partition: %v", err)
		return fmt.Errorf("failed to get partition: %v", err)
	}

	fileUuid = filePartition.FileUuid

	file, err = os.Create(fileName)
	if err != nil {
		log.Printf("failed to create file: %v", err)
		return fmt.Errorf("failed to create file: %v", err)
	}

	jobs := initJobs(filePartition)
	nbWorkers := min(len(filePartition.Chunks), 500)

	var wg sync.WaitGroup
	for i := 0; i < nbWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(jobs)
		}()
	}
	wg.Wait()

	log.Printf("File %s Downloaded successfully", file.Name())
	return nil
}
