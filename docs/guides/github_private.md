> New to Git API server? Please start [here](/docs/concepts/README.md).

# Sync Public Repository

This tutorial will show you how to use Git API server to sync a public repository hosted in Github. 

Before we start, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). Now, install Git API server in your cluster following the steps [here](/docs/setup/install.md). Also, configure Github webhook following the steps [here](/docs/guides/webhook.md).

## Create Secret

For syncing private repository, you have to create a secret with key `token` and value of Github API token.

```yaml
$ kubectl create secret generic github-credential --from-literal=token={github-api-token}
secret/github-credential created
```

## Create Repository CRD

Now, create the repository custom resource by specifying clone URL and secret.

```console
$ kubectl apply -f docs/examples/repository-private.yaml 
repository.git.kube.ci/private-test-repo created
```

```yaml
apiVersion: git.kube.ci/v1alpha1
kind: Repository
metadata:
  name: private-test-repo
  namespace: default
spec:
  host: github
  owner: tamalsaha
  repo: private-test-repo
  cloneUrl: https://github.com/tamalsaha/private-test-repo.git
  tokenFormSecret: github-credential
```

## Get Synced Resources

```console
$ kubectl get all -l repository=private-test-repo
NAME                                          AGE
pullrequest.git.kube.ci/private-test-repo-1   16m

NAME                                          AGE
branch.git.kube.ci/private-test-repo-b001     16m
branch.git.kube.ci/private-test-repo-master   16m
```

## Cleanup

```console
$ kubectl delete -f docs/examples/repository-private.yaml
repository.git.kube.ci "private-test-repo" deleted
```
