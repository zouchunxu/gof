package svr

import (
	"github.com/zouchunxu/deployment/api"
	"github.com/zouchunxu/deployment/internal/service"
	"github.com/zouchunxu/gof"
)

func Init(app *gof.App) {
	deploySvc := service.NewDeployService(app)
	// 部署服务
	api.RegisterDeployServer(app.GrpcSever, deploySvc)
	api.NewDeployHandler(deploySvc, app.G)
}
