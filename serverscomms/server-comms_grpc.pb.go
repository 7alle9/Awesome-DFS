// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.28.2
// source: server-comms.proto

package __

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Comms_Ping_FullMethodName = "/serverscomms.Comms/ping"
)

// CommsClient is the client API for Comms service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CommsClient interface {
	Ping(ctx context.Context, in *PingPayload, opts ...grpc.CallOption) (*PingResponse, error)
}

type commsClient struct {
	cc grpc.ClientConnInterface
}

func NewCommsClient(cc grpc.ClientConnInterface) CommsClient {
	return &commsClient{cc}
}

func (c *commsClient) Ping(ctx context.Context, in *PingPayload, opts ...grpc.CallOption) (*PingResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PingResponse)
	err := c.cc.Invoke(ctx, Comms_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CommsServer is the server API for Comms service.
// All implementations must embed UnimplementedCommsServer
// for forward compatibility.
type CommsServer interface {
	Ping(context.Context, *PingPayload) (*PingResponse, error)
	mustEmbedUnimplementedCommsServer()
}

// UnimplementedCommsServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedCommsServer struct{}

func (UnimplementedCommsServer) Ping(context.Context, *PingPayload) (*PingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedCommsServer) mustEmbedUnimplementedCommsServer() {}
func (UnimplementedCommsServer) testEmbeddedByValue()               {}

// UnsafeCommsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CommsServer will
// result in compilation errors.
type UnsafeCommsServer interface {
	mustEmbedUnimplementedCommsServer()
}

func RegisterCommsServer(s grpc.ServiceRegistrar, srv CommsServer) {
	// If the following call pancis, it indicates UnimplementedCommsServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Comms_ServiceDesc, srv)
}

func _Comms_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PingPayload)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CommsServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Comms_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CommsServer).Ping(ctx, req.(*PingPayload))
	}
	return interceptor(ctx, in, info, handler)
}

// Comms_ServiceDesc is the grpc.ServiceDesc for Comms service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Comms_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "serverscomms.Comms",
	HandlerType: (*CommsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ping",
			Handler:    _Comms_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server-comms.proto",
}
