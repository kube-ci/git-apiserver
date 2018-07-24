package controller

import (
	"github.com/appscode/go/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kube.ci/git-apiserver/apis/git/v1alpha1"
	"kube.ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1/util"
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
		log.Errorf("failed to update status of repository %s/%s, reason: %s", namespace, name, err.Error())
	}

	return err
}

func (c *Controller) updateBindingLastObservedGen(name, namespace string, generation int64) error {
	_, err := util.TryUpdateBindingStatus(
		c.gitAPIServerClient.GitV1alpha1(),
		metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		func(binding *v1alpha1.Binding) *v1alpha1.Binding {
			binding.Status.LastObservedGeneration = &generation
			return binding
		},
	)
	if err != nil {
		log.Errorf(err.Error())
	}
	return err
}
