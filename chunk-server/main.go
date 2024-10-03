package main

import (
	pb "Awesome-DFS/storage"
	"context"
	"crypto/sha256"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type ChunkServer struct {
	pb.UnimplementedStorageServer
}

func (s *ChunkServer) store(_ context.Context, in *pb.Chunk) (*pb.StoreResponse, error) {
	err := createChunk(in)
	if err != nil {
		errorMessage := fmt.Sprintf("Error: %v", err)
		return &pb.StoreResponse{Status: pb.Status_STATUS_ERROR, Message: errorMessage}, nil
	}
	return &pb.StoreResponse{Status: pb.Status_STATUS_OK}, nil
}

func createChunk(chunk *pb.Chunk) error {
	checksum, err := getChecksum(chunk.Data)
	if err != nil {
		return err
	}
	checksumPath := fmt.Sprintf("metedata/%s.checksum", chunk.UniqueName)
	err = writeData(checksumPath, checksum)
	if err != nil {
		return err
	}

	dataPath := fmt.Sprintf("chunks/%s", chunk.UniqueName)
	err = writeData(dataPath, chunk.Data)
	if err != nil {
		return err
	}

	return nil
}

func getChecksum(data []byte) ([]byte, error) {
	hasher := sha256.New()
	_, err := hasher.Write(data)
	if err != nil {
		return nil, err
	}
	return hasher.Sum(nil), nil
}

func writeData(filepath string, data []byte) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	err = file.Close()
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
	pb.RegisterStorageServer(grpcServer, new(ChunkServer))
	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
