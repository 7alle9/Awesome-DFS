package main

import (
	pb "Awesome-DFS/storage"
	"context"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"os"
	"sync"
	"time"
)

var (
	serverAddr = flag.String("addr", "localhost:8080", "The server address in the format of host:port")
	ramUsage   = flag.Int64("ram", 512, "The amount of RAM to use in MB")
	filepath   = flag.String("file", "", "The path of the file to store")
)

type job struct {
	file       *os.File
	offset     int64
	size       int
	uniqueName string
}

func worker(jobs <-chan job, results chan<- *pb.StoreResponse, client pb.StorageClient) {
	for job := range jobs {
		chunk, err := readChunk(job.file, job.offset, job.size)
		if err != nil {
			fmt.Println(err)
			return
		}

		//log.Printf("Sending chunk %s\n", job.uniqueName)
		response, err := client.Store(
			context.Background(),
			&pb.Chunk{UniqueName: job.uniqueName, Data: chunk},
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		results <- response
	}
}

func readChunk(file *os.File, offset int64, size int) ([]byte, error) {
	chunk := make([]byte, size)
	_, err := file.ReadAt(chunk, offset)
	if err != nil {
		return nil, err
	}
	return chunk, nil
}

func main() {
	flag.Parse()
	if *filepath == "" {
		panic("Please provide the path of the file to store")
	}

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(*serverAddr, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
	}
	defer conn.Close()
	log.Printf("Connected to %s\n", *serverAddr)
	client := pb.NewStorageClient(conn)

	file, err := os.Open(*filepath)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	fileSize := info.Size()

	chunkSize := int64(64 * 1024)

	nbJobs := fileSize / chunkSize
	if (fileSize % chunkSize) != 0 {
		nbJobs++
	}
	nbWorkers := int64(50)

	jobs := make(chan job, nbJobs)
	results := make(chan *pb.StoreResponse, nbJobs)

	for i := int64(0); i < fileSize; i += chunkSize {
		size := chunkSize
		if i+chunkSize > fileSize {
			size = fileSize - i
		}
		uniqueName := fmt.Sprintf("chunk_%d", i/chunkSize)
		jobs <- job{file, i, int(size), uniqueName}
	}
	close(jobs)

	start := time.Now()

	var wg sync.WaitGroup
	for i := int64(0); i < nbWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(jobs, results, client)
		}()
	}
	wg.Wait()
	close(results)

	elapsed := time.Since(start)
	log.Printf("Elapsed time: %s\n", elapsed)

	ok := 0
	bad := 0
	for result := range results {
		if result.Status == pb.Status_STATUS_OK {
			ok++
		} else {
			log.Fatalf("Error: %s", result.Message)
			bad++
		}
	}
	log.Printf("OK: %d | bad: %d", ok, bad)
}
