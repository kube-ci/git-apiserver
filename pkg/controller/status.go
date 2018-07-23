package controller

import (
	"github.com/appscode/go/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (c *Controller) updateRepositoryLastObservedGen(name, namespace string, generation int64) error {
	repo, err := c.gitAPIServerClient.GitV1alpha1().Repositories(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	repo.Status.LastObservedGeneration = &generation

	_, err = c.gitAPIServerClient.GitV1alpha1().Repositories(namespace).UpdateStatus(repo)
	if err != nil {
		log.Errorf("failed to update status of repository %s/%s, reason: %s", namespace, name, err.Error())
	}

	return err
}

func (c *Controller) updateBindingLastObservedGen(name, namespace string, generation int64) error {
	binding, err := c.gitAPIServerClient.GitV1alpha1().Bindings(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		return err
	}

	binding.Status.LastObservedGeneration = &generation

	_, err = c.gitAPIServerClient.GitV1alpha1().Bindings(namespace).UpdateStatus(binding)
	if err != nil {
		log.Errorf("failed to update status of repository %s/%s, reason: %s", namespace, name, err.Error())
	}

	return err
}
