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
	protoregistry.GlobalFiles.RangeFiles(func(fd protoreflect.FileDescriptor) bool {
		dep := protodesc.ToFileDescriptorProto(fd)
		for _, str := range dep.Dependency {
			fmt.Println(str)
		}
		return true
	})
}
