# Repository

## What is Repository

A `Repository` is a Kubernetes `CustomResourceDefinition` (CRD). It provides configuration for a remote Git repository. When a Repository CRD is created, all branches and tags are synced periodically with the remote of that repository. Pull requests are also synced based on webhook events.

## Repository Spec

As with all other Kubernetes objects, a Repository needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example Repository object:

```yaml
apiVersion: git.kube.ci/v1alpha1
kind: Repository
metadata:
  name: kubeci-gpig
  namespace: default
spec:
  host: github
  owner: tamalsaha
  repo: kubeci-gpig
  cloneUrl: https://github.com/kube-ci/kubeci-gpig.git
```

The `.spec` section has following parts:

### spec.host

Remote host for your Git repository. Currently, only `github` is supported.

### spec.owner

Owner of your Github repository, required for calling Github [API](https://developer.github.com/v3/pulls/#list-pull-requests) to list pull-requests.

### spec.repo

Name of your Github repository, required for calling Github [API](https://developer.github.com/v3/pulls/#list-pull-requests) to list pull-requests.

### spec.cloneUrl

Clone URL of your remote git repository to fetch branches and tags using `$ git ls-remote` command.

### spec.tokenFormSecret

Name of the Kubernetes secret in the same namespace containing token for accessing private repository under key `token`.

Example secret:

```yaml
apiVersion: v1
kind: Secret
metadata:
  creationTimestamp: 2018-11-30T05:42:24Z
  name: github-credential
  namespace: default
  resourceVersion: "866"
  selfLink: /api/v1/namespaces/default/secrets/github-credential
  uid: b67acf60-f462-11e8-b31f-0800271b8897
type: Opaque
data:
  token: {API_TOKEN}
```