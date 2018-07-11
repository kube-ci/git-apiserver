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
	kubeciinformers "kube.ci/git-apiserver/client/informers/externalversions"
	"kube.ci/git-apiserver/pkg/eventer"
)

type config struct {
	EnableRBAC     bool
	StashImageTag  string
	DockerRegistry string
	MaxNumRequeues int
	NumThreads     int
	ResyncPeriod   time.Duration
}

type Config struct {
	config

	ClientConfig *rest.Config
	KubeClient   kubernetes.Interface
	StashClient  cs.Interface
	CRDClient    crd_cs.ApiextensionsV1beta1Interface
}

func NewConfig(clientConfig *rest.Config) *Config {
	return &Config{
		ClientConfig: clientConfig,
	}
}

func (c *Config) New() (*StashController, error) {
	tweakListOptions := func(opt *metav1.ListOptions) {
		opt.IncludeUninitialized = true
	}
	ctrl := &StashController{
		config:                c.config,
		kubeClient:            c.KubeClient,
		kubeciClient:          c.StashClient,
		crdClient:             c.CRDClient,
		kubeInformerFactory:   informers.NewFilteredSharedInformerFactory(c.KubeClient, c.ResyncPeriod, core.NamespaceAll, tweakListOptions),
		kubeciInformerFactory: kubeciinformers.NewSharedInformerFactory(c.StashClient, c.ResyncPeriod),
		recorder:              eventer.NewEventRecorder(c.KubeClient, "kubeci-controller"),
	}

	if err := ctrl.ensureCustomResourceDefinitions(); err != nil {
		return nil, err
	}

	ctrl.initRepositoryWatcher()

	return ctrl, nil
}
