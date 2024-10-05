package main

import (
	"Awesome-DFS/master-node/comms"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	addressBook := []string{
		"localhost:8080",
		"localhost:8081",
		"localhost:8082",
		"localhost:8083",
	}
	availableNodes := comms.GetAvailableNodes(addressBook, 1024, opts)
	fmt.Printf("%v", availableNodes)
}
