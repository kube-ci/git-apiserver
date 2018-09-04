package controller

import (
	"time"

	. "github.com/appscode/go/encoding/json/types"
	"github.com/appscode/go/log"
	"github.com/appscode/go/types"
	"github.com/appscode/kubernetes-webhook-util/admission"
	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	webhook "github.com/appscode/kubernetes-webhook-util/admission/v1beta1/generic"
	meta_util "github.com/appscode/kutil/meta"
	"github.com/appscode/kutil/tools/queue"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kube.ci/git-apiserver/apis/git"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
	"kube.ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1/util"
	"kube.ci/git-apiserver/pkg/git-repo"
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
	c.repoInformer.AddEventHandler(queue.NewObservableHandler(c.repoQueue.GetQueue(), api.EnableStatusSubresource))
	c.repoLister = c.gitAPIServerInformerFactory.Git().V1alpha1().Repositories().Lister()
	c.repoSyncChannels = make(map[string]chan struct{})
}

func (c *Controller) runRepositoryInjector(key string) error {
	obj, exist, err := c.repoInformer.GetIndexer().GetByKey(key)
	if err != nil {
		log.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		log.Warningf("Repository %s does not exist anymore\n", key)
		if stopCh, ok := c.repoSyncChannels[key]; ok { // send stop signal and delete from map
			log.Infof("Closing sync for repository %s", key)
			close(stopCh)
			delete(c.repoSyncChannels, key)
		}
	} else {
		repo := obj.(*api.Repository).DeepCopy()
		log.Infof("Sync/Add/Update for Repository %s\n", key)

		if stopCh, ok := c.repoSyncChannels[key]; ok { // send stop signal
			log.Infof("Restarting sync for repository %s", key)
			close(stopCh)
		}

		c.repoSyncChannels[key] = make(chan struct{}) // create new stop channel
		if err := c.reconcileForRepository(repo, c.repoSyncChannels[key]); err != nil {
			return err
		}

		// update LastObservedGeneration
		_, err = util.UpdateRepositoryStatus(
			c.gitAPIServerClient.GitV1alpha1(),
			repo.ObjectMeta,
			func(r *api.RepositoryStatus) *api.RepositoryStatus {
				r.ObservedGeneration = NewIntHash(repo.Generation, meta_util.GenerationHash(repo))
				return r
			},
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Controller) reconcileForRepository(repository *api.Repository, stopCh <-chan struct{}) error {
	// fetch and reconcile open prs initially
	// next updates will be handled through webhook
	if repository.Spec.Host == "github" {
		log.Infof("Syncing github PRs for repository %s/%s", repository.Namespace, repository.Name)
		if err := c.fetchAndReconcileGithubPRs(repository); err != nil {
			return err
		}
	}

	// periodically fetch and reconcile branches, tags
	// stop using stopCh
	go func() {
		t := time.NewTicker(30 * time.Second)
	loop:
		for {
			select {
			case <-stopCh:
				log.Infof("Stop signal received for repository %s/%s", repository.Namespace, repository.Name)
				break loop
			case <-t.C:
				if err := c.fetchAndReconcileRefs(repository); err != nil {
					log.Errorln(err)
				}
			}
		}
	}()

	return nil
}

func (c *Controller) fetchAndReconcileRefs(repository *api.Repository) error {
	log.Infof("Fetching repository %s/%s", repository.Namespace, repository.Name)

	// repository token, empty if repository.Spec.TokenFormSecret is nil
	token, err := repository.GetToken(c.kubeClient)
	if err != nil {
		return err
	}

	repo, err := git_repo.Fetch(repository.Spec.CloneUrl, token)
	if err != nil {
		return err
	}

	log.Infof("Reconciling branches for repository %s/%s", repository.Namespace, repository.Name)
	if err = c.reconcileBranches(repository, repo.Branches); err != nil {
		return err
	}

	log.Infof("Reconciling tags for repository %s/%s", repository.Namespace, repository.Name)
	if err = c.reconcileTags(repository, repo.Tags); err != nil {
		return err
	}

	return nil
}

func (c *Controller) reconcileBranches(repository *api.Repository, branches []git_repo.Reference) error {
	// create or patch branch CRDs
	for _, gitBranch := range branches {
		meta := metav1.ObjectMeta{
			Name:      repository.Name + "-" + gitBranch.Name,
			Namespace: repository.Namespace,
			OwnerReferences: []metav1.OwnerReference{
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
			log.Infof("Deleting Branch %s/%s", branch.Namespace, branch.Name)
			err = c.gitAPIServerClient.GitV1alpha1().Branches(branch.Namespace).Delete(branch.Name, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (c *Controller) reconcileTags(repository *api.Repository, tags []git_repo.Reference) error {
	// create or patch tag CRDs
	for _, gitTag := range tags {
		meta := metav1.ObjectMeta{
			Name:      repository.Name + "-" + gitTag.Name,
			Namespace: repository.Namespace,
			OwnerReferences: []metav1.OwnerReference{
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
			log.Infof("Deleting Tag %s/%s", tag.Namespace, tag.Name)
			err = c.gitAPIServerClient.GitV1alpha1().Tags(tag.Namespace).Delete(tag.Name, nil)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
