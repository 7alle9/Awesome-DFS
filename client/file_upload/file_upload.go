package file_upload

import (
	pool "Awesome-DFS/client/connection_pool_manager"
	fp "Awesome-DFS/client/file_partition"
	part "Awesome-DFS/partition"
	up "Awesome-DFS/transfer"
	"context"
	"log"
	"os"
	"sync"
)

var (
	payloadSize int64 = 2 * 1024 * 1024
	file        *os.File
	fileUuid    string
)

func readData(offset int64, data []byte) {
	_, err := file.ReadAt(data, offset)
	if err != nil {
		panic(err)
	}
}

func worker(jobs <-chan *part.Chunk) {
	for info := range jobs {
		conn, connId := pool.ConnectTo(info.SendTo)
		client := up.NewFileTransferClient(conn)
		stream, err := client.Upload(context.Background())
		if err != nil {
			log.Fatalf("error while opening stream: %v", err)
		}

		metadata := &up.MetaData{
			FileUuid:     fileUuid,
			UniqueName:   info.Name,
			Size:         info.Size,
			ReplicaChain: info.ReplicaChain,
		}
		chunkMeta := &up.Chunk_Meta{Meta: metadata}
		chunk := &up.Chunk{Payload: chunkMeta}

		err = stream.Send(chunk)
		if err != nil {
			log.Fatalf("error sending metadata: %v", err)
		}

		data := make([]byte, payloadSize)
		offset := info.Offset
		limit := info.Offset + info.Size
		for i := offset; i < limit; i += payloadSize {
			if i+payloadSize > limit {
				data = data[:limit-i]
			}

			readData(i, data)

			payloadData := &up.Data{RawBytes: data, Number: i / payloadSize}
			chunkData := &up.Chunk_Data{Data: payloadData}
			chunk.Payload = chunkData

			err = stream.Send(chunk)
			if err != nil {
				log.Fatalf("error sending data: %v", err)
			}
		}

		reply, err := stream.CloseAndRecv()
		if err != nil {
			log.Fatalf("error while closing stream: %v", err)
		}
		if reply.Status == up.Status_STATUS_OK {
			log.Printf("%s: %s", info.Name, reply.Message)
		} else {
			log.Fatalf("chunk %s failed to upload: %s", info.Name, reply.Message)
		}

		pool.ReleaseConn(connId)
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

func UploadFile(fileToUp *os.File, chunkSize int64, nbReplicas int) error {
	file = fileToUp

	filePartition, err := fp.GetFilePartition(file, chunkSize, nbReplicas)
	if err != nil {
		return err
	}

	fileUuid = filePartition.FileUuid

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

	log.Printf("File %s uploaded successfully", file.Name())

	return nil
}
