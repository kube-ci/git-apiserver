package controller

import (
	"path/filepath"
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/go/types"
	"github.com/appscode/kutil/tools/queue"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
	"kube.ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1/util"
	"kube.ci/git-apiserver/pkg/git-repo"
)

func (c *Controller) initBindingWatcher() {
	c.bindingInformer = c.gitAPIServerInformerFactory.Git().V1alpha1().Bindings().Informer()
	c.bindingQueue = queue.New("Binding", c.MaxNumRequeues, c.NumThreads, c.runBindingInjector)
	c.bindingLister = c.gitAPIServerInformerFactory.Git().V1alpha1().Bindings().Lister()
	c.bindingInformer.AddEventHandler(
		queue.NewFilteredHandler(
			queue.DefaultEventHandler(c.bindingQueue.GetQueue()),
			labels.SelectorFromSet(map[string]string{
				NodeLabelKey: NodeMinikube, // TODO: get node-name from pod's env variable
			}),
		),
	)

	c.bindingMap = make(map[string]struct{})
}

func (c *Controller) runBindingInjector(key string) error {
	obj, exist, err := c.bindingInformer.GetIndexer().GetByKey(key)
	if err != nil {
		glog.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		glog.Warningf("Binding %s does not exist anymore\n", key)
		delete(c.bindingMap, key)
	} else {
		binding := obj.(*api.Binding).DeepCopy()

		// use key map instead of LastObservedGeneration to check binding already reconciled or not
		// it will help to restart git-watcher when operator is restarted
		if _, ok := c.bindingMap[key]; !ok {
			glog.Infof("Sync/Add/Update for Binding %s\n", key)
			if err = c.reconcileForBinding(binding); err != nil {
				return err
			}
			c.bindingMap[key] = struct{}{}
		}

		/*if binding.Status.LastObservedGeneration == nil || binding.Generation > *binding.Status.LastObservedGeneration {
			glog.Infof("Sync/Add/Update for Binding %s\n", key)
			if err = c.reconcileForBinding(binding); err != nil {
				return err
			}
			// update LastObservedGeneration // TODO: errors ?
			c.updateBindingLastObservedGen(binding.Name, binding.Namespace, binding.Generation)
		}*/
	}
	return nil
}

func (c *Controller) reconcileForBinding(binding *api.Binding) error {
	go func() {
		for {
			// TODO: write error events to binding or repository ?
			// if repository not found, we should stop the git watcher
			// TODO: use a stop channel instead ?
			if err := c.runOnce(binding.Name, binding.Namespace); kerr.IsNotFound(err) {
				log.Errorf("Stopping git watcher, reason: %s", err)
				break
			} else if err != nil {
				log.Errorln(err)
			}
			time.Sleep(time.Second * 30) // TODO: period ?
		}
	}()
	return nil
}

func (c *Controller) runOnce(name, namespace string) error {
	// get repository CRD
	repository, err := c.gitAPIServerClient.GitV1alpha1().Repositories(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	// repository token, empty if repository.Spec.TokenFormSecret is nil
	token, err := repository.GetToken(c.kubeClient)
	if err != nil {
		return err
	}

	path := filepath.Join("/tmp/git-apiserver", repository.Name) // TODO: use constant
	repo := git_repo.New(repository.Spec.CloneUrl, path, token)
	if err := repo.CloneOrFetch(false); err != nil { // TODO: true if repository crd changes
		return err
	}

	log.Infoln("Cloning/Fetching done...", repo)

	if err = c.reconcileBranches(repository, repo); err != nil {
		return err
	}
	if err = c.reconcileTags(repository, repo); err != nil {
		return err
	}

	return nil
}

func (c *Controller) reconcileBranches(repository *api.Repository, repo *git_repo.Repository) error {
	branches, err := repo.GetBranches()
	if err != nil {
		return err
	}

	// create or patch branch CRDs
	for _, gitBranch := range branches {
		meta := metav1.ObjectMeta{
			Name:      repository.Name + "-" + gitBranch.Name,
			Namespace: repository.Namespace,
			OwnerReferences: []metav1.OwnerReference{ // TODO: owner ref repository or binding ?
				{
					APIVersion:         api.SchemeGroupVersion.Group + "/" + api.SchemeGroupVersion.Version,
					Kind:               api.ResourceKindRepository,
					Name:               repository.Name,
					UID:                repository.UID,
					BlockOwnerDeletion: types.TrueP(),
				},
			},
		}

		transform := func(branch *api.Branch) *api.Branch {
			if branch.Labels == nil {
				branch.Labels = make(map[string]string, 0)
			}
			branch.Labels["repository"] = repository.Name
			branch.Spec.LastCommitHash = gitBranch.Hash
			return branch
		}

		_, _, err := util.CreateOrPatchBranch(c.gitAPIServerClient.GitV1alpha1(), meta, transform)
		if err != nil {
			return err
		}
	}

	// delete old branches that don't exist now
	branchList, err := c.gitAPIServerClient.GitV1alpha1().Branches(repository.Namespace).List(
		metav1.ListOptions{
			LabelSelector: labels.FormatLabels(
				map[string]string{
					"repository": repository.Name,
				},
			),
		},
	)
	if err != nil {
		return err
	}

	for _, branch := range branchList.Items {
		found := false
		for _, gitBranch := range branches {
			if branch.Name == repository.Name+"-"+gitBranch.Name {
				found = true
				break
			}
		}
		if !found {
			err = c.gitAPIServerClient.GitV1alpha1().Branches(repository.Namespace).Delete(branch.Name, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Controller) reconcileTags(repository *api.Repository, repo *git_repo.Repository) error {
	tags, err := repo.GetTags()
	if err != nil {
		return err
	}

	// create or patch tag CRDs
	for _, gitTag := range tags {
		meta := metav1.ObjectMeta{
			Name:      repository.Name + "-" + gitTag.Name,
			Namespace: repository.Namespace,
			OwnerReferences: []metav1.OwnerReference{ // TODO: owner ref repository or binding ?
				{
					APIVersion:         api.SchemeGroupVersion.Group + "/" + api.SchemeGroupVersion.Version,
					Kind:               api.ResourceKindRepository,
					Name:               repository.Name,
					UID:                repository.UID,
					BlockOwnerDeletion: types.TrueP(),
				},
			},
		}

		transform := func(tag *api.Tag) *api.Tag {
			if tag.Labels == nil {
				tag.Labels = make(map[string]string, 0)
			}
			tag.Labels["repository"] = repository.Name
			tag.Spec.LastCommitHash = gitTag.Hash
			return tag
		}

		_, _, err := util.CreateOrPatchTag(c.gitAPIServerClient.GitV1alpha1(), meta, transform)
		if err != nil {
			return err
		}
	}

	// delete old tags that don't exist now
	tagList, err := c.gitAPIServerClient.GitV1alpha1().Tags(repository.Namespace).List(
		metav1.ListOptions{
			LabelSelector: labels.FormatLabels(
				map[string]string{
					"repository": repository.Name,
				},
			),
		},
	)
	if err != nil {
		return err
	}

	for _, tag := range tagList.Items {
		found := false
		for _, gitTag := range tags {
			if tag.Name == repository.Name+"-"+gitTag.Name {
				found = true
				break
			}
		}
		if !found {
			err = c.gitAPIServerClient.GitV1alpha1().Tags(repository.Namespace).Delete(tag.Name, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
