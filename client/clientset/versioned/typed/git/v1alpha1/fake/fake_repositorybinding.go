/*
Copyright 2018 The KubeCI Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
	v1alpha1 "kube.ci/git-apiserver/apis/git/v1alpha1"
)

// FakeRepositoryBindings implements RepositoryBindingInterface
type FakeRepositoryBindings struct {
	Fake *FakeGitV1alpha1
	ns   string
}

var repositorybindingsResource = schema.GroupVersionResource{Group: "git.kube.ci", Version: "v1alpha1", Resource: "repositorybindings"}

var repositorybindingsKind = schema.GroupVersionKind{Group: "git.kube.ci", Version: "v1alpha1", Kind: "RepositoryBinding"}

// Get takes name of the repositoryBinding, and returns the corresponding repositoryBinding object, and an error if there is any.
func (c *FakeRepositoryBindings) Get(name string, options v1.GetOptions) (result *v1alpha1.RepositoryBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(repositorybindingsResource, c.ns, name), &v1alpha1.RepositoryBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RepositoryBinding), err
}

// List takes label and field selectors, and returns the list of RepositoryBindings that match those selectors.
func (c *FakeRepositoryBindings) List(opts v1.ListOptions) (result *v1alpha1.RepositoryBindingList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(repositorybindingsResource, repositorybindingsKind, c.ns, opts), &v1alpha1.RepositoryBindingList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.RepositoryBindingList{ListMeta: obj.(*v1alpha1.RepositoryBindingList).ListMeta}
	for _, item := range obj.(*v1alpha1.RepositoryBindingList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested repositoryBindings.
func (c *FakeRepositoryBindings) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(repositorybindingsResource, c.ns, opts))

}

// Create takes the representation of a repositoryBinding and creates it.  Returns the server's representation of the repositoryBinding, and an error, if there is any.
func (c *FakeRepositoryBindings) Create(repositoryBinding *v1alpha1.RepositoryBinding) (result *v1alpha1.RepositoryBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(repositorybindingsResource, c.ns, repositoryBinding), &v1alpha1.RepositoryBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RepositoryBinding), err
}

// Update takes the representation of a repositoryBinding and updates it. Returns the server's representation of the repositoryBinding, and an error, if there is any.
func (c *FakeRepositoryBindings) Update(repositoryBinding *v1alpha1.RepositoryBinding) (result *v1alpha1.RepositoryBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(repositorybindingsResource, c.ns, repositoryBinding), &v1alpha1.RepositoryBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RepositoryBinding), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeRepositoryBindings) UpdateStatus(repositoryBinding *v1alpha1.RepositoryBinding) (*v1alpha1.RepositoryBinding, error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateSubresourceAction(repositorybindingsResource, "status", c.ns, repositoryBinding), &v1alpha1.RepositoryBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RepositoryBinding), err
}

// Delete takes name of the repositoryBinding and deletes it. Returns an error if one occurs.
func (c *FakeRepositoryBindings) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(repositorybindingsResource, c.ns, name), &v1alpha1.RepositoryBinding{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeRepositoryBindings) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(repositorybindingsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.RepositoryBindingList{})
	return err
}

// Patch applies the patch and returns the patched repositoryBinding.
func (c *FakeRepositoryBindings) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RepositoryBinding, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(repositorybindingsResource, c.ns, name, data, subresources...), &v1alpha1.RepositoryBinding{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.RepositoryBinding), err
}