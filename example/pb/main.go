package main

import (
	"context"
	"fmt"
	"github.com/zouchunxu/gof/example/pb/api"
	"github.com/zouchunxu/gof/pkg/api_errors"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
)

type GreetSvc struct {
}

func (g GreetSvc) SayHello(ctx context.Context, req *api.SayHelloReq) (*api.SayHelloRsp, error) {
	if req.Name == "error" {
		return nil, api_errors.BadRequest("helloworld", "custom_error", fmt.Sprintf("invalid argument %s", req.Name))
	}
	return &api.SayHelloRsp{
		Name: req.Name,
	}, nil
}

func main() {
	srv := grpc.NewServer()
	api.RegisterGreetServer(srv, &GreetSvc{})
	info := srv.GetServiceInfo()
	for _, item := range info {
		desc, _ := protoregistry.GlobalFiles.FindFileByPath(item.Metadata.(string))

		fmt.Println(desc.FullName())
		fmt.Println(desc.Services().Get(0).Methods().Get(0).FullName())
		//fmt.Println(desc.Messages())
		//fmt.Println(protodesc.ToFileDescriptorProto(desc).Service)
	}
	if true {
		return
	}

	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		dep := protodesc.ToFileDescriptorProto(fd)
		fmt.Println(dep.Name)
		//for _, str := range dep.Dependency {
		//	fmt.Println(str)
		//	desc, _ := protoregistry.GlobalFiles.FindFileByPath(str)
		//
		//	dp := protodesc.ToFileDescriptorProto(desc)
		//
		//	for _, s := range dp.Dependency {
		//		fmt.Println(s)
		//	}
		//}
		return true
	})
}
