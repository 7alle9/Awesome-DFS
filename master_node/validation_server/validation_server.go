package validation_server

import (
	ms "Awesome-DFS/master_node/metadata_service"
	"Awesome-DFS/protobuf/validation"
	"context"
	"google.golang.org/grpc"
)

type ValidationServer struct {
	__.UnimplementedValidationServer
}

func (s *ValidationServer) Validate(_ context.Context, req *__.ValidationRequest) (*__.Empty, error) {
	ms.Validate(req.FileUuid)
	return &__.Empty{}, nil
}

func RegisterValidationServer(server *grpc.Server) {
	__.RegisterValidationServer(server, &ValidationServer{})
}
