[![Go Report Card](https://goreportcard.com/badge/kubeops.dev/falco-ui-server)](https://goreportcard.com/report/kubeops.dev/falco-ui-server)
[![Build Status](https://github.com/kubeops/falco-ui-server/workflows/CI/badge.svg)](https://github.com/kubeops/falco-ui-server/actions?workflow=CI)
[![Twitter](https://img.shields.io/twitter/follow/kubeops.svg?style=social&logo=twitter&label=Follow)](https://twitter.com/intent/follow?screen_name=Kubeops)

# falco-ui-server

falco-ui-server is an [extended api server](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/apiserver-aggregation/) that reports [Falco](https://falco.org/) events as a Kubernetes native resource.

## Deploy into a Kubernetes Cluster

You can deploy `falco-ui-server` using Helm chart found [here](https://github.com/kubeops/installer/tree/master/charts/falco-ui-server).

```console
helm repo add appscode https://charts.appscode.com/stable/
helm repo update

helm install falco-ui-server appscode/falco-ui-server
```
