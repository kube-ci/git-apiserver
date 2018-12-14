# Git API Server
[Git API server by AppsCode](https://github.com/kube-ci/git-apiserver) - Sync git repositories as Kubernetes resources.

## TL;DR;

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm install appscode/git-apiserver --name git-apiserver --namespace kube-system
```

## Introduction

This chart bootstraps a [Git API server controller](https://github.com/kube-ci/git-apiserver) deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Kubernetes 1.9+

## Installing the Chart

To install the chart with the release name `git-apiserver`:

```console
$ helm install appscode/git-apiserver --name git-apiserver
```

The command deploys Git API server operator on the Kubernetes cluster in the default configuration. The [configuration](#configuration) section lists the parameters that can be configured during installation.

> **Tip**: List all releases using `helm list`

## Uninstalling the Chart

To uninstall/delete the `git-apiserver`:

```console
$ helm delete git-apiserver
```

The command removes all the Kubernetes components associated with the chart and deletes the release.

## Configuration

The following table lists the configurable parameters of the Git API server chart and their default values.

| Parameter                            | Description                                                       | Default            |
| ------------------------------------ | ----------------------------------------------------------------- | ------------------ |
| `replicaCount`                       | Number of Git API server operator replicas to create (only 1 is supported) | `1`                |
| `operator.registry`                  | Docker registry used to pull operator image                       | `kubeci`         |
| `operator.repository`                | Operator container image                                          | `git-apiserver`            |
| `operator.tag`                       | Operator container image tag                                      | `0.1.0`            |
| `cleaner.registry`                   | Docker registry used to pull Webhook cleaner image                | `appscode`         |
| `cleaner.repository`                 | Webhook cleaner container image                                   | `kubectl`          |
| `cleaner.tag`                        | Webhook cleaner container image tag                               | `v1.11`            |
| `imagePullPolicy`                    | Container image pull policy                                       | `IfNotPresent`     |
| `criticalAddon`                      | If true, installs Git API server operator as critical addon                | `false`            |
| `logLevel`                           | Log level for operator                                            | `3`                |
| `affinity`                           | Affinity rules for pod assignment                                 | `{}`               |
| `annotations`                        | Annotations applied to operator pod(s)                            | `{}`               |
| `nodeSelector`                       | Node labels for pod assignment                                    | `{}`               |
| `tolerations`                        | Tolerations used pod assignment                                   | `{}`               |
| `rbac.create`                        | If `true`, create and use RBAC resources                          | `true`             |
| `serviceAccount.create`              | If `true`, create a new service account                           | `true`             |
| `serviceAccount.name`                | Service account to be used. If not set and `serviceAccount.create` is `true`, a name is generated using the fullname template | `` |
| `apiserver.groupPriorityMinimum`     | The minimum priority the group should have.                       | 10000              |
| `apiserver.versionPriority`          | The ordering of this API inside of the group.                     | 15                 |
| `apiserver.enableValidatingWebhook`  | Enable validating webhooks for Git API server CRDs                         | true               |
| `apiserver.enableMutatingWebhook`    | Enable mutating webhooks for Kubernetes workloads                 | true               |
| `apiserver.ca`                       | CA certificate used by main Kubernetes api server                 | `not-ca-cert`      |
| `apiserver.disableStatusSubresource` | If true, disables status sub resource for crds. Otherwise enables based on Kubernetes version | `false`            |
| `enableAnalytics`                    | Send usage events to Google Analytics                             | `true`             |


Specify each parameter using the `--set key=value[,key=value]` argument to `helm install`. For example:

```console
$ helm install --name git-apiserver --set image.tag=v0.2.1 appscode/git-apiserver
```

Alternatively, a YAML file that specifies the values for the parameters can be provided while
installing the chart. For example:

```console
$ helm install --name git-apiserver --values values.yaml appscode/git-apiserver
```

## RBAC

By default the chart will not install the recommended RBAC roles and rolebindings.

You need to have the flag `--authorization-mode=RBAC` on the api server. See the following document for how to enable [RBAC](https://kubernetes.io/docs/admin/authorization/rbac/).

To determine if your cluster supports RBAC, run the following command:

```console
$ kubectl api-versions | grep rbac
```

If the output contains "beta", you may install the chart with RBAC enabled (see below).

### Enable RBAC role/rolebinding creation

To enable the creation of RBAC resources (On clusters with RBAC). Do the following:

```console
$ helm install --name git-apiserver appscode/git-apiserver --set rbac.create=true
```
