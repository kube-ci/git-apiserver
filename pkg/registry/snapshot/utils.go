package snapshot

import (
	kubeci "kube.ci/git-apiserver/apis/git/v1alpha1"
	"kube.ci/git-apiserver/apis/repositories"
)

func (r *REST) GetSnapshots(repository *kubeci.Repository, snapshotIDs []string) ([]repositories.Snapshot, error) {
	snapshots := make([]repositories.Snapshot, 0)
	return snapshots, nil
}

func (r *REST) ForgetSnapshots(repository *kubeci.Repository, snapshotIDs []string) error {

	return nil
}
