package main

import (
	"context"
	"fmt"
	gin "github.com/gin-gonic/gin"
	"github.com/zouchunxu/gof/example/helloworld/api"
	"github.com/zouchunxu/gof/pkg/api_errors"
	"github.com/zouchunxu/gof/server"
	"log"
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
	s := server.New("helloworld")
	g := gin.New()
	svc := GreetSvc{}
	api.RegisterGreetServer(s.GrpcSever, svc)
	api.NewGreetHandler(svc, g)
	go func() {
		g.Run(":5906")
	}()
	s.Log.Info("aaaa")
	if err := s.Run(":5903"); err != nil {
		log.Fatalln(err.Error())
	}
}
