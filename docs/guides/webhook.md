> New to Git API server? Please start [here](/docs/concepts/README.md).

# Configure Github Webhook

This tutorial will show you how to configure Github webhook to sync pull-request events.

Before we start, you need to have a Kubernetes cluster, and the kubectl command-line tool must be configured to communicate with your cluster. If you do not already have a cluster, you can create one by using [Minikube](https://github.com/kubernetes/minikube). Now, install Git API server in your cluster following the steps [here](/docs/setup/install.md).

First, go to `Webhooks` from settings tab of your Github repository. Now, select `Add webhook` to create a new webhook for the repository. Then, set the followings:

- Payload URL: `https://{master-ip}/apis/webhook.git.kube.ci/v1alpha1/githubevents`
- Content Type: `application/json`
- SSL verification: Disable

Also, select `Pull Requests` from individual events.