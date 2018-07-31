package controller

import (
	"context"

	"github.com/TamalSaha/go-oneliners"
	"github.com/appscode/go/log"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apiserver/pkg/registry/rest"
	repo_v1alpha1 "kube.ci/git-apiserver/apis/git/v1alpha1"
	"kube.ci/git-apiserver/apis/webhook/v1alpha1"
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
				c.githubPRHandler(event.PullRequest, repository)
			}
		}
	}
}

func (c *Controller) githubPRHandler(pr *v1alpha1.PullRequest, repo *repo_v1alpha1.Repository) {
	// create or patch PR CRD
}
