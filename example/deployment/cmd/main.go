package main

import (
	"github.com/zouchunxu/deployment/api"
	"github.com/zouchunxu/deployment/internal/service"
	"github.com/zouchunxu/gof/server"
	"log"
)

func main() {
	app := server.New("/Users/zouchunxu/web/docker/wwwroot/go/gof/example/deployment/app.yaml")
	deploySvc := service.NewDeployService(app)
	// 部署服务
	api.RegisterDeployServer(app.GrpcSever, deploySvc)
	api.NewDeployHandler(deploySvc, app.G)

	if err := app.Run(); err != nil {
		log.Fatal(err.Error())
	}
}
