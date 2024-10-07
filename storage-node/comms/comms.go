package comms

import (
	pb "Awesome-DFS/servers-comms"
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

	log.Printf("Received ping from %s. Payload Size %dKB\n", p.Addr.String(), len(in.Payload)/1024)

	return &pb.PingResponse{Status: pb.Status_STATUS_READY}, nil
}

func RegisterCommsServer(server *grpc.Server) {
	pb.RegisterCommsServer(server, new(commsNode))
}
