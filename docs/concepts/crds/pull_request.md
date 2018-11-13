# PullRequest

## What is PullRequest

A `PullRequest` is a Kubernetes `CustomResourceDefinition` (CRD) that represents pull request in your Git repository host. Initially when a Repository CRD is created, all pull requests are fetched from remote. After that, they are created/updated based on [pull request events](https://developer.github.com/v3/activity/events/types/#pullrequestevent) with matching clone URL. The naming format is `{repository-CRD-name}-{pull-request-id}`.

## PullRequest Spec

As with all other Kubernetes objects, a PullRequest needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example PullRequest object:

```yaml
apiVersion: git.kube.ci/v1alpha1
kind: PullRequest
metadata:
  creationTimestamp: 2018-11-09T06:51:00Z
  generation: 1
  labels:
    repository: kubeci-gpig
    state: open
    ok-to-tesk:
  name: kubeci-gpig-1
  namespace: default
  ownerReferences:
  - apiVersion: git.kube.ci/v1alpha1
    blockOwnerDeletion: true
    kind: Repository
    name: kubeci-gpig
    uid: d0a491e9-e3eb-11e8-a7e0-080027868e9e
  resourceVersion: "110983"
  selfLink: /apis/git.kube.ci/v1alpha1/namespaces/default/pullrequests/kubeci-gpig-1
  uid: d1975ba0-e3eb-11e8-a7e0-080027868e9e
spec:
  headRef: b001
  headSHA: 46859a78afa9f895962ccf111c7982f66e9d3b72
  number: 1
```

The `repository` and `state` labels are set by default. Apart from them, all existing labels in the remote associated with the pull request are also set with empty value (`ok-to-test` label in the example above).

The `.spec` section has following parts:

### spec.number

Represents the ID number of the pull request.

### spec.headRef

Represents the head reference associated with this pull request.

### spec.headSHA

Represents the SHA of the head (i.e. the latest commit hash) associated with this pull request.

