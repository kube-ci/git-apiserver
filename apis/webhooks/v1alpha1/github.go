package v1alpha1

import (
	"github.com/google/go-github/github"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindGithubEvent = "GithubEvent"
	ResourceGithubEvents    = "githubevents"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type GithubEvent struct {
	metav1.TypeMeta `json:",inline,omitempty"`

	Action *string            `json:"action,omitempty"`
	Repo   *github.Repository `json:"repository,omitempty"`
	Sender *github.User       `json:"sender,omitempty"`

	PullRequest *github.PullRequest `json:"pull_request,omitempty"`
}
