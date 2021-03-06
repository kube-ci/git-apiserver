/*
Copyright 2019 The KubeCI Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package versioned

import (
	gitv1alpha1 "github.com/kube-ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1"
	webhooksv1alpha1 "github.com/kube-ci/git-apiserver/client/clientset/versioned/typed/webhooks/v1alpha1"
	discovery "k8s.io/client-go/discovery"
	rest "k8s.io/client-go/rest"
	flowcontrol "k8s.io/client-go/util/flowcontrol"
)

type Interface interface {
	Discovery() discovery.DiscoveryInterface
	GitV1alpha1() gitv1alpha1.GitV1alpha1Interface
	WebhooksV1alpha1() webhooksv1alpha1.WebhooksV1alpha1Interface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	*discovery.DiscoveryClient
	gitV1alpha1      *gitv1alpha1.GitV1alpha1Client
	webhooksV1alpha1 *webhooksv1alpha1.WebhooksV1alpha1Client
}

// GitV1alpha1 retrieves the GitV1alpha1Client
func (c *Clientset) GitV1alpha1() gitv1alpha1.GitV1alpha1Interface {
	return c.gitV1alpha1
}

// WebhooksV1alpha1 retrieves the WebhooksV1alpha1Client
func (c *Clientset) WebhooksV1alpha1() webhooksv1alpha1.WebhooksV1alpha1Interface {
	return c.webhooksV1alpha1
}

// Discovery retrieves the DiscoveryClient
func (c *Clientset) Discovery() discovery.DiscoveryInterface {
	if c == nil {
		return nil
	}
	return c.DiscoveryClient
}

// NewForConfig creates a new Clientset for the given config.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c
	if configShallowCopy.RateLimiter == nil && configShallowCopy.QPS > 0 {
		configShallowCopy.RateLimiter = flowcontrol.NewTokenBucketRateLimiter(configShallowCopy.QPS, configShallowCopy.Burst)
	}
	var cs Clientset
	var err error
	cs.gitV1alpha1, err = gitv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	cs.webhooksV1alpha1, err = webhooksv1alpha1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	cs.DiscoveryClient, err = discovery.NewDiscoveryClientForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.gitV1alpha1 = gitv1alpha1.NewForConfigOrDie(c)
	cs.webhooksV1alpha1 = webhooksv1alpha1.NewForConfigOrDie(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClientForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.gitV1alpha1 = gitv1alpha1.New(c)
	cs.webhooksV1alpha1 = webhooksv1alpha1.New(c)

	cs.DiscoveryClient = discovery.NewDiscoveryClient(c)
	return &cs
}
