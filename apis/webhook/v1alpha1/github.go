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

	// TODO: similar to TypeMeta, but deepcopy-gen not working
	// ActivityEvent   `json:",inline,omitempty"`

	Action      *string      `json:"action,omitempty"`
	Repo        *Repository  `json:"repository,omitempty"`
	Sender      *User        `json:"sender,omitempty"`
	Issue       *Issue       `json:"issue,omitempty"` // TODO: for test only, we don't need issue events
	PullRequest *PullRequest `json:"pull_request,omitempty"`
}

/*
// TODO: similar to TypeMeta, but deepcopy-gen not working
// +k8s:deepcopy-gen=false
type ActivityEvent struct {
	Action *string            `json:"action,omitempty"`
	Repo   *github.Repository `json:"repository,omitempty"`
	Sender *github.User       `json:"sender,omitempty"`

	Issue       *github.Issue       `json:"issue,omitempty"`
	PullRequest *github.PullRequest `json:"pull_request,omitempty"`
}
*/

// +k8s:deepcopy-gen=false
type Repository struct {
	*github.Repository `json:",inline,omitempty"`
}

// +k8s:deepcopy-gen=false
type User struct {
	*github.User `json:",inline,omitempty"`
}

// +k8s:deepcopy-gen=false
type Issue struct {
	*github.Issue `json:",inline,omitempty"`
}

// +k8s:deepcopy-gen=false
type PullRequest struct {
	*github.PullRequest `json:",inline,omitempty"`
}
