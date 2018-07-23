## K8s setup

```bash
eval $(minikube docker-env)
make build-nounparser build-frontend
make push-nounparser push-frontend

make k8s-create
# make k8s-test-<any word>, e.g.
make k8s-test-Ísafjörður
```

## Local setup

```bash
make build-nounparser build-frontend
make run-nounparser run-frontend
# curl localhost:8000/<any word>, e.g.
curl localhost:8000/Ísafjörður

docker logs frontend  # etc.

make stop
```
