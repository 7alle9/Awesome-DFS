package main

import (
	"Awesome-DFS/master-node/comms"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	addressBook := []string{
		"192.168.100.7:8080",
		"localhost:8080",
		"192.168.100.6",
	}
	comms.GetAvailableNodes(addressBook, 1024, opts)
}
