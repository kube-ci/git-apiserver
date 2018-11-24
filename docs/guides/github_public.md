> New to Git API server? Please start [here](/docs/concepts/README.md).

# Sync Public Repository

This tutorial will show you how to use Git API server to sync a public repository hosted in Github.

Before we start, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). Now, install Git API server in your cluster following the steps [here](/docs/setup/install.md). Also, configure Github webhook following the steps [here](/docs/guides/webhook.md).

## Create Repository CRD

Now, create the repository custom resource by specifying clone URL.

```console
$ kubectl apply -f docs/examples/repository.yaml
repository.git.kube.ci/kubeci-gpig created
```

```yaml
apiVersion: git.kube.ci/v1alpha1
kind: Repository
metadata:
  name: kubeci-gpig
  namespace: default
spec:
  host: github
  owner: diptadas
  repo: kubeci-gpig
  cloneUrl: https://github.com/diptadas/kubeci-gpig.git
```

## Get Synced Resources

List all resources associated with the repository:

```console
$ kubectl get all -l repository=kubeci-gpig
NAME                                 AGE
tag.git.kube.ci/kubeci-gpig-v0.0.1   7s

NAME                                    AGE
pullrequest.git.kube.ci/kubeci-gpig-1   39s

NAME                                    AGE
branch.git.kube.ci/kubeci-gpig-b001     8s
branch.git.kube.ci/kubeci-gpig-master   7s
```

List of branches:

```console
$ kubectl get branches -l repository=kubeci-gpig
NAME                 AGE
kubeci-gpig-b001     41s
kubeci-gpig-master   41s
```

List of tags:

```yaml
$ kubectl get tags -l repository=kubeci-gpig
NAME                 AGE
kubeci-gpig-v0.0.1   39s
```

List of open pull requests:

```console
$ kubectl get pullrequests -l repository=kubeci-gpig,state=open
NAME            AGE
kubeci-gpig-1   7m
```

Get specific resource:

```yaml
$ kubectl get branch kubeci-gpig-master -o yaml
apiVersion: git.kube.ci/v1alpha1
kind: Branch
metadata:
  creationTimestamp: 2018-11-13T04:00:54Z
  generation: 1
  labels:
    repository: kubeci-gpig
  name: kubeci-gpig-master
  namespace: default
  ownerReferences:
  - apiVersion: git.kube.ci/v1alpha1
    blockOwnerDeletion: true
    kind: Repository
    name: kubeci-gpig
    uid: a3b4ad46-e6f8-11e8-bd7b-080027c0efdb
  resourceVersion: "11542"
  selfLink: /apis/git.kube.ci/v1alpha1/namespaces/default/branches/kubeci-gpig-master
  uid: b7de0815-e6f8-11e8-bd7b-080027c0efdb
spec:
  lastCommitHash: ef96193e5bb9b3d95e859300670a19f0de38ed7f
```

## Cleanup

```console
$ kubectl delete -f docs/examples/repository.yaml
repository.git.kube.ci "kubeci-gpig" deleted
```
