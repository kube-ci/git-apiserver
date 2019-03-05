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

// Code generated by informer-gen. DO NOT EDIT.

package v1alpha1

import (
	time "time"

	gitv1alpha1 "github.com/kube-ci/git-apiserver/apis/git/v1alpha1"
	versioned "github.com/kube-ci/git-apiserver/client/clientset/versioned"
	internalinterfaces "github.com/kube-ci/git-apiserver/client/informers/externalversions/internalinterfaces"
	v1alpha1 "github.com/kube-ci/git-apiserver/client/listers/git/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	runtime "k8s.io/apimachinery/pkg/runtime"
	watch "k8s.io/apimachinery/pkg/watch"
	cache "k8s.io/client-go/tools/cache"
)

// BranchInformer provides access to a shared informer and lister for
// Branches.
type BranchInformer interface {
	Informer() cache.SharedIndexInformer
	Lister() v1alpha1.BranchLister
}

type branchInformer struct {
	factory          internalinterfaces.SharedInformerFactory
	tweakListOptions internalinterfaces.TweakListOptionsFunc
	namespace        string
}

// NewBranchInformer constructs a new informer for Branch type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewBranchInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers) cache.SharedIndexInformer {
	return NewFilteredBranchInformer(client, namespace, resyncPeriod, indexers, nil)
}

// NewFilteredBranchInformer constructs a new informer for Branch type.
// Always prefer using an informer factory to get a shared informer instead of getting an independent
// one. This reduces memory footprint and number of connections to the server.
func NewFilteredBranchInformer(client versioned.Interface, namespace string, resyncPeriod time.Duration, indexers cache.Indexers, tweakListOptions internalinterfaces.TweakListOptionsFunc) cache.SharedIndexInformer {
	return cache.NewSharedIndexInformer(
		&cache.ListWatch{
			ListFunc: func(options v1.ListOptions) (runtime.Object, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GitV1alpha1().Branches(namespace).List(options)
			},
			WatchFunc: func(options v1.ListOptions) (watch.Interface, error) {
				if tweakListOptions != nil {
					tweakListOptions(&options)
				}
				return client.GitV1alpha1().Branches(namespace).Watch(options)
			},
		},
		&gitv1alpha1.Branch{},
		resyncPeriod,
		indexers,
	)
}

func (f *branchInformer) defaultInformer(client versioned.Interface, resyncPeriod time.Duration) cache.SharedIndexInformer {
	return NewFilteredBranchInformer(client, f.namespace, resyncPeriod, cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc}, f.tweakListOptions)
}

func (f *branchInformer) Informer() cache.SharedIndexInformer {
	return f.factory.InformerFor(&gitv1alpha1.Branch{}, f.defaultInformer)
}

func (f *branchInformer) Lister() v1alpha1.BranchLister {
	return v1alpha1.NewBranchLister(f.Informer().GetIndexer())
}
