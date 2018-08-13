package controller

import (
	"github.com/appscode/go/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// TODO: use TryUpdateBindingStatus
func (c *Controller) updateRepositoryLastObservedGen(name, namespace string, generation int64) error {
	repo, err := c.gitAPIServerClient.GitV1alpha1().Repositories(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	repo.Status.LastObservedGeneration = &generation

	_, err = c.gitAPIServerClient.GitV1alpha1().Repositories(namespace).UpdateStatus(repo)
	if err != nil {
		log.Errorf("Failed to update status of repository %s/%s, reason: %s", namespace, name, err.Error())
	}

	return err
}
