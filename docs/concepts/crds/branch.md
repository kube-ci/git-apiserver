# Branch

## What is Branch

A `Branch` is a Kubernetes `CustomResourceDefinition` (CRD) that represents a Git branch. Branches are synced(create/update/delete) with remote for existing repository CRDs. The naming format is `{repository-CRD-name}-{branch-name}`.

## Branch Spec

As with all other Kubernetes objects, a Branch needs `apiVersion`, `kind`, and `metadata` fields. It also needs a `.spec` section. Below is an example Branch object:

```yaml
apiVersion: git.kube.ci/v1alpha1
kind: Branch
metadata:
  creationTimestamp: 2018-11-09T06:51:33Z
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
    uid: d0a491e9-e3eb-11e8-a7e0-080027868e9e
  resourceVersion: "111023"
  selfLink: /apis/git.kube.ci/v1alpha1/namespaces/default/branches/kubeci-gpig-master
  uid: e4e53e0e-e3eb-11e8-a7e0-080027868e9e
spec:
  lastCommitHash: ef96193e5bb9b3d95e859300670a19f0de38ed7f
```

The `.spec` section has following parts:

### spec.lastCommitHash

The latest commit hash associated with this Git branch.