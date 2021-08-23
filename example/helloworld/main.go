package main

import (
	"context"
	"github.com/zouchunxu/gof/example/helloworld/api"
	"github.com/zouchunxu/gof/server"
	"log"
)

type GreetSvc struct {
}

func (g GreetSvc) SayHello(ctx context.Context, req *api.SayHelloReq) (*api.SayHelloRsp, error) {
	panic("implement me")
}

func main() {
	s := server.New("helloworld")
	svc := GreetSvc{}
	api.RegisterGreetServer(s.GrpcSever, svc)
	s.Log.Info("aaaa")
	if err := s.Run(":5903"); err != nil {
		log.Fatalln(err.Error())
	}
}
