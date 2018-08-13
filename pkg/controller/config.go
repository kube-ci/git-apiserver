package controller

import (
	"time"

	core "k8s.io/api/core/v1"
	crd_cs "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/typed/apiextensions/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	cs "kube.ci/git-apiserver/client/clientset/versioned"
	git_apiserver_informers "kube.ci/git-apiserver/client/informers/externalversions"
	"kube.ci/git-apiserver/pkg/eventer"
)

type config struct {
	EnableRBAC           bool
	GitAPIServerImageTag string
	DockerRegistry       string
	MaxNumRequeues       int
	NumThreads           int
	ResyncPeriod         time.Duration
}

type Config struct {
	config

	ClientConfig       *rest.Config
	KubeClient         kubernetes.Interface
	GitAPIServerClient cs.Interface
	CRDClient          crd_cs.ApiextensionsV1beta1Interface
}

func NewConfig(clientConfig *rest.Config) *Config {
	return &Config{
		ClientConfig: clientConfig,
	}
}

func (c *Config) New() (*Controller, error) {
	tweakListOptions := func(opt *metav1.ListOptions) {
		opt.IncludeUninitialized = true
	}
	ctrl := &Controller{
		config:                      c.config,
		kubeClient:                  c.KubeClient,
		gitAPIServerClient:          c.GitAPIServerClient,
		crdClient:                   c.CRDClient,
		kubeInformerFactory:         informers.NewFilteredSharedInformerFactory(c.KubeClient, c.ResyncPeriod, core.NamespaceAll, tweakListOptions),
		gitAPIServerInformerFactory: git_apiserver_informers.NewSharedInformerFactory(c.GitAPIServerClient, c.ResyncPeriod),
		recorder:                    eventer.NewEventRecorder(c.KubeClient, "kubeci-controller"),
	}

	if err := ctrl.ensureCustomResourceDefinitions(); err != nil {
		return nil, err
	}

	ctrl.initRepositoryWatcher()

	return ctrl, nil
}
