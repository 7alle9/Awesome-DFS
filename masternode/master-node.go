package main

import (
	comms "Awesome-DFS/serverscomms"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

func connect(address string, opts []grpc.DialOption) (comms.CommsClient, error) {
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, nil
	}
	return comms.NewCommsClient(conn), nil
}

func ping(address string, payload []byte, opts []grpc.DialOption) (comms.Status, time.Duration, error) {
	storageNode, err := connect(address, opts)
	if err != nil {
		return comms.Status_STATUS_NOT_READY, 0, nil
	}

	start := time.Now()
	pingPayload := &comms.PingPayload{Payload: payload}
	response, err := storageNode.Ping(context.Background(), pingPayload)
	if err != nil {
		return comms.Status_STATUS_NOT_READY, 0, nil
	}
	elapsed := time.Since(start)

	return response.Status, elapsed, nil
}

func main() {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	addressBook := []string{
		"192.168.100.7:8080",
		"localhost:8080",
		"192.168.100.6",
	}
}
