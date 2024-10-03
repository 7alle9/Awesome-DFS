package main

import (
	pb "Awesome-DFS/dispatcher"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

func main() {
	conn, err := grpc.NewClient("192.168.100.6:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	masterNode := pb.NewDispatcherClient(conn)
	response, err := masterNode.Dispatch(context.Background(), &pb.Chunk{
		UniqueName: "chunk_1",
		Data:       []byte("Hello, World!"),
	})
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
	log.Printf("Response: %v", response.Status)
}
