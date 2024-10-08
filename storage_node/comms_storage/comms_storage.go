package comms_storage

import (
	pb "Awesome-DFS/servers_comms"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/peer"
	"log"
)

type commsNode struct {
	pb.UnimplementedCommsServer
}

func (s *commsNode) Ping(ctx context.Context, in *pb.PingPayload) (*pb.PingResponse, error) {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return nil, fmt.Errorf("failed to get peer from context")
	}

	log.Printf("Received ping from %s. Payload Size %s\n", p.Addr.String(), payloadSizeString(len(in.Payload)))

	return &pb.PingResponse{Status: pb.Status_STATUS_READY}, nil
}

func RegisterCommsServer(server *grpc.Server) {
	pb.RegisterCommsServer(server, new(commsNode))
}

func payloadSizeString(size int) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	}
	if size < 1024*1024 {
		return fmt.Sprintf("%dKB", size/1024)
	}
	return fmt.Sprintf("%dMB", size/1024/1024)
}
