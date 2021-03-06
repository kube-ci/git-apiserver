/*
Copyright 2019 The KubeCI Authors.

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
	v1alpha1 "github.com/kube-ci/git-apiserver/apis/git/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeTags implements TagInterface
type FakeTags struct {
	Fake *FakeGitV1alpha1
	ns   string
}

var tagsResource = schema.GroupVersionResource{Group: "git.kube.ci", Version: "v1alpha1", Resource: "tags"}

var tagsKind = schema.GroupVersionKind{Group: "git.kube.ci", Version: "v1alpha1", Kind: "Tag"}

// Get takes name of the tag, and returns the corresponding tag object, and an error if there is any.
func (c *FakeTags) Get(name string, options v1.GetOptions) (result *v1alpha1.Tag, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(tagsResource, c.ns, name), &v1alpha1.Tag{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Tag), err
}

// List takes label and field selectors, and returns the list of Tags that match those selectors.
func (c *FakeTags) List(opts v1.ListOptions) (result *v1alpha1.TagList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(tagsResource, tagsKind, c.ns, opts), &v1alpha1.TagList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.TagList{ListMeta: obj.(*v1alpha1.TagList).ListMeta}
	for _, item := range obj.(*v1alpha1.TagList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested tags.
func (c *FakeTags) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(tagsResource, c.ns, opts))

}

// Create takes the representation of a tag and creates it.  Returns the server's representation of the tag, and an error, if there is any.
func (c *FakeTags) Create(tag *v1alpha1.Tag) (result *v1alpha1.Tag, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(tagsResource, c.ns, tag), &v1alpha1.Tag{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Tag), err
}

// Update takes the representation of a tag and updates it. Returns the server's representation of the tag, and an error, if there is any.
func (c *FakeTags) Update(tag *v1alpha1.Tag) (result *v1alpha1.Tag, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(tagsResource, c.ns, tag), &v1alpha1.Tag{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Tag), err
}

// Delete takes name of the tag and deletes it. Returns an error if one occurs.
func (c *FakeTags) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(tagsResource, c.ns, name), &v1alpha1.Tag{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeTags) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(tagsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.TagList{})
	return err
}

// Patch applies the patch and returns the patched tag.
func (c *FakeTags) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Tag, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(tagsResource, c.ns, name, pt, data, subresources...), &v1alpha1.Tag{})

	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.Tag), err
}
