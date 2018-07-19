package controller

import (
	"log"

	"github.com/appscode/kubernetes-webhook-util/admission"
	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	webhook "github.com/appscode/kubernetes-webhook-util/admission/v1beta1/generic"
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kube.ci/git-apiserver/apis/git"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
)

func (c *Controller) NewRepositoryWebhook() hooks.AdmissionHook {
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
func (c *Controller) initRepositoryWatcher() {
	c.repoInformer = c.gitAPIServerInformerFactory.Git().V1alpha1().Repositories().Informer()
	c.repoQueue = queue.New("Repository", c.MaxNumRequeues, c.NumThreads, c.runRepositoryInjector)
	c.repoInformer.AddEventHandler(queue.DefaultEventHandler(c.repoQueue.GetQueue()))
	c.repoLister = c.gitAPIServerInformerFactory.Git().V1alpha1().Repositories().Lister()
}

func (c *Controller) runRepositoryInjector(key string) error {
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
		log.Println(repo)

		// finalizer
		/*if repo.DeletionTimestamp != nil {
			if core_util.HasFinalizer(repo.ObjectMeta, util.RepositoryFinalizer) {
				err = c.deleteRepository(repo)
				if err != nil {
					return err
				}
				_, _, err = git_apiserver_util.PatchRepository(c.gitAPIServerClient.GitV1alpha1(), repo, func(in *api.Repository) *api.Repository {
					in.ObjectMeta = core_util.RemoveFinalizer(in.ObjectMeta, util.RepositoryFinalizer)
					return in
				})
				return err
			}
		} else {
			_, _, err = git_apiserver_util.PatchRepository(c.gitAPIServerClient.GitV1alpha1(), repo, func(in *api.Repository) *api.Repository {
				in.ObjectMeta = core_util.AddFinalizer(in.ObjectMeta, util.RepositoryFinalizer)
				return in
			})
			return err
		}*/
	}
	return nil
}

func (c *Controller) deleteRepository(repository *api.Repository) error {
	return nil
}
