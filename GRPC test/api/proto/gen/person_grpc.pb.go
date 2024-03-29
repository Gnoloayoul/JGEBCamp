// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: person.proto

package personv1

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

const (
	PersonAction_SayHello_FullMethodName   = "/person.v1.PersonAction/SayHello"
	PersonAction_SayGoodBye_FullMethodName = "/person.v1.PersonAction/SayGoodBye"
)

// PersonActionClient is the client API for PersonAction service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PersonActionClient interface {
	SayHello(ctx context.Context, in *SayHelloRequest, opts ...grpc.CallOption) (*SayHelloResponse, error)
	SayGoodBye(ctx context.Context, in *SayGoodByeRequest, opts ...grpc.CallOption) (*SayGoodByeResponse, error)
}

type personActionClient struct {
	cc grpc.ClientConnInterface
}

func NewPersonActionClient(cc grpc.ClientConnInterface) PersonActionClient {
	return &personActionClient{cc}
}

func (c *personActionClient) SayHello(ctx context.Context, in *SayHelloRequest, opts ...grpc.CallOption) (*SayHelloResponse, error) {
	out := new(SayHelloResponse)
	err := c.cc.Invoke(ctx, PersonAction_SayHello_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *personActionClient) SayGoodBye(ctx context.Context, in *SayGoodByeRequest, opts ...grpc.CallOption) (*SayGoodByeResponse, error) {
	out := new(SayGoodByeResponse)
	err := c.cc.Invoke(ctx, PersonAction_SayGoodBye_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PersonActionServer is the server API for PersonAction service.
// All implementations must embed UnimplementedPersonActionServer
// for forward compatibility
type PersonActionServer interface {
	SayHello(context.Context, *SayHelloRequest) (*SayHelloResponse, error)
	SayGoodBye(context.Context, *SayGoodByeRequest) (*SayGoodByeResponse, error)
	mustEmbedUnimplementedPersonActionServer()
}

// UnimplementedPersonActionServer must be embedded to have forward compatible implementations.
type UnimplementedPersonActionServer struct {
}

func (UnimplementedPersonActionServer) SayHello(context.Context, *SayHelloRequest) (*SayHelloResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedPersonActionServer) SayGoodBye(context.Context, *SayGoodByeRequest) (*SayGoodByeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayGoodBye not implemented")
}
func (UnimplementedPersonActionServer) mustEmbedUnimplementedPersonActionServer() {}

// UnsafePersonActionServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PersonActionServer will
// result in compilation errors.
type UnsafePersonActionServer interface {
	mustEmbedUnimplementedPersonActionServer()
}

func RegisterPersonActionServer(s grpc.ServiceRegistrar, srv PersonActionServer) {
	s.RegisterService(&PersonAction_ServiceDesc, srv)
}

func _PersonAction_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SayHelloRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonActionServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PersonAction_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonActionServer).SayHello(ctx, req.(*SayHelloRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PersonAction_SayGoodBye_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SayGoodByeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PersonActionServer).SayGoodBye(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PersonAction_SayGoodBye_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PersonActionServer).SayGoodBye(ctx, req.(*SayGoodByeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PersonAction_ServiceDesc is the grpc.ServiceDesc for PersonAction service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PersonAction_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "person.v1.PersonAction",
	HandlerType: (*PersonActionServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _PersonAction_SayHello_Handler,
		},
		{
			MethodName: "SayGoodBye",
			Handler:    _PersonAction_SayGoodBye_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "person.proto",
}
