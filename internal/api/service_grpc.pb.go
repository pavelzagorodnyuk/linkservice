// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LinkServiceClient is the client API for LinkService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LinkServiceClient interface {
	Create(ctx context.Context, in *URL, opts ...grpc.CallOption) (*Link, error)
	Get(ctx context.Context, in *Link, opts ...grpc.CallOption) (*URL, error)
}

type linkServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewLinkServiceClient(cc grpc.ClientConnInterface) LinkServiceClient {
	return &linkServiceClient{cc}
}

func (c *linkServiceClient) Create(ctx context.Context, in *URL, opts ...grpc.CallOption) (*Link, error) {
	out := new(Link)
	err := c.cc.Invoke(ctx, "/api.LinkService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *linkServiceClient) Get(ctx context.Context, in *Link, opts ...grpc.CallOption) (*URL, error) {
	out := new(URL)
	err := c.cc.Invoke(ctx, "/api.LinkService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LinkServiceServer is the server API for LinkService service.
// All implementations must embed UnimplementedLinkServiceServer
// for forward compatibility
type LinkServiceServer interface {
	Create(context.Context, *URL) (*Link, error)
	Get(context.Context, *Link) (*URL, error)
	mustEmbedUnimplementedLinkServiceServer()
}

// UnimplementedLinkServiceServer must be embedded to have forward compatible implementations.
type UnimplementedLinkServiceServer struct {
}

func (UnimplementedLinkServiceServer) Create(context.Context, *URL) (*Link, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedLinkServiceServer) Get(context.Context, *Link) (*URL, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedLinkServiceServer) mustEmbedUnimplementedLinkServiceServer() {}

// UnsafeLinkServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LinkServiceServer will
// result in compilation errors.
type UnsafeLinkServiceServer interface {
	mustEmbedUnimplementedLinkServiceServer()
}

func RegisterLinkServiceServer(s grpc.ServiceRegistrar, srv LinkServiceServer) {
	s.RegisterService(&LinkService_ServiceDesc, srv)
}

func _LinkService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(URL)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.LinkService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServiceServer).Create(ctx, req.(*URL))
	}
	return interceptor(ctx, in, info, handler)
}

func _LinkService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Link)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LinkServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.LinkService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LinkServiceServer).Get(ctx, req.(*Link))
	}
	return interceptor(ctx, in, info, handler)
}

// LinkService_ServiceDesc is the grpc.ServiceDesc for LinkService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LinkService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.LinkService",
	HandlerType: (*LinkServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _LinkService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _LinkService_Get_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/service.proto",
}