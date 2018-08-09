package controller

import (
	"context"
	"fmt"

	"github.com/TamalSaha/go-oneliners"
	"github.com/appscode/go/log"
	"github.com/appscode/go/types"
	"github.com/google/go-github/github"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
	repo_v1alpha1 "kube.ci/git-apiserver/apis/git/v1alpha1"
	"kube.ci/git-apiserver/apis/webhook/v1alpha1"
	"kube.ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1/util"
)

type GithubREST struct {
	controller *Controller
}

var _ rest.Creater = &GithubREST{}
var _ rest.Scoper = &GithubREST{}

func NewGithubREST(controller *Controller) *GithubREST {
	return &GithubREST{
		controller: controller,
	}
}

func (r *GithubREST) New() runtime.Object {
	return &v1alpha1.GithubEvent{}
}

func (r *GithubREST) NamespaceScoped() bool {
	return false
}

func (r *GithubREST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return v1alpha1.SchemeGroupVersion.WithKind(v1alpha1.ResourceKindGithubEvent)
}

// curl -k -H 'Content-Type: application/json' -d '{"action":"labeled"}' https://192.168.99.100:8443/apis/webhook.git.kube.ci/v1alpha1/githubpullrequests
func (r *GithubREST) Create(ctx context.Context, obj runtime.Object, createValidation rest.ValidateObjectFunc, includeUninitialized bool) (runtime.Object, error) {
	event := obj.(*v1alpha1.GithubEvent)
	oneliners.PrettyJson(event, "Github Webhook Event")
	r.controller.githubEventHandler(event)
	return event, nil // TODO: error ?
}

func (c *Controller) githubEventHandler(event *v1alpha1.GithubEvent) {
	repositories, err := c.repoLister.List(labels.Everything())
	if err != nil {
		log.Errorln(err)
	}

	// find matching repository
	for _, repository := range repositories {
		if event.Repo != nil && event.Repo.CloneURL != nil && repository.Spec.Url == *event.Repo.CloneURL {
			log.Infof("Event for repository %s/%s", repository.Namespace, repository.Name)
			if event.PullRequest != nil {
				err := c.githubPRHandler(event.PullRequest, repository) // TODO: errors ?
				if err != nil {
					log.Errorln(err)
				}
			}
		}
	}
}

func (c *Controller) githubPRHandler(githubPR *github.PullRequest, repository *repo_v1alpha1.Repository) error {
	// create or patch PR CRD
	meta := metav1.ObjectMeta{
		Name:      fmt.Sprintf("%s-%d", repository.Name, *githubPR.Number),
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

	transform := func(pr *api.PullRequest) *api.PullRequest {
		//if pr.Labels == nil {
		//	pr.Labels = make(map[string]string, 0)
		//}
		// TODO: always create new ?
		pr.Labels = make(map[string]string, 0)
		pr.Labels["repository"] = repository.Name

		// add PR labels
		for _, label := range githubPR.Labels {
			if label != nil && label.Name != nil {
				pr.Labels[*label.Name] = ""
			}
		}
		// add state as label
		if githubPR.State != nil {
			pr.Labels["state"] = *githubPR.State
		}

		if githubPR.Head != nil {
			if githubPR.Head.Ref != nil {
				pr.Spec.HeadRef = *githubPR.Head.Ref
			}
			if githubPR.Head.SHA != nil {
				pr.Spec.HeadSHA = *githubPR.Head.SHA
			}
		}

		if githubPR.ID != nil {
			pr.Spec.Number = *githubPR.Number
		}

		return pr
	}

	_, _, err := util.CreateOrPatchPullRequest(c.gitAPIServerClient.GitV1alpha1(), meta, transform)
	if err != nil {
		return err
	}

	return nil
}

func (c *Controller) fetchGithubPRs(repository *repo_v1alpha1.Repository) error {
	var githubPR *github.PullRequest // TODO: fetch list using github-api
	return c.githubPRHandler(githubPR, repository)
}