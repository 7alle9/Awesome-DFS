package validation_server

import (
	ms "Awesome-DFS/master_node/metadata_service"
	pb "Awesome-DFS/validation"
	"context"
	"google.golang.org/grpc"
)

type ValidationServer struct {
	pb.UnimplementedValidationServer
}

func (s *ValidationServer) Validate(_ context.Context, req *pb.ValidationRequest) (*pb.Empty, error) {
	ms.Validate(req.FileUuid)
	return &pb.Empty{}, nil
}

func RegisterValidationServer(server *grpc.Server) {
	pb.RegisterValidationServer(server, &ValidationServer{})
}
