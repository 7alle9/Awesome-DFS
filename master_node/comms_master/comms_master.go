package comms_master

import (
	pb "Awesome-DFS/servers_comms"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

type availability struct {
	Node         string
	Status       pb.Status
	ResponseTime time.Duration
}

func connect(address string) (pb.CommsClient, *grpc.ClientConn, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(address, opts...)
	if err != nil {
		return nil, nil, nil
	}
	return pb.NewCommsClient(conn), conn, nil
}

func ping(address string, payload []byte) (pb.Status, time.Duration) {
	storageNode, conn, err := connect(address)
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

func pingWorker(address string, size int, nodeStatus chan<- *availability) {
	payload := make([]byte, size)
	status, responseTime := ping(address, payload)

	if status == pb.Status_STATUS_READY {
		log.Printf("Node %s is READY. Response time %v\n", address, responseTime)
	} else {
		log.Printf("Node %s is NOT READY. Response time %v\n", address, responseTime)
	}

	nodeStatus <- &availability{address, status, responseTime}
}

func GetAvailableNodes(addressBook []string, chunkSize int) map[string]time.Duration {
	nodeStatus := make(chan *availability, len(addressBook))
	for _, nodeAddr := range addressBook {
		go pingWorker(nodeAddr, chunkSize, nodeStatus)
	}

	availableNodes := make(map[string]time.Duration)
	for i := 0; i < len(addressBook); i++ {
		nodeAvailability := <-nodeStatus
		if nodeAvailability.Status == pb.Status_STATUS_READY {
			availableNodes[nodeAvailability.Node] = nodeAvailability.ResponseTime
		}
	}

	return availableNodes
}
