package framework

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"time"

	shell "github.com/codeskyblue/go-sh"
	srvr "github.com/kube-ci/git-apiserver/pkg/cmds/server"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	genericapiserver "k8s.io/apiserver/pkg/server"
	kapi "k8s.io/kube-aggregator/pkg/apis/apiregistration/v1beta1"
)

func (f *Framework) NewTestGitAPIServerOptions(kubeConfigPath string, controllerOptions *srvr.ExtraOptions) *srvr.GitAPIServerOptions {
	opt := srvr.NewGitAPIServerOptions(os.Stdout, os.Stderr)
	opt.RecommendedOptions.Authentication.RemoteKubeConfigFile = kubeConfigPath
	//opt.RecommendedOptions.Authentication.SkipInClusterLookup = true
	opt.RecommendedOptions.Authorization.RemoteKubeConfigFile = kubeConfigPath
	opt.RecommendedOptions.CoreAPI.CoreAPIKubeconfigPath = kubeConfigPath
	opt.RecommendedOptions.SecureServing.BindPort = 6443
	opt.RecommendedOptions.SecureServing.BindAddress = net.ParseIP("127.0.0.1")
	opt.ExtraOptions = controllerOptions
	opt.StdErr = os.Stderr
	opt.StdOut = os.Stdout

	return opt
}

func (f *Framework) StartAPIServerAndOperator(kubeConfigPath string, extraOptions *srvr.ExtraOptions) {
	defer GinkgoRecover()

	sh := shell.NewSession()
	args := []interface{}{"--namespace", f.Namespace(), "--test=true"}
	if !f.WebhookEnabled {
		args = append(args, "--enable-webhook=false")
	}
	runScript := filepath.Join("..", "..", "hack", "dev", "run.sh")

	By("Creating API server and webhook stuffs")
	cmd := sh.Command(runScript, args...)
	err := cmd.Run()
	Expect(err).ShouldNot(HaveOccurred())

	By("Starting Server and Operator")
	stopCh := genericapiserver.SetupSignalHandler()
	gitAPIServerOptions := f.NewTestGitAPIServerOptions(kubeConfigPath, extraOptions)
	err = gitAPIServerOptions.Run(stopCh)
	Expect(err).ShouldNot(HaveOccurred())
}

func (f *Framework) EventuallyAPIServerReady() GomegaAsyncAssertion {
	apiServices := []string{
		"v1alpha1.admission.git.kube.ci",
		"v1alpha1.webhooks.git.kube.ci",
	}

	return Eventually(
		func() error {
			for _, apiService := range apiServices {
				apiservice, err := f.KAClient.ApiregistrationV1beta1().APIServices().Get(apiService, metav1.GetOptions{})
				if err != nil {
					return err
				}
				for _, cond := range apiservice.Status.Conditions {
					if cond.Type == kapi.Available && cond.Status == kapi.ConditionTrue && cond.Reason == "Passed" {
						return nil
					}
				}
				return fmt.Errorf("ApiService not ready yet")
			}
			return nil
		},
		time.Minute*5,
		time.Microsecond*10,
	)
}
