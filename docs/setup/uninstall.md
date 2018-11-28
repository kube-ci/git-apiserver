# Uninstall Git API server

To uninstall Git API server, run the following command:

```console
$ curl -fsSL https://raw.githubusercontent.com/kube-ci/git-apiserver/0.1.0/hack/deploy/install.sh \
    | bash -s -- --uninstall [--namespace=NAMESPACE]

validatingwebhookconfiguration.admissionregistration.k8s.io "admission.git.kube.ci" deleted
No resources found
apiservice.apiregistration.k8s.io "v1alpha1.admission.git.kube.ci" deleted
apiservice.apiregistration.k8s.io "v1alpha1.webhooks.git.kube.ci" deleted
deployment.extensions "git-apiserver" deleted
service "git-apiserver" deleted
secret "git-apiserver-cert" deleted
serviceaccount "git-apiserver" deleted
clusterrolebinding.rbac.authorization.k8s.io "git-apiserver" deleted
clusterrolebinding.rbac.authorization.k8s.io "git-apiserver-auth-delegator" deleted
clusterrole.rbac.authorization.k8s.io "git-apiserver" deleted
rolebinding.rbac.authorization.k8s.io "git-apiserver-extension-server-authentication-reader" deleted
No resources found
waiting for git-apiserver operator pod to stop running

Successfully uninstalled GIT-APISERVER!
```

The above command will leave the Git API server CRD objects as-is. If you wish to **nuke** all Git API server CRD objects, also pass the `--purge` flag. This will keep a copy of Git API server CRD objects in your current directory.