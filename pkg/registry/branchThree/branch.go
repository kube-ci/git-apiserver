package branchThree

import (
	"context"
	"fmt"
	"log"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
	genericregistry "k8s.io/apiserver/pkg/registry/generic/registry"
	"k8s.io/apiserver/pkg/registry/rest"
	repo_v1alpha1 "kube.ci/git-apiserver/apis/repositories/v1alpha1"
)

type REST struct {
	*genericregistry.Store
}

func NewREST(scheme *runtime.Scheme, optsGetter generic.RESTOptionsGetter) (*REST, error) {
	strategy := NewStrategy(scheme)
	store := &genericregistry.Store{
		NewFunc:                  func() runtime.Object { return &repo_v1alpha1.Branch{} },
		NewListFunc:              func() runtime.Object { return &repo_v1alpha1.BranchList{} },
		PredicateFunc:            MatchFischer,
		DefaultQualifiedResource: repo_v1alpha1.Resource(repo_v1alpha1.ResourceBranches),
		CreateStrategy:           strategy,
		UpdateStrategy:           strategy,
		DeleteStrategy:           strategy,
	}
	options := &generic.StoreOptions{RESTOptions: optsGetter, AttrFunc: GetAttrs}
	if err := store.CompleteWithOptions(options); err != nil {
		return nil, err
	}
	return &REST{store}, nil
}

func (r *REST) CreateOrUpdateBranch(branch *repo_v1alpha1.Branch) {
	log.Println("CreateOrUpdateBranch...")
	_, err := r.Create(context.Background(),
		branch,
		rest.ValidateAllObjectFunc,
		true,
	)
	if err != nil {
		fmt.Println("Error...CreateOrUpdateBranch...", err)
	}
}
