package controller

import (
	"time"

	"github.com/appscode/go/log"
	"github.com/appscode/go/types"
	"github.com/appscode/kubernetes-webhook-util/admission"
	hooks "github.com/appscode/kubernetes-webhook-util/admission/v1beta1"
	webhook "github.com/appscode/kubernetes-webhook-util/admission/v1beta1/generic"
	"github.com/appscode/kutil/tools/queue"
	kerr "k8s.io/apimachinery/pkg/api/errors"
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
	c.repoInformer.AddEventHandler(queue.DefaultEventHandler(c.repoQueue.GetQueue()))
	c.repoLister = c.gitAPIServerInformerFactory.Git().V1alpha1().Repositories().Lister()
}

func (c *Controller) runRepositoryInjector(key string) error {
	obj, exist, err := c.repoInformer.GetIndexer().GetByKey(key)
	if err != nil {
		log.Errorf("Fetching object with key %s from store failed with %v", key, err)
		return err
	}

	if !exist {
		log.Warningf("Repository %s does not exist anymore\n", key)
	} else {
		repo := obj.(*api.Repository).DeepCopy()

		// TODO: periodically reconcile or use a node-watcher ?
		// don't use LastObservedGeneration, always reconcile repository
		// it will help us to check binding is valid or not periodically

		log.Infof("Sync/Add/Update for Repository %s\n", key)
		if err := c.reconcileForRepository(repo); err != nil {
			return err
		}

		/*if repo.Status.LastObservedGeneration == nil || repo.Generation > *repo.Status.LastObservedGeneration {
			log.Infof("Sync/Add/Update for Repository %s\n", key)
			if err := c.reconcileForRepository(repo); err != nil {
				return err
			}
			// update LastObservedGeneration // TODO: errors ?
			c.updateRepositoryLastObservedGen(repo.Name, repo.Namespace, repo.Generation)
		}*/
	}
	return nil
}

func (c *Controller) reconcileForRepository(repository *api.Repository) error {
	// fetch all open prs initially
	if repository.Spec.Host == "github" {
		log.Infof("Syncing github PRs for repository %s", repository.Name)
		err := c.initGithubPRs(repository)
		if err != nil {
			return err
		}
	}

	go func() {
		for {
			// TODO: write error events to repository
			// if repository not found, we should stop the git watcher
			if err := c.runOnce(repository); kerr.IsNotFound(err) {
				log.Errorf("Stopping git watcher for repository %s/%s, reason: %s", repository.Namespace, repository.Name, err)
				break
			} else if err != nil {
				log.Errorln(err)
			}
			time.Sleep(time.Second * 30) // TODO: period ?
		}
	}()

	return nil
}

func (c *Controller) runOnce(repository *api.Repository) error {
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
