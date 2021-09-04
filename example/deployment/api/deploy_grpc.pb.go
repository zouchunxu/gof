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

// DeployClient is the client API for Deploy service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type DeployClient interface {
	List(ctx context.Context, in *DeployListReq, opts ...grpc.CallOption) (*DeployListRsp, error)
	Create(ctx context.Context, in *DeployCreateReq, opts ...grpc.CallOption) (*DeployCreateRsp, error)
	Update(ctx context.Context, in *DeployUpdateReq, opts ...grpc.CallOption) (*DeployUpdateRsp, error)
	Delete(ctx context.Context, in *DeployDeleteReq, opts ...grpc.CallOption) (*DeployDeleteRsp, error)
}

type deployClient struct {
	cc grpc.ClientConnInterface
}

func NewDeployClient(cc grpc.ClientConnInterface) DeployClient {
	return &deployClient{cc}
}

func (c *deployClient) List(ctx context.Context, in *DeployListReq, opts ...grpc.CallOption) (*DeployListRsp, error) {
	out := new(DeployListRsp)
	err := c.cc.Invoke(ctx, "/api.Deploy/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deployClient) Create(ctx context.Context, in *DeployCreateReq, opts ...grpc.CallOption) (*DeployCreateRsp, error) {
	out := new(DeployCreateRsp)
	err := c.cc.Invoke(ctx, "/api.Deploy/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deployClient) Update(ctx context.Context, in *DeployUpdateReq, opts ...grpc.CallOption) (*DeployUpdateRsp, error) {
	out := new(DeployUpdateRsp)
	err := c.cc.Invoke(ctx, "/api.Deploy/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deployClient) Delete(ctx context.Context, in *DeployDeleteReq, opts ...grpc.CallOption) (*DeployDeleteRsp, error) {
	out := new(DeployDeleteRsp)
	err := c.cc.Invoke(ctx, "/api.Deploy/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeployServer is the server API for Deploy service.
// All implementations should embed UnimplementedDeployServer
// for forward compatibility
type DeployServer interface {
	List(context.Context, *DeployListReq) (*DeployListRsp, error)
	Create(context.Context, *DeployCreateReq) (*DeployCreateRsp, error)
	Update(context.Context, *DeployUpdateReq) (*DeployUpdateRsp, error)
	Delete(context.Context, *DeployDeleteReq) (*DeployDeleteRsp, error)
}

// UnimplementedDeployServer should be embedded to have forward compatible implementations.
type UnimplementedDeployServer struct {
}

func (UnimplementedDeployServer) List(context.Context, *DeployListReq) (*DeployListRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedDeployServer) Create(context.Context, *DeployCreateReq) (*DeployCreateRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedDeployServer) Update(context.Context, *DeployUpdateReq) (*DeployUpdateRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedDeployServer) Delete(context.Context, *DeployDeleteReq) (*DeployDeleteRsp, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}

// UnsafeDeployServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to DeployServer will
// result in compilation errors.
type UnsafeDeployServer interface {
	mustEmbedUnimplementedDeployServer()
}

func RegisterDeployServer(s grpc.ServiceRegistrar, srv DeployServer) {
	s.RegisterService(&Deploy_ServiceDesc, srv)
}

func _Deploy_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployListReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeployServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Deploy/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeployServer).List(ctx, req.(*DeployListReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deploy_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployCreateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeployServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Deploy/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeployServer).Create(ctx, req.(*DeployCreateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deploy_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployUpdateReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeployServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Deploy/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeployServer).Update(ctx, req.(*DeployUpdateReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Deploy_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeployDeleteReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeployServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Deploy/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeployServer).Delete(ctx, req.(*DeployDeleteReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Deploy_ServiceDesc is the grpc.ServiceDesc for Deploy service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Deploy_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Deploy",
	HandlerType: (*DeployServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "List",
			Handler:    _Deploy_List_Handler,
		},
		{
			MethodName: "Create",
			Handler:    _Deploy_Create_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _Deploy_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _Deploy_Delete_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/deploy.proto",
}