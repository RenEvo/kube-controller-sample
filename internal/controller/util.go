package controller

import (
	"fmt"
	"strings"
	"os"

	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"

	k8errors "k8s.io/apimachinery/pkg/api/errors"
)

func int32Ptr(i int32) *int32 { return &i }

func isNotFound(err error) bool {
	statusErr, ok := err.(*k8errors.StatusError)
	if !ok {
		// super dumb that the errors don't actaully line up with the errors ....
		return strings.HasSuffix(err.Error(), " not found")
	}

	fmt.Fprintf(os.Stdout, "Err: %+v\n", statusErr.Status())

	return statusErr.Status().Code == 404
}

func GetClientConfig(kubeconfig string) (*rest.Config, error) {
	if kubeconfig != "" {
		return clientcmd.BuildConfigFromFlags("", kubeconfig)
	}
	return rest.InClusterConfig()
}
