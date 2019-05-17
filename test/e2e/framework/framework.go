package framework

import (
	"path/filepath"

	"github.com/appscode/go/crypto/rand"
	cs "github.com/kube-ci/git-apiserver/client/clientset/versioned"
	. "github.com/onsi/gomega"
	"github.com/spf13/afero"
	"gomodules.xyz/cert/certstore"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	ka "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
)

type Framework struct {
	KubeClient         kubernetes.Interface
	GitAPIServerClient cs.Interface
	KAClient           ka.Interface
	namespace          string
	CertStore          *certstore.CertStore
	WebhookEnabled     bool
	SelfHostedOperator bool
	ClientConfig       *rest.Config
}

func New(kubeClient kubernetes.Interface, extClient cs.Interface, kaClient ka.Interface, webhookEnabled bool, selfHostedOperator bool, clientConfig *rest.Config) *Framework {
	store, err := certstore.NewCertStore(afero.NewMemMapFs(), filepath.Join("", "pki"))
	Expect(err).NotTo(HaveOccurred())

	err = store.InitCA()
	Expect(err).NotTo(HaveOccurred())

	return &Framework{
		KubeClient:         kubeClient,
		GitAPIServerClient: extClient,
		KAClient:           kaClient,
		namespace:          rand.WithUniqSuffix("test-git-apiserver"),
		CertStore:          store,
		WebhookEnabled:     webhookEnabled,
		SelfHostedOperator: selfHostedOperator,
		ClientConfig:       clientConfig,
	}
}

func (f *Framework) Invoke() *Invocation {
	return &Invocation{
		Framework: f,
		app:       rand.WithUniqSuffix("git-apiserver-e2e"),
	}
}

func (f *Invocation) AppLabel() string {
	return "app=" + f.app
}

func (f *Invocation) App() string {
	return f.app
}

type Invocation struct {
	*Framework
	app string
}
