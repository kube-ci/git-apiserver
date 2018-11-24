# Installation Guide

Git API server can be installed via a script or as a Helm chart.

<ul class="nav nav-tabs" id="installerTab" role="tablist">
  <li class="nav-item">
    <a class="nav-link active" id="script-tab" data-toggle="tab" href="#script" role="tab" aria-controls="script" aria-selected="true">Script</a>
  </li>
  <li class="nav-item">
    <a class="nav-link" id="helm-tab" data-toggle="tab" href="#helm" role="tab" aria-controls="helm" aria-selected="false">Helm</a>
  </li>
</ul>
<div class="tab-content" id="installerTabContent">
  <div class="tab-pane fade show active" id="script" role="tabpanel" aria-labelledby="script-tab">

## Using Script

To install Git API server in your Kubernetes cluster, run the following command:

```console
$ curl -fsSL https://raw.githubusercontent.com/kube-ci/git-apiserver/0.1.0/hack/deploy/install.sh | bash
```

After successful installation, you should have a `git-apiserver-***` pod running in the `kube-system` namespace.

```console
$ kubectl get pods -n kube-system | grep git-apiserver
git-apiserver-846d47f489-jrb58       1/1       Running   0          48s
```

#### Customizing Installer

The installer script and associated yaml files can be found in the [/hack/deploy](https://github.com/kube-ci/git-apiserver/tree/0.1.0/hack/deploy) folder. You can see the full list of flags available to installer using `-h` flag.

```console
$ curl -fsSL https://raw.githubusercontent.com/kube-ci/git-apiserver/0.1.0/hack/deploy/install.sh | bash -s -- -h
git-apiserver.sh - install git-apiserver operator

git-apiserver.sh [options]

options:
-h, --help                         show brief help
-n, --namespace=NAMESPACE          specify namespace (default: kube-system)
    --rbac                         create RBAC roles and bindings (default: true)
    --docker-registry              docker registry used to pull git-apiserver images (default: kubeci)
    --image-pull-secret            name of secret used to pull git-apiserver operator images
    --run-on-master                run git-apiserver operator on master
    --enable-validating-webhook    enable/disable validating webhooks for git-apiserver crds
    --enable-mutating-webhook      enable/disable mutating webhooks for Kubernetes workloads
    --enable-status-subresource    If enabled, uses status sub resource for crds
    --enable-analytics             send usage events to Google Analytics (default: true)
    --uninstall                    uninstall git-apiserver
    --purge                        purges git-apiserver crd objects and crds
```

If you would like to run Git API server operator pod in `master` instances, pass the `--run-on-master` flag:

```console
$ curl -fsSL https://raw.githubusercontent.com/kube-ci/git-apiserver/0.1.0/hack/deploy/install.sh \
    | bash -s -- --run-on-master [--rbac]
```

Git API server operator will be installed in a `kube-system` namespace by default. If you would like to run operator pod in `git-apiserver` namespace, pass the `--namespace=git-apiserver` flag:

```console
$ kubectl create namespace git-apiserver
$ curl -fsSL https://raw.githubusercontent.com/kube-ci/git-apiserver/0.1.0/hack/deploy/install.sh \
    | bash -s -- --namespace=git-apiserver [--run-on-master] [--rbac]
```

If you are using a private Docker registry, you need to pull the following image:

 - [kubeci/git-apiserver](https://hub.docker.com/r/kubeci/git-apiserver)

To pass the address of your private registry and optionally a image pull secret use flags `--docker-registry` and `--image-pull-secret` respectively.

```console
$ kubectl create namespace git-apiserver
$ curl -fsSL https://raw.githubusercontent.com/kube-ci/git-apiserver/0.1.0/hack/deploy/install.sh \
    | bash -s -- --docker-registry=MY_REGISTRY [--image-pull-secret=SECRET_NAME] [--rbac]
```

Git API server implements [validating admission webhooks](https://kubernetes.io/docs/admin/admission-controllers/#validatingadmissionwebhook-alpha-in-18-beta-in-19) to validate Git API server CRDs. This is enabled by default for Kubernetes 1.9.0 or later releases. To disable this feature, pass the `--enable-validating-webhook=false` flag.

```console
$ curl -fsSL https://raw.githubusercontent.com/kube-ci/git-apiserver/0.1.0/hack/deploy/install.sh \
    | bash -s -- --enable-validating-webhook=false [--rbac]
```

Git API server 0.1.0 or later releases can use status sub resource for CustomResourceDefinitions. This is enabled by default for Kubernetes 1.11.0 or later releases. To disable this feature, pass the `--enable-status-subresource=false` flag.

</div>
<div class="tab-pane fade" id="helm" role="tabpanel" aria-labelledby="helm-tab">

## Using Helm
Git API server can be installed via [Helm](https://helm.sh/) using the [chart](https://github.com/kube-ci/git-apiserver/tree/0.1.0/chart/git-apiserver) from [AppsCode Charts Repository](https://github.com/appscode/charts). To install the chart with the release name `my-release`:

```console
$ helm repo add appscode https://charts.appscode.com/stable/
$ helm repo update
$ helm search appscode/git-apiserver
NAME            CHART VERSION APP VERSION DESCRIPTION
appscode/git-apiserver  0.1.0    0.1.0  git-apiserver by AppsCode - Kuberenetes native CI system

$ helm install appscode/git-apiserver --name my-release --version 0.1.0 --namespace kube-system
```

To see the detailed configuration options, visit [here](https://github.com/kube-ci/git-apiserver/tree/master/chart/git-apiserver).

</div>

### Installing in GKE Cluster

If you are installing Git API server on a GKE cluster, you will need cluster admin permissions to install Git API server operator. Run the following command to grant admin permission to the cluster.

```console
$ kubectl create clusterrolebinding "cluster-admin-$(whoami)" \
  --clusterrole=cluster-admin \
  --user="$(gcloud config get-value core/account)"
```


## Verify installation

To check if Git API server operator pods have started, run the following command:
```console
$ kubectl get pods --all-namespaces -l app=git-apiserver --watch

NAMESPACE     NAME                              READY     STATUS    RESTARTS   AGE
kube-system   git-apiserver-859d6bdb56-m9br5    2/2       Running   2          5s
```

Once the operator pods are running, you can cancel the above command by typing `Ctrl+C`.

Now, to confirm CRD groups have been registered by the operator, run the following command:
```console
$ kubectl get crd -l app=git-apiserver

NAME                               CREATED AT
branches.git.kube.ci               2018-11-09T06:49:02Z
pullrequests.git.kube.ci           2018-11-09T06:49:03Z
repositories.git.kube.ci           2018-11-09T06:49:01Z
tags.git.kube.ci                   2018-11-09T06:49:02Z
```

Now, you are ready to [run your first workflow](/docs/guides/README.md) using git-apiserver.

## Configuring RBAC

Git API server introduces resources, such as, `Branch`, `Tag`, `PullRequest`. Git API server installer will create 2 user facing cluster roles:

| ClusterRole                 | Aggregates To | Description                            |
|-----------------------------|---------------|----------------------------------------|
| appscode:git-apiserver:edit | admin, edit   | Allows edit access to git-apiserver CRDs, intended to be granted within a namespace using a RoleBinding. |
| appscode:git-apiserver:view | view          | Allows read-only access to git-apiserver CRDs, intended to be granted within a namespace using a RoleBinding. |

These user facing roles supports [ClusterRole Aggregation](https://kubernetes.io/docs/admin/authorization/rbac/#aggregated-clusterroles) feature in Kubernetes 1.9 or later clusters.


## Using kubectl for Repositories

```console
# List all Repository objects
$ kubectl get repository --all-namespaces

# List Repository objects for a namespace
$ kubectl get repository -n <namespace>

# Get Repository YAML
$ kubectl get repository -n <namespace> <name> -o yaml

# Describe Repository. Very useful to debug problems.
$ kubectl describe repository -n <namespace> <name>
```

## Detect git-apiserver version

To detect git-apiserver version, exec into the operator pod and run `git-apiserver version` command.

```console
$ POD_NAMESPACE=kube-system
$ POD_NAME=$(kubectl get pods -n $POD_NAMESPACE -l app=git-apiserver -o jsonpath={.items[0].metadata.name})
$ kubectl exec -it $POD_NAME -c operator -n $POD_NAMESPACE git-apiserver version

Version = 0.1.0
VersionStrategy = tag
Os = alpine
Arch = amd64
CommitHash = 85b0f16ab1b915633e968aac0ee23f877808ef49
GitBranch = release-0.1.0
GitTag = 0.1.0
CommitTimestamp = 2018-10-10T05:24:23
```