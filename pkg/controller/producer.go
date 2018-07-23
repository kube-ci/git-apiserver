package controller

import (
	"time"

	"github.com/appscode/go/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
)

func (c *Controller) ReconcileProducer(name, namespace string) {
	log.Infoln("Reconcile Producer...")
	go func() {
		for {
			repository, err := c.gitAPIServerClient.GitV1alpha1().Repositories(namespace).Get(name, metav1.GetOptions{})
			if err != nil {
				log.Errorln(err)
				break
			}
			if err = runGitProducer(repository); err != nil {
				log.Errorln(err)
				break
			}
			time.Sleep(time.Minute)
		}
	}()
}

func runGitProducer(repository *api.Repository) error {
	return nil
}
