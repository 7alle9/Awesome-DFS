package main

import (
	pb "Awesome-DFS/dispatcher"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type StorageNode struct {
	pb.UnimplementedDispatcherServer
}

func (s *StorageNode) Dispatch(_ context.Context, in *pb.Chunk) (*pb.Response, error) {
	err := createChunk(in)
	if err != nil {
		return &pb.Response{Status: pb.Status_STATUS_ERROR, Message: "Error"}, err
	}
	return &pb.Response{Status: pb.Status_STATUS_OK}, nil
}

func createChunk(chunk *pb.Chunk) error {
	file, err := os.OpenFile(
		fmt.Sprintf("chunks/%s", chunk.UniqueName),
		os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(chunk.Data)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	lis, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}
	grpcServer := grpc.NewServer()
	pb.RegisterDispatcherServer(grpcServer, new(StorageNode))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
