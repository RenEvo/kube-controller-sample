package main

import (
	"flag"
	"os"
	"os/signal"
	"time"

	"k8s.io/client-go/kubernetes"
	meta "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/golang/glog"

	"github.com/renevo/kube-controller/internal/controller"
)

func main() {

	// When running as a pod in-cluster, a kubeconfig is not needed. Instead this will make use of the service account injected into the pod.
	// However, allow the use of a local kubeconfig as this can make local development & testing easier.
	kubeconfig := flag.String("kubeconfig", "", "Path to a kubeconfig file")

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

	signalCh := make(chan os.Signal, 10)
	signal.Notify(signalCh, os.Interrupt, os.Kill)
	
	t := time.NewTicker(time.Second)
	defer t.Stop()

	LOOP:
	for {
		select {
		case <-signalCh:
			break LOOP

		case <-t.C:
			list, err := client.Core().Nodes().List(meta.ListOptions{})
			if err != nil {
				glog.Fatalf("Failed to list nodes: %v", err)
			}
		
			for _, node := range list.Items {
				glog.Infof("%s/%s", node.Namespace, node.Name)
			}
		}
	}

	glog.Infof("Done")
}
