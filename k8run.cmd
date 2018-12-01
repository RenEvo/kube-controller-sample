go mod vendor
docker -H tcp://127.0.0.1:2375 build --force-rm --tag "renevo/kube-controller:latest" .
kubectl delete deployments kube-controller
kubectl run kube-controller --image renevo/kube-controller:latest --image-pull-policy=Never -- /usr/local/bin/kube-controller
rd /s /q .\vendor