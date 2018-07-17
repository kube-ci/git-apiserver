package util

import (
	"github.com/appscode/go/log/golog"
)

const (
	RepositoryFinalizer = "git-apiserver"
)

var (
	AnalyticsClientID string
	EnableAnalytics   = true
	LoggerOptions     golog.Options
)
