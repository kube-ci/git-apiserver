package branchTwo

import (
	"context"
	"log"
	"sync"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/tools/cache"
	repo_v1alpha1 "kube.ci/git-apiserver/apis/repositories/v1alpha1"
)

type REST struct {
	Indexer        cache.Indexer
	BranchWatchers []*BranchWatcher
}

var _ rest.Getter = &REST{}
var _ rest.Lister = &REST{}
var _ rest.Watcher = &REST{}
var _ rest.GroupVersionKindProvider = &REST{}
var _ rest.Scoper = &REST{}

func NewREST() *REST {
	log.Println("NewREST...")
	return &REST{
		Indexer: cache.NewIndexer(
			cache.DeletionHandlingMetaNamespaceKeyFunc,
			cache.Indexers{cache.NamespaceIndex: cache.MetaNamespaceIndexFunc},
		),
	}
}

func (r *REST) NamespaceScoped() bool {
	return true
}

func (r *REST) New() runtime.Object {
	return &repo_v1alpha1.Branch{}
}

func (p *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return repo_v1alpha1.SchemeGroupVersion.WithKind(repo_v1alpha1.ResourceKindBranch)
}

func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	log.Println("Get...")

	ns, ok := apirequest.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("missing namespace")
	}
	if len(name) == 0 {
		return nil, errors.NewBadRequest("missing search query")
	}

	obj, exists, err := r.Indexer.GetByKey(ns + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(repo_v1alpha1.Resource(repo_v1alpha1.ResourceBranches), name)
	}
	return obj.(*repo_v1alpha1.Branch), nil
}

func (r *REST) NewList() runtime.Object {
	return &repo_v1alpha1.BranchList{}
}

func (r *REST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	log.Println("List...")

	ns, ok := apirequest.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.NewBadRequest("missing namespace")
	}

	resp := &repo_v1alpha1.BranchList{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "repositories.git.kube.ci/v1alpha1",
			Kind:       repo_v1alpha1.ResourceKindBranch,
		},
	}

	var selector labels.Selector
	if options != nil {
		selector = labels.Everything() // TODO: fix nil pointer panic
	} else {
		selector = labels.Everything()
	}

	err := cache.ListAllByNamespace(r.Indexer, ns, selector, func(m interface{}) {
		branch := m.(*repo_v1alpha1.Branch)
		resp.Items = append(resp.Items, *branch)
	})

	return resp, err
}

func (r *REST) Watch(ctx context.Context, options *metainternalversion.ListOptions) (watch.Interface, error) {
	log.Println("Watch...")
	return r.NewBranchWatcher(options), nil
}

// watch interface

type BranchWatcher struct {
	result  chan watch.Event
	Stopped bool
	options *metainternalversion.ListOptions
	sync.Mutex
}

func (f *BranchWatcher) Stop() {
	f.Lock()
	defer f.Unlock()
	if !f.Stopped {
		log.Println("Stopping branch watcher...")
		close(f.result)
		f.Stopped = true
	}
}

func (f *BranchWatcher) ResultChan() <-chan watch.Event {
	return f.result
}

func (r *REST) NewBranchWatcher(options *metainternalversion.ListOptions) *BranchWatcher {
	branchWatcher := &BranchWatcher{
		options: options,
		result:  make(chan watch.Event),
	}
	r.BranchWatchers = append(r.BranchWatchers, branchWatcher)
	return branchWatcher
}

func (r *REST) CreateOrUpdateBranch(branch *repo_v1alpha1.Branch) {
	log.Println("CreateOrUpdateBranch...")

	key := branch.Namespace + "/" + branch.Name
	event := watch.Event{}

	_, exists, _ := r.Indexer.GetByKey(key)
	if exists {
		if err := r.Indexer.Update(branch); err != nil {
			log.Println("Error...CreateOrUpdateBranch...", err)
			return
		}
		event = watch.Event{
			Type:   watch.Modified,
			Object: branch,
		}
	} else {
		if err := r.Indexer.Add(branch); err != nil {
			log.Println("Error...CreateOrUpdateBranch...", err)
			return
		}
		event = watch.Event{
			Type:   watch.Added,
			Object: branch,
		}
	}

	// push event to channels
	for _, branchWatcher := range r.BranchWatchers {
		if !branchWatcher.Stopped {
			branchWatcher.result <- event
		}
	}
}

func (r *REST) DeleteBranch(branch *repo_v1alpha1.Branch) {
	log.Println("DeleteBranch...")

	key := branch.Namespace + "/" + branch.Name

	_, exists, _ := r.Indexer.GetByKey(key)
	if !exists {
		return
	}

	if err := r.Indexer.Delete(branch); err != nil {
		log.Println("Error...DeleteBranch...", err)
		return
	}

	event := watch.Event{
		Type:   watch.Deleted,
		Object: branch,
	}

	// push event to channels
	for _, branchWatcher := range r.BranchWatchers {
		if !branchWatcher.Stopped {
			branchWatcher.result <- event
		}
	}
}
