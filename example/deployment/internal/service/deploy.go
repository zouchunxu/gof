package service

import (
	"context"
	"github.com/zouchunxu/deployment/api"
	"github.com/zouchunxu/deployment/internal/biz"
	"github.com/zouchunxu/deployment/internal/data"
	"github.com/zouchunxu/gof"
	"github.com/zouchunxu/gof/pkg/api_errors"
)

func NewDeployService(app *gof.App) api.DeployServer {
	return &DeployService{
		app:    app,
		deploy: biz.NewDeployUsecase(app, data.NewDeployRepo(app)),
	}
}

type DeployService struct {
	app    *gof.App
	deploy *biz.DeployUsecase
}

func (d *DeployService) List(ctx context.Context, req *api.DeployListReq) (*api.DeployListRsp, error) {
	list, err := d.deploy.List(req.Page, req.PageSize)
	if err != nil {
		return nil, api_errors.FromError(err)
	}
	rsp := &api.DeployListRsp{}
	for _, item := range list {
		rsp.List = append(rsp.List, &api.DeployListRsp_Row{
			Name:  item.Name,
			Image: item.Image,
		})
	}
	return rsp, nil
}

func (d *DeployService) Create(ctx context.Context, req *api.DeployCreateReq) (*api.DeployCreateRsp, error) {
	err := d.deploy.Create(&biz.Deploy{
		Name:  req.Name,
		Image: req.Image,
	})
	if err != nil {
		return nil, api_errors.FromError(err)
	}
	return &api.DeployCreateRsp{}, nil
}

func (d *DeployService) Update(ctx context.Context, req *api.DeployUpdateReq) (*api.DeployUpdateRsp, error) {
	err := d.deploy.Update(req.ID, &biz.Deploy{
		Name:  req.Name,
		Image: req.Image,
	})
	if err != nil {
		return nil, api_errors.FromError(err)
	}
	return &api.DeployUpdateRsp{}, nil
}

func (d *DeployService) Delete(ctx context.Context, req *api.DeployDeleteReq) (*api.DeployDeleteRsp, error) {
	err := d.deploy.Delete(req.Id)
	if err != nil {
		return nil, api_errors.FromError(err)
	}
	return &api.DeployDeleteRsp{}, nil
}
