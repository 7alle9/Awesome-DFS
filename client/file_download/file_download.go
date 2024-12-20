package file_download

import (
	pool "Awesome-DFS/client/connection_pool_manager"
	fp "Awesome-DFS/client/file_partition"
	hs "Awesome-DFS/client/hashing_service"
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
	dehasher *hs.HashingService
)

func handleJobError(errMessage string, connId int) {
	log.Println(errMessage)
	log.Println("trying next storage node")
	pool.ReleaseConn(connId)
}

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
				msg := fmt.Sprintf("could not reach storage node %s: %v", storageNode, err)
				handleJobError(msg, connId)
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
					msg := fmt.Sprintf("error while downloading chunk: %v", err)
					handleJobError(msg, connId)
					chunkIsGood = false
					break
				}

				switch payload := chunk.Payload.(type) {
				case *down.Chunk_Data:
					data, err := dehasher.DecryptByteArray(payload.Data.RawBytes)
					if err != nil {
						msg := fmt.Sprintf("error while decrypting data: %v", err)
						handleJobError(msg, connId)
						chunkIsGood = false
						break
					}

					err = writeData(offset, data)
					if err != nil {
						msg := fmt.Sprintf("error while writing data: %v", err)
						handleJobError(msg, connId)
						chunkIsGood = false
						break
					}

					hasher.Write(payload.Data.RawBytes)

					offset += int64(len(payload.Data.RawBytes))
				case *down.Chunk_IntegrityCheck:
					checksum := fmt.Sprintf("%x", hasher.Sum(nil))
					if checksum != payload.IntegrityCheck.Checksum {
						msg := fmt.Sprintf("checksum mismatch: chunk is corrupted")
						handleJobError(msg, connId)
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

	dehasher, err = hs.GetHasher()
	if err != nil {
		return err
	}

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
