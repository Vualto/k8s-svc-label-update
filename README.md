## Docker Cloud - Kubernetes Continuous Delivery API

#### go build

platforms: `linux`, `macos`, `windows`

`./build.sh {platform}`

#### usage

* Add needed resources to cluster (see examples/continuous-delivery.yaml)

```
kubectl apply -f examples/continuous-delivery-install.yaml
```

* Create a ContinuousDelivery resource linking Docker Hub repo with container(s) of a cluster resource (pods, deployments, ...) (see examples/my-docker-repo-continuous-delivery.yaml)

```
kubectl apply -f examples/my-docker-repo-continuous-delivery.yaml
```

> Important: ContinuousDelivery.metadata.name must be set to docker repo name (example is `my-docker-repo`)

#### TODO

SECURITY, VERIFY DOCKER BUILD:
- check posted build tag/pushDate matches https://hub.docker.com/v2/repositories/vualto/reponame/tags/tag name/pushDate
