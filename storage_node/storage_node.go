package main

import (
	comms "Awesome-DFS/storage_node/comms_storage"
	store "Awesome-DFS/storage_node/file_storage"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"net"
)

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
	store.RegisterFileTransferServer(grpcServer)

	log.Printf("Listening on %s\n", lis.Addr().String())

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
