package data

import (
	"context"
	"github.com/zouchunxu/deployment/internal/biz"
	"github.com/zouchunxu/gof/server"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type deployRepo struct {
	app *server.App
}

func (d deployRepo) ListDeploy(page, pageSize uint32) ([]*biz.Deploy, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	l, err := clientset.AppsV1().Deployments("").List(context.Background(), metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	var res []*biz.Deploy
	for _, item := range l.Items {
		res = append(res, &biz.Deploy{
			Name:  item.Name,
			Image: item.Name,
		})
	}
	return res, nil
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
