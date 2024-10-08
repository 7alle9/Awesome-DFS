package main

import (
	ps "Awesome-DFS/master_node/partiton_server"
	"Awesome-DFS/storage_node/comms_storage"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	lis, err := net.Listen("tcp", ":8079")
	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	grpcServer := grpc.NewServer()
	ps.RegisterPartitionServer(grpcServer)

	comms_storage.RegisterCommsServer(grpcServer)

	log.Printf("Listening on %s\n", lis.Addr().String())

	err = grpcServer.Serve(lis)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}
}
