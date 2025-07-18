// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             (unknown)
// source: oauth2/v1/oauth2.proto

package v1

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
	OAuth2Service_GetAuthURL_FullMethodName     = "/OAuth2Service/GetAuthURL"
	OAuth2Service_HandleCallback_FullMethodName = "/OAuth2Service/HandleCallback"
)

// OAuth2ServiceClient is the client API for OAuth2Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type OAuth2ServiceClient interface {
	// 获取授权 URL
	GetAuthURL(ctx context.Context, in *GetAuthURLRequest, opts ...grpc.CallOption) (*GetAuthURLResponse, error)
	// 处理回调并获取 Token
	HandleCallback(ctx context.Context, in *HandleCallbackRequest, opts ...grpc.CallOption) (*HandleCallbackResponse, error)
}

type oAuth2ServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewOAuth2ServiceClient(cc grpc.ClientConnInterface) OAuth2ServiceClient {
	return &oAuth2ServiceClient{cc}
}

func (c *oAuth2ServiceClient) GetAuthURL(ctx context.Context, in *GetAuthURLRequest, opts ...grpc.CallOption) (*GetAuthURLResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetAuthURLResponse)
	err := c.cc.Invoke(ctx, OAuth2Service_GetAuthURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *oAuth2ServiceClient) HandleCallback(ctx context.Context, in *HandleCallbackRequest, opts ...grpc.CallOption) (*HandleCallbackResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HandleCallbackResponse)
	err := c.cc.Invoke(ctx, OAuth2Service_HandleCallback_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// OAuth2ServiceServer is the server API for OAuth2Service service.
// All implementations must embed UnimplementedOAuth2ServiceServer
// for forward compatibility.
type OAuth2ServiceServer interface {
	// 获取授权 URL
	GetAuthURL(context.Context, *GetAuthURLRequest) (*GetAuthURLResponse, error)
	// 处理回调并获取 Token
	HandleCallback(context.Context, *HandleCallbackRequest) (*HandleCallbackResponse, error)
	mustEmbedUnimplementedOAuth2ServiceServer()
}

// UnimplementedOAuth2ServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedOAuth2ServiceServer struct{}

func (UnimplementedOAuth2ServiceServer) GetAuthURL(context.Context, *GetAuthURLRequest) (*GetAuthURLResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAuthURL not implemented")
}
func (UnimplementedOAuth2ServiceServer) HandleCallback(context.Context, *HandleCallbackRequest) (*HandleCallbackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleCallback not implemented")
}
func (UnimplementedOAuth2ServiceServer) mustEmbedUnimplementedOAuth2ServiceServer() {}
func (UnimplementedOAuth2ServiceServer) testEmbeddedByValue()                       {}

// UnsafeOAuth2ServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to OAuth2ServiceServer will
// result in compilation errors.
type UnsafeOAuth2ServiceServer interface {
	mustEmbedUnimplementedOAuth2ServiceServer()
}

func RegisterOAuth2ServiceServer(s grpc.ServiceRegistrar, srv OAuth2ServiceServer) {
	// If the following call pancis, it indicates UnimplementedOAuth2ServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&OAuth2Service_ServiceDesc, srv)
}

func _OAuth2Service_GetAuthURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAuthURLRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OAuth2ServiceServer).GetAuthURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OAuth2Service_GetAuthURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OAuth2ServiceServer).GetAuthURL(ctx, req.(*GetAuthURLRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _OAuth2Service_HandleCallback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandleCallbackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(OAuth2ServiceServer).HandleCallback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: OAuth2Service_HandleCallback_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(OAuth2ServiceServer).HandleCallback(ctx, req.(*HandleCallbackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// OAuth2Service_ServiceDesc is the grpc.ServiceDesc for OAuth2Service service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var OAuth2Service_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "OAuth2Service",
	HandlerType: (*OAuth2ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetAuthURL",
			Handler:    _OAuth2Service_GetAuthURL_Handler,
		},
		{
			MethodName: "HandleCallback",
			Handler:    _OAuth2Service_HandleCallback_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "oauth2/v1/oauth2.proto",
}
