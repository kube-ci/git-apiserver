package snapshot

import (
	"context"

	"github.com/pkg/errors"
	metainternalversion "k8s.io/apimachinery/pkg/apis/meta/internalversion"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	apirequest "k8s.io/apiserver/pkg/endpoints/request"
	"k8s.io/apiserver/pkg/registry/rest"
	"k8s.io/client-go/kubernetes"
	restconfig "k8s.io/client-go/rest"
	kubeci "kube.ci/git-apiserver/apis/git/v1alpha1"
	"kube.ci/git-apiserver/apis/repositories"
	"kube.ci/git-apiserver/client/clientset/versioned"
	"kube.ci/git-apiserver/pkg/util"
)

type REST struct {
	kubeciClient versioned.Interface
	kubeClient   kubernetes.Interface
	config       *restconfig.Config
}

var _ rest.Scoper = &REST{}
var _ rest.Getter = &REST{}
var _ rest.Lister = &REST{}
var _ rest.GracefulDeleter = &REST{}

func NewREST(config *restconfig.Config) *REST {
	return &REST{
		kubeciClient: versioned.NewForConfigOrDie(config),
		kubeClient:   kubernetes.NewForConfigOrDie(config),
		config:       config,
	}
}

func (r *REST) NamespaceScoped() bool {
	return true
}

func (r *REST) New() runtime.Object {
	return &repositories.Snapshot{}
}

func (r *REST) Get(ctx context.Context, name string, options *metav1.GetOptions) (runtime.Object, error) {
	ns, ok := apirequest.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.New("missing namespace")
	}
	if len(name) < 9 {
		return nil, errors.New("invalid snapshot name")
	}

	repoName, snapshotId, err := util.GetRepoNameAndSnapshotID(name)
	if err != nil {
		return nil, err
	}

	repo, err := r.kubeciClient.GitV1alpha1().Repositories(ns).Get(repoName, metav1.GetOptions{})
	if err != nil {
		return nil, errors.New("respective repository not found. error:" + err.Error())
	}

	snapshots := make([]repositories.Snapshot, 0)
	snapshots, err = r.GetSnapshots(repo, []string{snapshotId})
	if err != nil {
		return nil, err
	}

	if len(snapshots) == 0 {
		return nil, errors.New("no resource found")
	}

	snapshot := &repositories.Snapshot{}
	snapshot = &snapshots[0]
	return snapshot, nil
}

func (r *REST) NewList() runtime.Object {
	return &repositories.SnapshotList{}
}

func (r *REST) List(ctx context.Context, options *metainternalversion.ListOptions) (runtime.Object, error) {
	ns, ok := apirequest.NamespaceFrom(ctx)
	if !ok {
		return nil, errors.New("missing namespace")
	}

	repos, err := r.kubeciClient.GitV1alpha1().Repositories(ns).List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}

	var selectedRepos []kubeci.Repository
	if options.LabelSelector != nil {
		for _, r := range repos.Items {
			repoLabels := make(map[string]string)
			repoLabels = r.Labels
			repoLabels["repository"] = r.Name
			if options.LabelSelector.Matches(labels.Set(repoLabels)) {
				selectedRepos = append(selectedRepos, r)
			}
		}
	} else {
		selectedRepos = repos.Items
	}

	snapshotList := &repositories.SnapshotList{
		Items: make([]repositories.Snapshot, 0),
	}
	for _, repo := range selectedRepos {
		var snapshots []repositories.Snapshot
		snapshots, err = r.GetSnapshots(&repo, nil)
		if err != nil {
			return nil, err
		}
		snapshotList.Items = append(snapshotList.Items, snapshots...)
	}

	// k8s.io/apimachinery/pkg/apis/meta/v1/unstructured/unstructured_list.go
	// unstructured.UnstructuredList{}
	return snapshotList, nil
}

func (r *REST) Delete(ctx context.Context, name string, options *metav1.DeleteOptions) (runtime.Object, bool, error) {
	ns, ok := apirequest.NamespaceFrom(ctx)
	if !ok {
		return nil, false, errors.New("missing namespace")
	}
	repoName, snapshotId, err := util.GetRepoNameAndSnapshotID(name)
	if err != nil {
		return nil, false, err
	}
	repo, err := r.kubeciClient.GitV1alpha1().Repositories(ns).Get(repoName, metav1.GetOptions{})
	if err != nil {
		return nil, false, errors.New("respective repository not found. error:" + err.Error())
	}

	err = r.ForgetSnapshots(repo, []string{snapshotId})
	if err != nil {
		return nil, false, err
	}

	return nil, true, nil
}
