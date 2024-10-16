package master_connection

import (
	part "Awesome-DFS/protobuf/partition"
	val "Awesome-DFS/protobuf/validation"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"sync"
)

var (
	mu      sync.Mutex
	address                  = "localhost:8079"
	conn    *grpc.ClientConn = nil
	opts                     = []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		//grpc.WithKeepaliveParams(keepalive.ClientParameters{Time: 10, Timeout: 20, PermitWithoutStream: true}),
	}
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
	mu.Lock()
	defer mu.Unlock()

	if conn == nil {
		log.Printf("Creating new connection to %s\n", address)
		err := makeConnection()
		if err != nil {
			panic(err)
		}
	}

	return conn
}

func GetPartitionClient() part.PartitionClient {
	serverConn := getMasterConnection()

	log.Printf("Creating new partition client\n")
	return part.NewPartitionClient(serverConn)
}

func GetValidationClient() val.ValidationClient {
	serverConn := getMasterConnection()

	log.Printf("Creating new validation client\n")
	return val.NewValidationClient(serverConn)
}
