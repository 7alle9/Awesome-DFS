package main

import (
	fp "Awesome-DFS/client/file_partition"
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

func tempMain() {
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

type chunkDescriptor struct {
	uniqueName      string
	offset          int64
	size            int
	storageLocation []string
}

func getChunkPartition(file *os.File, fileSize int64, chunkSize int64) []chunkDescriptor {
	nbChunks := fileSize / chunkSize
	if (fileSize % chunkSize) != 0 {
		nbChunks++
	}

	partition := make([]chunkDescriptor, nbChunks)

	for i := int64(0); i < fileSize; i += chunkSize {
		size := chunkSize
		if i+chunkSize > fileSize {
			size = fileSize - i
		}
		uniqueName := fmt.Sprintf("chunk_%d", i/chunkSize)
		partition[i] = chunkDescriptor{
			uniqueName,
			i,
			int(size),
			[]string{"localhost:8080", "localhost:8081", "localhost:8082"},
		}
	}

	return partition
}

func makeJobs(jobsTo map[string]chan job, partition []chunkDescriptor, file *os.File) {
	for _, chunkDesc := range partition {
		worker := chunkDesc.storageLocation[0]
		if _, ok := jobsTo[worker]; !ok {
			jobsTo[worker] = make(chan job)
		}
		jobsTo[worker] <- job{file, chunkDesc.offset, chunkDesc.size, chunkDesc.uniqueName}
	}
}

func main() {
	f, err := os.Open("test.txt")
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	part, err := fp.GetFilePartition(f, 4, 3)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	for _, chunk := range part.Chunks {
		fmt.Println("----------------")
		fmt.Printf("Chunk: %s\n", chunk.Name)
		fmt.Printf("Size: %d\n", chunk.Size)
		fmt.Printf("Offset: %v\n", chunk.Offset)
		fmt.Printf("Send To: %v\n", chunk.SendTo)
		fmt.Printf("Replicas: %v\n", chunk.ReplicaChain)
	}
}
