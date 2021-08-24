package main

import (
	"context"
	"fmt"
	"github.com/zouchunxu/gof/client"
	"github.com/zouchunxu/gof/errors"
	"github.com/zouchunxu/gof/example/helloworld/api"
	"github.com/zouchunxu/gof/server"
	"log"
	"time"
)

type GreetSvc struct {
}

func (g GreetSvc) SayHello(ctx context.Context, req *api.SayHelloReq) (*api.SayHelloRsp, error) {
	if req.Name == "error" {
		return nil, errors.BadRequest("helloworld", "custom_error", fmt.Sprintf("invalid argument %s", req.Name))
	}
	return &api.SayHelloRsp{
		Name: req.Name,
	}, nil
}

func main() {
	s := server.New("helloworld")
	svc := GreetSvc{}
	api.RegisterGreetServer(s.GrpcSever, svc)
	s.Log.Info("aaaa")
	go func() {
		if err := s.Run(":5903"); err != nil {
			log.Fatalln(err.Error())
		}
	}()
	time.Sleep(1 * time.Second)
	conn, err := client.NewConnect("127.0.0.1:5903")
	if err != nil {
		log.Fatalln(err.Error())
	}
	cs := api.NewGreetClient(conn)
	rsp, err := cs.SayHello(context.Background(), &api.SayHelloReq{
		Name: "error",
	})
	if err != nil {
		s.Log.Errorf("err: %+v", err)
	} else {
		fmt.Println(rsp.Name)
	}
}
