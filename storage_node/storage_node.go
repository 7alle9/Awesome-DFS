package main

import (
	comms "Awesome-DFS/storage_node/comms_storage"
	pb "Awesome-DFS/transfer"
	"context"
	"crypto/sha256"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

type StorageNode struct {
	pb.UnimplementedStorageServer
}

func (s *StorageNode) Store(ctx context.Context, in *pb.Chunk) (*pb.StoreResponse, error) {
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

	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	MetadataDir := fmt.Sprintf("%s\\metdata", wd)
	ChunkDir := fmt.Sprintf("%s\\chunks", wd)

	checksumPath := fmt.Sprintf("%s\\%s.checksum", MetadataDir, chunk.UniqueName)
	err = writeData(checksumPath, checksum)
	if err != nil {
		return err
	}

	dataPath := fmt.Sprintf("%s\\%s", ChunkDir, chunk.UniqueName)
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

// port flag
var port = flag.String("port", "8080", "Port to listen on")

func main() {
	flag.Parse()
	addr := fmt.Sprintf(":%s", *port)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	grpcServer := grpc.NewServer()

	comms.RegisterCommsServer(grpcServer)
	pb.RegisterStorageServer(grpcServer, new(StorageNode))

	log.Printf("Listening on %s\n", lis.Addr().String())

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
