package controller

import (
	"fmt"

	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
	cs "kube.ci/git-apiserver/client/clientset/versioned"
	kubeciinformers "kube.ci/git-apiserver/client/informers/externalversions"
	kubeci_listers "kube.ci/git-apiserver/client/listers/git/v1alpha1"
)

type StashController struct {
	config

	kubeClient   kubernetes.Interface
	kubeciClient cs.Interface
	crdClient    crd_cs.ApiextensionsV1beta1Interface
	recorder     record.EventRecorder

	kubeInformerFactory   informers.SharedInformerFactory
	kubeciInformerFactory kubeciinformers.SharedInformerFactory

	// Repository
	repoQueue    *queue.Worker
	repoInformer cache.SharedIndexInformer
	repoLister   kubeci_listers.RepositoryLister
}

func (c *StashController) ensureCustomResourceDefinitions() error {
	crds := []*crd_api.CustomResourceDefinition{
		api.Repository{}.CustomResourceDefinition(),
	}
	return crdutils.RegisterCRDs(c.crdClient, crds)
}

func (c *StashController) RunInformers(stopCh <-chan struct{}) {
	defer runtime.HandleCrash()

	glog.Info("Starting Stash controller")
	c.kubeInformerFactory.Start(stopCh)
	c.kubeciInformerFactory.Start(stopCh)

	// Wait for all involved caches to be synced, before processing items from the queue is started
	for _, v := range c.kubeInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}
	for _, v := range c.kubeciInformerFactory.WaitForCacheSync(stopCh) {
		if !v {
			runtime.HandleError(fmt.Errorf("timed out waiting for caches to sync"))
			return
		}
	}

	c.repoQueue.Run(stopCh)
}
