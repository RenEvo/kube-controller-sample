package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"k8s.io/client-go/kubernetes"
	api "k8s.io/api/core/v1"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/golang/glog"

	"github.com/renevo/kube-controller/internal/controller"
)

func main() {

	// When running as a pod in-cluster, a kubeconfig is not needed. Instead this will make use of the service account injected into the pod.
	// However, allow the use of a local kubeconfig as this can make local development & testing easier.
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig file")
	namespace := flag.String("namespace", api.NamespaceAll, "Namespace to use")

	// We log to stderr because glog will default to logging to a file.
	// By setting this debugging is easier via `kubectl logs`
	flag.Set("logtostderr", "true")
	flag.Parse()

	// Build the client config - optionally using a provided kubeconfig file.
	config, err := controller.GetClientConfig(*kubeconfig)
	if err != nil {
		glog.Fatalf("Failed to load client config: %v", err)
	}

	// Construct the Kubernetes client
	client, err := kubernetes.NewForConfig(config)
	if err != nil {
		glog.Fatalf("Failed to create kubernetes client: %v", err)
	}

	nodeClient := client.Core().Nodes()
	deploymentsClient := client.Apps().Deployments(*namespace)

	signalCh := make(chan os.Signal, 10)
	signal.Notify(signalCh, os.Interrupt, os.Kill)
	
	t := time.NewTicker(time.Second)
	defer t.Stop()

	if _, err := controller.CreateWeb(deploymentsClient);  err != nil {
		glog.Fatalf("Failed to create web deployment: %v", err)
	}

	if service, err := controller.CreateService(client.Core().Services(*namespace));  err != nil {
		glog.Fatalf("Failed to create service deployment: %v", err)
	} else {
		glog.Infof("Service created: %v", service.Spec.Ports[0].NodePort)
	}

	LOOP:
	for {
		select {
		case <-signalCh:
			break LOOP

		case <-t.C:
			nodeList, err := nodeClient.List(meta.ListOptions{})
			if err != nil {
				glog.Fatalf("Failed to list nodes: %v", err)
			}
		
			for _, node := range nodeList.Items {
				glog.Infof("Node: %s/%s", node.Namespace, node.Name)
			}

			deployList, err := deploymentsClient.List(meta.ListOptions{})
			if err != nil {
				glog.Fatalf("Failed to list deployments: %v", err)
			}

			for _, deploy := range deployList.Items {
				glog.Infof("Deployment: %s/%s", deploy.Namespace, deploy.Name)
			}
		}
	}

	if err := controller.DeleteService(client.Core().Services(*namespace)); err != nil{
		glog.Errorf("Failed to delete service deployment: %v", err)
	}

	if err := controller.DeleteWeb(deploymentsClient); err != nil{
		glog.Errorf("Failed to delete web deployment: %v", err)
	}

	glog.Infof("Done")
}
