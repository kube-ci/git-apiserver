package branch

import (
	"context"
	"log"
	"sync"

	"k8s.io/apimachinery/pkg/api/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	repo_v1alpha1 "kube.ci/git-apiserver/apis/repositories/v1alpha1"
)

type REST struct {
	Branches       map[string]*repo_v1alpha1.Branch
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
		Branches: make(map[string]*repo_v1alpha1.Branch, 0),
	}
}

func (r *REST) NamespaceScoped() bool {
	return true
}

func (r *REST) New() runtime.Object {
	return &repo_v1alpha1.Branch{}
}

func (p *REST) GroupVersionKind(containingGV schema.GroupVersion) schema.GroupVersionKind {
	return repo_v1alpha1.SchemeGroupVersion.WithKind("Branch")
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

	key := ns + "/" + name
	log.Println("Key...", key)

	if branch, ok := r.Branches[key]; ok {
		return branch, nil
	}

	return nil, errors.NewNotFound(repo_v1alpha1.Resource(repo_v1alpha1.ResourceBranches), name)
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

	for _, branch := range r.Branches {
		if branch.Namespace == ns {
			resp.Items = append(resp.Items, *branch)
		}
	}

	return resp, nil
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

	if _, ok := r.Branches[key]; !ok {
		r.Branches[key] = branch
		event = watch.Event{
			Type:   watch.Added,
			Object: branch,
		}
	} else {
		r.Branches[key] = branch
		event = watch.Event{
			Type:   watch.Modified,
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
	if _, ok := r.Branches[key]; !ok {
		return
	}

	delete(r.Branches, key)
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
