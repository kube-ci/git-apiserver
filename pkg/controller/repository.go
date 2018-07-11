package controller

import (
	"github.com/appscode/kubernetes-webhook-util/admission"
	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	webhook "github.com/appscode/kubernetes-webhook-util/admission/v1beta1/generic"
	core_util "github.com/appscode/kutil/core/v1"
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kube.ci/git-apiserver/apis/git"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
	kubeci_util "kube.ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1/util"
	"kube.ci/git-apiserver/pkg/util"
)

func (c *StashController) NewRepositoryWebhook() hooks.AdmissionHook {
	return webhook.NewGenericWebhook(
		schema.GroupVersionResource{
			Group:    "admission.git.kube.ci",
			Version:  "v1alpha1",
			Resource: "repositories",
		},
		"repository",
		[]string{git.GroupName},
		api.SchemeGroupVersion.WithKind("Repository"),
		nil,
		&admission.ResourceHandlerFuncs{
			CreateFunc: func(obj runtime.Object) (runtime.Object, error) {
				return nil, obj.(*api.Repository).IsValid()
			},
			UpdateFunc: func(oldObj, newObj runtime.Object) (runtime.Object, error) {
				return nil, newObj.(*api.Repository).IsValid()
			},
		},
	)
}
func (c *StashController) initRepositoryWatcher() {
	c.repoInformer = c.kubeciInformerFactory.Git().V1alpha1().Repositories().Informer()
	c.repoQueue = queue.New("Repository", c.MaxNumRequeues, c.NumThreads, c.runRepositoryInjector)
	c.repoInformer.AddEventHandler(queue.DefaultEventHandler(c.repoQueue.GetQueue()))
	c.repoLister = c.kubeciInformerFactory.Git().V1alpha1().Repositories().Lister()
}

func (c *StashController) runRepositoryInjector(key string) error {
	obj, exist, err := c.repoInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		glog.Warningf("Repository %s does not exist anymore\n", key)
	} else {
		glog.Infof("Sync/Add/Update for Repository %s\n", key)

		repo := obj.(*api.Repository)

		if repo.DeletionTimestamp != nil {
			if core_util.HasFinalizer(repo.ObjectMeta, util.RepositoryFinalizer) {
				if repo.Spec.WipeOut {
					err = c.deleteResticRepository(repo)
					if err != nil {
						return err
					}
				}
				_, _, err = kubeci_util.PatchRepository(c.kubeciClient.GitV1alpha1(), repo, func(in *api.Repository) *api.Repository {
					in.ObjectMeta = core_util.RemoveFinalizer(in.ObjectMeta, util.RepositoryFinalizer)
					return in
				})
				return err
			}
		} else {
			_, _, err = kubeci_util.PatchRepository(c.kubeciClient.GitV1alpha1(), repo, func(in *api.Repository) *api.Repository {
				in.ObjectMeta = core_util.AddFinalizer(in.ObjectMeta, util.RepositoryFinalizer)
				return in
			})
			return err
		}
	}
	return nil
}

func (c *StashController) deleteResticRepository(repository *api.Repository) error {
	return nil
}
