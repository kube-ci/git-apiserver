package controller

import (
	"fmt"

	"github.com/appscode/go/log"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/tools/queue"
	api "github.com/kube-ci/git-apiserver/apis/git/v1alpha1"
	cs "github.com/kube-ci/git-apiserver/client/clientset/versioned"
	git_apiserver_informers "github.com/kube-ci/git-apiserver/client/informers/externalversions"
	git_apiserver_listers "github.com/kube-ci/git-apiserver/client/listers/git/v1alpha1"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
)

type Controller struct {
	config

	kubeClient         kubernetes.Interface
	gitAPIServerClient cs.Interface
	crdClient          crd_cs.ApiextensionsV1beta1Interface
	recorder           record.EventRecorder

	kubeInformerFactory         informers.SharedInformerFactory
	gitAPIServerInformerFactory git_apiserver_informers.SharedInformerFactory

	// Repository
	repoQueue        *queue.Worker
	repoInformer     cache.SharedIndexInformer
	repoLister       git_apiserver_listers.RepositoryLister
	repoSyncChannels map[string]chan struct{}
}

func (c *Controller) ensureCustomResourceDefinitions() error {
	crds := []*crd_api.CustomResourceDefinition{
		api.Repository{}.CustomResourceDefinition(),
		api.Branch{}.CustomResourceDefinition(),
		api.Tag{}.CustomResourceDefinition(),
		api.PullRequest{}.CustomResourceDefinition(),
	}
	return crdutils.RegisterCRDs(c.crdClient, crds)
}

func (c *Controller) RunInformers(stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	log.Info("Starting git apiserver")
	c.kubeInformerFactory.Start(stopCh)
	c.gitAPIServerInformerFactory.Start(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	for _, v := range c.kubeInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}
	for _, v := range c.gitAPIServerInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	c.repoQueue.Run(stopCh)
}
