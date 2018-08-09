package v1alpha1

import (
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
		Kind:          ResourceKindRepository,
		ShortNames:    []string{"repo"},
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
		SpecDefinitionName:      "kube.ci/git-apiserver/apis/git/v1alpha1.Repository",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (c Branch) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceBranches,
		Kind:          ResourceKindBranch,
		ShortNames:    []string{"branch"},
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
		SpecDefinitionName:      "kube.ci/git-apiserver/apis/git/v1alpha1.Branch",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (c Tag) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceTags,
		Kind:          ResourceKindTag,
		ShortNames:    []string{"tag"},
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
		SpecDefinitionName:      "kube.ci/git-apiserver/apis/git/v1alpha1.Tag",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (c Binding) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourceBindings,
		Kind:          ResourceKindBinding,
		ShortNames:    []string{"binding"},
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
		SpecDefinitionName:      "kube.ci/git-apiserver/apis/git/v1alpha1.Binding",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}

func (c PullRequest) CustomResourceDefinition() *apiextensions.CustomResourceDefinition {
	return crdutils.NewCustomResourceDefinition(crdutils.Config{
		Group:         SchemeGroupVersion.Group,
		Plural:        ResourcePullRequests,
		Kind:          ResourceKindPullRequest,
		ShortNames:    []string{"pr"},
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
		SpecDefinitionName:      "kube.ci/git-apiserver/apis/git/v1alpha1.PullRequest",
		EnableValidation:        true,
		GetOpenAPIDefinitions:   GetOpenAPIDefinitions,
		EnableStatusSubresource: EnableStatusSubresource,
	})
}
