package comms

import (
	pb "Awesome-DFS/serverscomms"
	"context"
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
		return &pb.PingResponse{Status: pb.Status_STATUS_NOT_READY}, nil
	}

	log.Printf("Received ping from %s\n", p.Addr)

	return &pb.PingResponse{Status: pb.Status_STATUS_READY}, nil
}

func RegisterCommsServer(s *grpc.Server) {
	pb.RegisterCommsServer(s, new(commsNode))
}
