package v1alpha1

import (
	"github.com/appscode/go/encoding/json/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

const (
	ResourceKindRepository = "Repository"
	ResourceRepositories   = "repositories"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RepositorySpec   `json:"spec,omitempty"`
	Status            RepositoryStatus `json:"status,omitempty"`
}

type RepositorySpec struct {
	Host            string  `json:"host,omitempty"` // github, gitlab // TODO: use type
	Owner           string  `json:"owner,omitempty"`
	Repo            string  `json:"repo,omitempty"`
	CloneUrl        string  `json:"cloneUrl,omitempty"`
	TokenFormSecret *string `json:"tokenFormSecret,omitempty"` // secret name, secret must have field 'token'
}

type RepositoryStatus struct {
	// observedGeneration is the most recent generation observed for this resource. It corresponds to the
	// resource's generation, which is updated on mutation by the API Server.
	// +optional
	ObservedGeneration *types.IntHash `json:"observedGeneration,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repository `json:"items,omitempty"`
}

func (repo *Repository) GetToken(kubeClient kubernetes.Interface) (string, error) {
	if repo.Spec.TokenFormSecret == nil {
		return "", nil
	}
	secret, err := kubeClient.CoreV1().Secrets(repo.Namespace).Get(*repo.Spec.TokenFormSecret, metav1.GetOptions{})
	if err != nil {
		return "", err
	}
	return string(secret.Data["token"]), nil
}
