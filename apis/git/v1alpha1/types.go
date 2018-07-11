package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	ResourceKindRepository     = "Repository"
	ResourcePluralRepository   = "repositories"
	ResourceSingularRepository = "repository"
)

// +genclient
// +k8s:openapi-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              RepositorySpec   `json:"spec,omitempty"`
	Status            RepositoryStatus `json:"status,omitempty"`
}

type RepositorySpec struct {
	Backend Backend `json:"backend,omitempty"`
	// If true, delete respective restic repository
	// +optional
	WipeOut bool `json:"wipeOut,omitempty"`
}

type RepositoryStatus struct {
	FirstBackupTime          *metav1.Time `json:"firstBackupTime,omitempty"`
	LastBackupTime           *metav1.Time `json:"lastBackupTime,omitempty"`
	LastSuccessfulBackupTime *metav1.Time `json:"lastSuccessfulBackupTime,omitempty"`
	LastBackupDuration       string       `json:"lastBackupDuration,omitempty"`
	BackupCount              int64        `json:"backupCount,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repository `json:"items,omitempty"`
}

type Backend struct {
}
