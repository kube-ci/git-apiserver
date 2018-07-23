package controller

import (
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
)

func (c *Controller) initBindingWatcher() {
	c.bindingInformer = c.gitAPIServerInformerFactory.Git().V1alpha1().Bindings().Informer()
	c.bindingQueue = queue.New("Binding", c.MaxNumRequeues, c.NumThreads, c.runBindingInjector)
	c.bindingInformer.AddEventHandler(queue.DefaultEventHandler(c.bindingQueue.GetQueue()))
	c.bindingLister = c.gitAPIServerInformerFactory.Git().V1alpha1().Bindings().Lister()
}

func (c *Controller) runBindingInjector(key string) error {
	obj, exist, err := c.bindingInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		glog.Warningf("Binding %s does not exist anymore\n", key)
	} else {
		binding := obj.(*api.Binding).DeepCopy()
		if binding.Status.LastObservedGeneration == nil || binding.Generation > *binding.Status.LastObservedGeneration {
			glog.Infof("Sync/Add/Update for Binding %s\n", key)
			// do tasks...
			c.updateBindingLastObservedGen(binding.Name, binding.Namespace, binding.Generation) // TODO: errors
		}
	}
	return nil
}

func (c *Controller) reconcileForBinding(binding *api.Binding) error {
	return nil
}
