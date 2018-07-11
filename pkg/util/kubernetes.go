package util

import (
	"strings"

	"github.com/appscode/go/log/golog"
	"github.com/pkg/errors"
)

const (
	RepositoryFinalizer            = "kubeci"
	SnapshotIDLength               = 8
	SnapshotIDLengthWithDashPrefix = 9
)

var (
	AnalyticsClientID string
	EnableAnalytics   = true
	LoggerOptions     golog.Options
)

type RepoLabelData struct {
	WorkloadKind string
	WorkloadName string
	PodName      string
	NodeName     string
}

func GetRepoNameAndSnapshotID(snapshotName string) (repoName, snapshotId string, err error) {
	if len(snapshotName) < 9 {
		err = errors.New("invalid snapshot name")
		return
	}
	snapshotId = snapshotName[len(snapshotName)-SnapshotIDLength:]

	repoName = strings.TrimSuffix(snapshotName, snapshotName[len(snapshotName)-SnapshotIDLengthWithDashPrefix:])
	return
}
