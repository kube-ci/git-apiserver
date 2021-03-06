package e2e_test

import (
	"strings"
	"testing"
	"time"

	logs "github.com/appscode/go/log/golog"
	api "github.com/kube-ci/git-apiserver/apis/git/v1alpha1"
	"github.com/kube-ci/git-apiserver/client/clientset/versioned/scheme"
	_ "github.com/kube-ci/git-apiserver/client/clientset/versioned/scheme"
	"github.com/kube-ci/git-apiserver/pkg/controller"
	"github.com/kube-ci/git-apiserver/test/e2e/framework"
	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/reporters"
	. "github.com/onsi/gomega"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/client-go/discovery"
	clientsetscheme "k8s.io/client-go/kubernetes/scheme"
	ka "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	"k8s.io/kubernetes/pkg/api/legacyscheme"
	crdutils "kmodules.xyz/client-go/apiextensions/v1beta1"
	discovery_util "kmodules.xyz/client-go/discovery"
	"kmodules.xyz/client-go/meta"
	"kmodules.xyz/client-go/tools/cli"
	"kmodules.xyz/client-go/tools/clientcmd"
)

const (
	TIMEOUT                  = 20 * time.Minute
	TestGitAPIServerImageTag = "canary"
)

var (
	ctrl *controller.Controller
	root *framework.Framework
)

func TestE2e(t *testing.T) {
	logs.InitLogs()
	RegisterFailHandler(Fail)
	SetDefaultEventuallyTimeout(TIMEOUT)
	junitReporter := reporters.NewJUnitReporter("junit.xml")
	RunSpecsWithDefaultAndCustomReporters(t, "e2e Suite", []Reporter{junitReporter})
}

var _ = BeforeSuite(func() {
	scheme.AddToScheme(clientsetscheme.Scheme)
	scheme.AddToScheme(legacyscheme.Scheme)
	cli.LoggerOptions.Verbosity = "5"

	clientConfig, err := clientcmd.BuildConfigFromContext(options.KubeConfig, options.KubeContext)
	Expect(err).NotTo(HaveOccurred())
	ctrlConfig := controller.NewConfig(clientConfig)

	discClient, err := discovery.NewDiscoveryClientForConfig(clientConfig)
	serverVersion, err := discovery_util.GetBaseVersion(discClient)
	Expect(err).NotTo(HaveOccurred())
	if strings.Compare(serverVersion, "1.11") >= 0 {
		api.EnableStatusSubresource = true
	}

	err = options.ApplyTo(ctrlConfig)
	Expect(err).NotTo(HaveOccurred())

	kaClient := ka.NewForConfigOrDie(clientConfig)

	root = framework.New(ctrlConfig.KubeClient, ctrlConfig.GitAPIServerClient, kaClient, options.EnableWebhook, options.SelfHostedOperator, clientConfig)
	err = root.CreateTestNamespace()
	Expect(err).NotTo(HaveOccurred())
	By("Using test namespace " + root.Namespace())

	if !options.SelfHostedOperator {
		crds := []*crd_api.CustomResourceDefinition{
			api.Repository{}.CustomResourceDefinition(),
			api.Tag{}.CustomResourceDefinition(),
			api.Branch{}.CustomResourceDefinition(),
		}

		By("Registering CRDs")
		err = crdutils.RegisterCRDs(ctrlConfig.CRDClient, crds)
		Expect(err).NotTo(HaveOccurred())

		go root.StartAPIServerAndOperator(options.KubeConfig, options.ExtraOptions)
	}

	By("Waiting for APIServer to be ready")
	root.EventuallyAPIServerReady().Should(Succeed())
	time.Sleep(time.Second * 5) // let's API server be warmed up
})

var _ = AfterSuite(func() {
	By("Cleaning API server and Webhook stuff")

	if options.EnableWebhook && !options.SelfHostedOperator {
		root.KubeClient.AdmissionregistrationV1beta1().MutatingWebhookConfigurations().Delete("admission.git.kube.ci", meta.DeleteInBackground())
		root.KubeClient.AdmissionregistrationV1beta1().ValidatingWebhookConfigurations().Delete("admission.git.kube.ci", meta.DeleteInBackground())
	}

	if !options.SelfHostedOperator {
		root.KubeClient.CoreV1().Endpoints(root.Namespace()).Delete("git-dev-apiserver", meta.DeleteInBackground())
		root.KubeClient.CoreV1().Services(root.Namespace()).Delete("git-dev-apiserver", meta.DeleteInBackground())
		root.KAClient.ApiregistrationV1beta1().APIServices().Delete("v1alpha1.admission.git.kube.ci", meta.DeleteInBackground())
		root.KAClient.ApiregistrationV1beta1().APIServices().Delete("v1alpha1.webhooks.git.kube.ci", meta.DeleteInBackground())
	}
	root.DeleteNamespace(root.Namespace())
})
