package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"github.com/zouchunxu/gof"
	"github.com/zouchunxu/gof/client"
	"github.com/zouchunxu/gof/example/helloworld/api"
	"github.com/zouchunxu/gof/example/helloworld/config"
	"github.com/zouchunxu/gof/pkg/api_errors"
	"github.com/zouchunxu/gof/registry"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
	"time"
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
	go func() {
		var cfg config.Config
		s := gof.New("/Users/zouchunxu/web/docker/wwwroot/go/gof/example/helloworld/app.yaml")
		//g := gin.New()
		svc := GreetSvc{}
		api.RegisterGreetServer(s.GrpcSever, svc)
		//api.NewGreetHandler(svc, g)
		if err := viper.Unmarshal(&cfg); err != nil {
			panic(err)
		}
		//go func() {
		//	g.Run(":5906")
		//}()
		go func() {
			s.Log.Info("aaaa")
			if err := s.Run(); err != nil {
				log.Fatalln(err.Error())
			}
		}()

		cli, cerr := clientv3.NewFromURL("http://127.0.0.1:2379")
		if cerr != nil {
			panic(cerr.Error())
		}

		if err := registry.Register(cli, "foo/bar/my-service", "127.0.0.1:5903"); err != nil {
			panic(err.Error())
		}
	}()

	time.Sleep(1 * time.Second)

	conn, err := client.NewConnect("etcd:///foo/bar/my-service")
	if err != nil {
		panic(err.Error())
	}

	gree := api.NewGreetClient(conn)
	rsp, err := gree.SayHello(context.Background(), &api.SayHelloReq{Name: "aa"})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(rsp.Name)
}
