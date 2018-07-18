/*
Copyright 2017 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package branchThree

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/apiserver/pkg/registry/generic"
	"k8s.io/apiserver/pkg/storage"
	"k8s.io/apiserver/pkg/storage/names"
	repo_v1alpha1 "kube.ci/git-apiserver/apis/repositories/v1alpha1"
)

// NewStrategy creates and returns a branchStrategy instance
func NewStrategy(typer runtime.ObjectTyper) branchStrategy {
	return branchStrategy{typer, names.SimpleNameGenerator}
}

// GetAttrs returns labels.Set, fields.Set, the presence of Initializers if any
// and error in case the given runtime.Object is not a Fischer
func GetAttrs(obj runtime.Object) (labels.Set, fields.Set, bool, error) {
	apiserver, ok := obj.(*repo_v1alpha1.Branch)
	if !ok {
		return nil, nil, false, fmt.Errorf("given object is not a Fischer")
	}
	return labels.Set(apiserver.ObjectMeta.Labels), SelectableFields(apiserver), apiserver.Initializers != nil, nil
}

// MatchFischer is the filter used by the generic etcd backend to watch events
// from etcd to clients of the apiserver only interested in specific labels/fields.
func MatchFischer(label labels.Selector, field fields.Selector) storage.SelectionPredicate {
	return storage.SelectionPredicate{
		Label:    label,
		Field:    field,
		GetAttrs: GetAttrs,
	}
}

// SelectableFields returns a field set that represents the object.
func SelectableFields(obj *repo_v1alpha1.Branch) fields.Set {
	return generic.ObjectMetaFieldsSet(&obj.ObjectMeta, true)
}

type branchStrategy struct {
	runtime.ObjectTyper
	names.NameGenerator
}

func (branchStrategy) NamespaceScoped() bool {
	return false
}

func (branchStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
}

func (branchStrategy) PrepareForUpdate(ctx context.Context, obj, old runtime.Object) {
}

func (branchStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	return field.ErrorList{}
}

func (branchStrategy) AllowCreateOnUpdate() bool {
	return false
}

func (branchStrategy) AllowUnconditionalUpdate() bool {
	return false
}

func (branchStrategy) Canonicalize(obj runtime.Object) {
}

func (branchStrategy) ValidateUpdate(ctx context.Context, obj, old runtime.Object) field.ErrorList {
	return field.ErrorList{}
}
