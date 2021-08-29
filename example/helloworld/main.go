package main

import (
	"context"
	"fmt"
	gin "github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/zouchunxu/gof/example/helloworld/api"
	"github.com/zouchunxu/gof/example/helloworld/config"
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
	var cfg config.Config
	s := server.New("/Users/zouchunxu/web/docker/wwwroot/go/gof/example/helloworld/app.yaml")
	g := gin.New()
	svc := GreetSvc{}
	api.RegisterGreetServer(s.GrpcSever, svc)
	api.NewGreetHandler(svc, g)
	if err := viper.Unmarshal(&cfg); err != nil {
		panic(err)
	}
	go func() {
		g.Run(":5906")
	}()
	s.Log.Info("aaaa")
	if err := s.Run(); err != nil {
		log.Fatalln(err.Error())
	}
}
