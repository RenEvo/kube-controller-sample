# Kube Controller Prototoype

This is just a simple kube service that will eventually control something. Right now it just outputs the name of each node in each namespace.

## Running locally (this assumes you have some sort of kubernetes running locally):

```bash
go run main.go -kubeconfig=<userpath>/.kube/config
```

## Running on cluster (local one)

```bash
./k8run
```

This will build and run the container in k8s, first deleting the existing deployment `kube-controller-sample` in the cluster.