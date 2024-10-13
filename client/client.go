package main

import (
	fp "Awesome-DFS/client/file_partition"
	"fmt"
	"log"
	"os"
)

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
