package controller

import (
	"github.com/appscode/kubernetes-webhook-util/admission"
	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	webhook "github.com/appscode/kubernetes-webhook-util/admission/v1beta1/generic"
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kube.ci/git-apiserver/apis/git"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
	"kube.ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1/util"
)

const (
	NodeLabelKey = "node"
	NodeMinikube = "minikube"
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
		return c.reconcileBindingForRepository(repo)

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

/*func (c *Controller) deleteRepository(repository *api.Repository) error {
	return nil
}*/

func (c *Controller) reconcileBindingForRepository(repository *api.Repository) error {
	/*binding, err := c.getBinding(repository.Name, repository.Namespace)
	if err != nil {
		return err
	}

	if ok, err := c.isBindingValid(binding); err != nil {
		return err
	} else if ok { // binding exists and valid, no nothing
		return nil
	} else { // binding do not exists or invalid, create/patch binding

	}*/

	meta := metav1.ObjectMeta{
		Name:      repository.Name,
		Namespace: repository.Namespace,
	}

	transform := func(binding *api.Binding) *api.Binding {
		if ok, _ := c.isBindingValid(binding); !ok {
			if binding.Labels == nil {
				binding.Labels = make(map[string]string, 0)
			}
			binding.Labels[NodeLabelKey] = c.nextNodeName()
		}
		binding.Status.LastSynced = metav1.Now()
		return binding
	}

	_, _, err := util.CreateOrPatchBinding(c.gitAPIServerClient.GitV1alpha1(), meta, transform)

	return err
}

/*func (c *Controller) getBinding(name, namespace string) (*api.Binding, error) {
	bindings, err := c.bindingLister.List(labels.Everything())
	if err != nil {
		return nil, err
	}
	for _, binding := range bindings {
		if binding.Name == name && binding.Namespace == namespace {
			return binding, nil
		}
	}
	return nil, nil
}*/

func (c *Controller) isBindingValid(binding *api.Binding) (bool, error) {
	if binding == nil || binding.Labels == nil || binding.Labels[NodeLabelKey] == "" {
		return false, nil
	}

	// check node exists or not
	_, err := c.kubeClient.CoreV1().Nodes().Get(binding.Labels[NodeLabelKey], metav1.GetOptions{})
	if err != nil {
		if kerr.IsNotFound(err) { // node not found
			return false, nil
		}
		return false, err // something wrong
	}

	return true, nil
}

func (c *Controller) nextNodeName() string { // TODO: use a node selector strategy
	return NodeMinikube
}
