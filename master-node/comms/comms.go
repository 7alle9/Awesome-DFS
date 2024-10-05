package comms

import (
	pb "Awesome-DFS/serverscomms"
	"context"
	"google.golang.org/grpc"
	"log"
	"time"
)

type availability struct {
	Node         string
	Status       pb.Status
	ResponseTime time.Duration
}

func connect(address string, opts []grpc.DialOption) (pb.CommsClient, *grpc.ClientConn, error) {
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, nil, nil
	}
	return pb.NewCommsClient(conn), conn, nil
}

func ping(address string, payload []byte, opts []grpc.DialOption) (pb.Status, time.Duration) {
	storageNode, conn, err := connect(address, opts)
	if err != nil {
		return pb.Status_STATUS_NOT_READY, 0
	}
	defer conn.Close()

	start := time.Now()
	pingPayload := &pb.PingPayload{Payload: payload}
	response, err := storageNode.Ping(context.Background(), pingPayload)
	if err != nil {
		return pb.Status_STATUS_NOT_READY, 0
	}
	elapsed := time.Since(start)

	return response.Status, elapsed
}

func pingWorker(address string, size int, opts []grpc.DialOption, nodeStatus chan<- *availability) {
	payload := make([]byte, size)
	status, responseTime := ping(address, payload, opts)

	if status == pb.Status_STATUS_READY {
		log.Printf("Node %s : READY. ", address)
	} else {
		log.Printf("Node %s : NOT READY. ", address)
	}
	log.Printf("Response time: %v\n", responseTime)

	nodeStatus <- &availability{address, status, responseTime}
}

func GetAvailableNodes(addressBook []string, chunkSize int, opts []grpc.DialOption) []string {
	nodeStatus := make(chan *availability, len(addressBook))
	for _, nodeAddr := range addressBook {
		go pingWorker(nodeAddr, chunkSize, opts, nodeStatus)
	}

	availableNodes := []string{}
	for i := 0; i < len(addressBook); i++ {
		availability := <-nodeStatus
		if availability.Status == pb.Status_STATUS_READY {
			availableNodes = append(availableNodes, availability.Node)
		}
	}

	return availableNodes
}
