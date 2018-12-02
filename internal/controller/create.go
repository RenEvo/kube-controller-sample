package controller

import (
	"github.com/pkg/errors"
	apps "k8s.io/api/apps/v1"
	api "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	appsTyped "k8s.io/client-go/kubernetes/typed/apps/v1"
	apiTyped "k8s.io/client-go/kubernetes/typed/core/v1"

	k8errors "k8s.io/apimachinery/pkg/api/errors"
)

const webname = "kube-controller-sample-web"

func CreateService(serviceClient apiTyped.ServiceInterface) (*api.Service, error) {
	execFunc := serviceClient.Create

	current, err := serviceClient.Get(webname, meta.GetOptions{IncludeUninitialized: true})
	if err != nil && !k8errors.IsNotFound(err) {
		return nil, errors.Wrapf(err, "failed to get current service for %s", webname)
	}

	if current != nil && err == nil {
		execFunc = serviceClient.Update
	}

	fresh := createService(webname)

	return execFunc(fresh)
}

func DeleteService(serviceClient apiTyped.ServiceInterface) error {
	return serviceClient.Delete(webname, &meta.DeleteOptions{})
}

func CreateWeb(deployClient appsTyped.DeploymentInterface) (*apps.Deployment, error) {
	execFunc := deployClient.Create

	current, err := deployClient.Get(webname, meta.GetOptions{IncludeUninitialized: true})
	if err != nil && !k8errors.IsNotFound(err) {
		return nil, errors.Wrapf(err, "failed to get current deployment for %s", webname)
	}

	// API returns an error, and a value
	if current != nil && err == nil {
		execFunc = deployClient.Update
	}

	fresh := createWebDeployment(webname, "alpine")

	return execFunc(fresh)
}

func DeleteWeb(deployClient appsTyped.DeploymentInterface) error {
	return deployClient.Delete(webname, &meta.DeleteOptions{})
}

func createService(name string) *api.Service {
	return &api.Service{
		ObjectMeta: meta.ObjectMeta{
			Name: webname,
		},
		Spec: api.ServiceSpec{
			Type: api.ServiceTypeNodePort,
			Selector: map[string]string{
				"app": "nginx",
			},
			Ports: []api.ServicePort{
				{
					Port: 80,
				},
			},
		},
	}
}

func createWebDeployment(name string, version string) *apps.Deployment {
	return &apps.Deployment{
		ObjectMeta: meta.ObjectMeta{
			Name: name,
		},
		Spec: apps.DeploymentSpec{
			Replicas: int32Ptr(2),
			Selector: &meta.LabelSelector{
				MatchLabels: map[string]string{
					"app": "nginx",
				},
			},
			Template: api.PodTemplateSpec{
				ObjectMeta: meta.ObjectMeta{
					Labels: map[string]string{
						"app": "nginx",
					},
				},
				Spec: api.PodSpec{
					Containers: []api.Container{
						{
							Name:  "web",
							Image: "nginx:" + version,
							Ports: []api.ContainerPort{
								{
									Name:          "http",
									Protocol:      api.ProtocolTCP,
									ContainerPort: 80,
								},
							},
						},
					},
				},
			},
		},
	}
}
