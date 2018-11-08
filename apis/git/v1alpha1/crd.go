package v1alpha1

import (
	"strings"

	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	apiextensions "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
)

var (
	EnableStatusSubresource bool
)

func (c Repository) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceRepositories,
		Singular:      strings.ToLower(ResourceKindRepository),
		Kind:          ResourceKindRepository,
		ShortNames:    []string{"repo"},
		Categories:    []string{"kubeci", "ci", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubeci"},
		},
		SpecDefinitionName:      "github.com/kube-ci/git-apiserver/apis/git/v1alpha1.Repository",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (c Branch) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceBranches,
		Singular:      strings.ToLower(ResourceKindBranch),
		Kind:          ResourceKindBranch,
		ShortNames:    []string{"branch"},
		Categories:    []string{"kubeci", "ci", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubeci"},
		},
		SpecDefinitionName:      "github.com/kube-ci/git-apiserver/apis/git/v1alpha1.Branch",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (c Tag) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceTags,
		Singular:      strings.ToLower(ResourceKindTag),
		Kind:          ResourceKindTag,
		ShortNames:    []string{"tag"},
		Categories:    []string{"kubeci", "ci", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubeci"},
		},
		SpecDefinitionName:      "github.com/kube-ci/git-apiserver/apis/git/v1alpha1.Tag",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (c PullRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePullRequests,
		Singular:      strings.ToLower(ResourceKindPullRequest),
		Kind:          ResourceKindPullRequest,
		ShortNames:    []string{"pr"},
		Categories:    []string{"kubeci", "ci", "appscode", "all"},
		ResourceScope: string(apiextensions.NamespaceScoped),
		Versions: []apiextensions.CustomResourceDefinitionVersion{
			{
				Name:    SchemeGroupVersion.Version,
				Served:  true,
				Storage: true,
			},
		},
		Labels: crdutils.Labels{
			LabelsMap: map[string]string{"app": "kubeci"},
		},
		SpecDefinitionName:      "github.com/kube-ci/git-apiserver/apis/git/v1alpha1.PullRequest",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}
