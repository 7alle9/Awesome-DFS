package server_connection

import (
	pb "Awesome-DFS/partition"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
)

var (
	address                  = "localhost:8079"
	conn    *grpc.ClientConn = nil
	opts                     = []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
)

func makeConnection() error {
	var err error

	conn, err = grpc.NewClient(address, opts...)
	if err != nil {
		return err
	}

	return nil
}

func getMasterConnection() *grpc.ClientConn {
	if conn == nil {
		log.Printf("Creating new connection to %s\n", address)
		err := makeConnection()
		if err != nil {
			panic(err)
		}
	}

	return conn
}

func GetPartitionClient() pb.PartitionClient {
	serverConn := getMasterConnection()
	return pb.NewPartitionClient(serverConn)
}
