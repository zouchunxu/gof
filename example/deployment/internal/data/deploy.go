package data

import (
	"github.com/zouchunxu/deployment/internal/biz"
	"github.com/zouchunxu/gof/server"
)

type deployRepo struct {
	app *server.App
}

func (d deployRepo) ListDeploy(page, pageSize uint32) ([]*biz.Deploy, error) {
	panic("implement me")
}

func (d deployRepo) CreateDeploy(deploy *biz.Deploy) error {
	panic("implement me")
}

func (d deployRepo) UpdateDeploy(id uint32, deploy *biz.Deploy) error {
	panic("implement me")
}

func (d deployRepo) DeleteDeploy(id uint32) error {
	panic("implement me")
}

func NewDeployRepo(app *server.App) biz.DeployRepo {
	return deployRepo{
		app: app,
	}
}
