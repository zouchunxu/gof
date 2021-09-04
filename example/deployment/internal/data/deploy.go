package data

import (
	"context"
	"github.com/zouchunxu/deployment/internal/biz"
	"github.com/zouchunxu/gof/server"
	v1 "k8s.io/api/apps/v1"
	cv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	dv1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type deployRepo struct {
	app *server.App
}

func (d deployRepo) ListDeploy(page, pageSize uint32) ([]*biz.Deploy, error) {
	// creates the clientset
	clientset, err := getK8sConfig()
	if err != nil {
		return nil, err
	}
	l, err := clientset.AppsV1().Deployments("default").List(context.Background(), metav1.ListOptions{})
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
	// creates the clientset
	clientset, err := getK8sConfig()
	if err != nil {
		d.app.Log.Errorf("get err: %+v", err)
		return err
	}
	var replicas int32 = 1
	ds := &v1.Deployment{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Deployment",
			APIVersion: "apps/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: deploy.Name,
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"app": deploy.Name,
				},
			},
			Replicas: &replicas,
			Template: cv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"app": deploy.Name,
					},
				},
				Spec: cv1.PodSpec{
					Containers: []cv1.Container{
						{
							Name:            deploy.Name,
							Image:           deploy.Image,
							ImagePullPolicy: "IfNotPresent",
							Ports: []cv1.ContainerPort{
								{ContainerPort: 80},
								{ContainerPort: 5903},
								{ContainerPort: 5904},
								{ContainerPort: 5906},
								{ContainerPort: 9909},
							},
						},
					},
				},
			},
		},
	}
	_, err = clientset.AppsV1().Deployments("default").
		Create(context.Background(), ds, metav1.CreateOptions{})
	if err != nil {
		d.app.Log.Errorf("create err: %+v", err)
		return err
	}
	svc := &cv1.Service{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Service",
			APIVersion: "v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: deploy.Name,
			Labels: map[string]string{
				"app": deploy.Name,
			},
		},
		Spec: cv1.ServiceSpec{
			Selector: map[string]string{
				"app": deploy.Name,
			},
			Ports: []cv1.ServicePort{
				{
					Port: 80,
					Name: "nginx",
				},
				{
					Port: 5903,
					Name: "grpc",
				},
				{
					Port: 5904,
					Name: "prometheus",
				},
				{
					Port: 5906,
					Name: "http",
				},
			},
			ClusterIP: "None",
		},
	}
	_, err = clientset.CoreV1().Services("default").
		Create(context.Background(), svc, metav1.CreateOptions{})
	if err != nil {
		d.app.Log.Errorf("svc err: %+v", err)
	}
	return err
}

func (d deployRepo) UpdateDeploy(id uint32, deploy *biz.Deploy) error {
	k8s, err := getK8sConfig()
	if err != nil {
		return err
	}
	kind := "Deployment"
	v := "apps/v1"
	var cfg = &appsv1.DeploymentApplyConfiguration{
		TypeMetaApplyConfiguration: dv1.TypeMetaApplyConfiguration{
			Kind:       &kind,
			APIVersion: &v,
		},
		ObjectMetaApplyConfiguration: &dv1.ObjectMetaApplyConfiguration{
			Name: &deploy.Name,
		},
		Spec: &appsv1.DeploymentSpecApplyConfiguration{
			Selector: &dv1.LabelSelectorApplyConfiguration{
				MatchLabels: map[string]string{
					"app": deploy.Name,
				},
			},
			Template: &corev1.PodTemplateSpecApplyConfiguration{
				Spec: &corev1.PodSpecApplyConfiguration{
					Containers: []corev1.ContainerApplyConfiguration{
						{
							Name:  &deploy.Name,
							Image: &deploy.Image,
						},
					},
				},
			},
		},
	}
	_, err = k8s.AppsV1().Deployments("default").Apply(context.Background(), cfg, metav1.ApplyOptions{})
	return err
}

func (d deployRepo) DeleteDeploy(id uint32) error {
	panic("implement me")
}

func NewDeployRepo(app *server.App) biz.DeployRepo {
	return deployRepo{
		app: app,
	}
}

func getK8sConfig() (*kubernetes.Clientset, error) {
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
	return clientset, nil
}
