package biz

import (
	"github.com/zouchunxu/gof"
)

type Deploy struct {
	Name  string
	Image string
}

type DeployRepo interface {
	ListDeploy(page, pageSize uint32) ([]*Deploy, error)
	CreateDeploy(deploy *Deploy) error
	UpdateDeploy(id uint32, deploy *Deploy) error
	DeleteDeploy(id uint32) error
}

type DeployUsecase struct {
	app  *gof.App
	repo DeployRepo
}

func NewDeployUsecase(app *gof.App, repo DeployRepo) *DeployUsecase {
	return &DeployUsecase{
		app:  app,
		repo: repo,
	}
}

func (uc *DeployUsecase) List(page, pageSize uint32) ([]*Deploy, error) {
	return uc.repo.ListDeploy(page, pageSize)
}

func (uc *DeployUsecase) Create(deploy *Deploy) error {
	return uc.repo.CreateDeploy(deploy)
}

func (uc *DeployUsecase) Update(id uint32, deploy *Deploy) error {
	return uc.repo.UpdateDeploy(id, deploy)
}

func (uc *DeployUsecase) Delete(id uint32) error {
	return uc.repo.DeleteDeploy(id)
}
