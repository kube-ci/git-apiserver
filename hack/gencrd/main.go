package main

import (
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/appscode/go/log"
	gort "github.com/appscode/go/runtime"
	crdutils "github.com/appscode/kutil/apiextensions/v1beta1"
	"github.com/appscode/kutil/openapi"
	"github.com/go-openapi/spec"
		crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kube-openapi/pkg/common"
	git_install "kube.ci/git-apiserver/apis/git/install"
	v1alpha1 "kube.ci/git-apiserver/apis/git/v1alpha1"
)

func generateCRDDefinitions() {
	v1alpha1.EnableStatusSubresource = true

	filename := gort.GOPath() + "/src/kube.ci/git-apiserver/apis/kubedb/v1alpha1/crds.yaml"
	os.Remove(filename)

	err := os.MkdirAll(filepath.Join(gort.GOPath(), "/src/kube.ci/git-apiserver/api/crds"), 0755)
	if err != nil {
		log.Fatal(err)
	}

	crds := []*crd_api.CustomResourceDefinition{
		v1alpha1.Repository{}.CustomResourceDefinition(),
		v1alpha1.Branch{}.CustomResourceDefinition(),
		v1alpha1.Tag{}.CustomResourceDefinition(),
		v1alpha1.PullRequest{}.CustomResourceDefinition(),
	}
	for _, crd := range crds {
		filename := filepath.Join(gort.GOPath(), "/src/kube.ci/git-apiserver/api/crds", crd.Spec.Names.Singular+".yaml")
		f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		crdutils.MarshallCrd(f, crd, "yaml")
		f.Close()
	}
}

func generateSwaggerJson() {
	var (
		Scheme = runtime.NewScheme()
		Codecs = serializer.NewCodecFactory(Scheme)
	)

	git_install.Install(Scheme)

	apispec, err := openapi.RenderOpenAPISpec(openapi.Config{
		Scheme: Scheme,
		Codecs: Codecs,
		Info: spec.InfoProps{
			Title:   "Kubeci",
			Version: "v0.1.0",
			Contact: &spec.ContactInfo{
				Name:  "AppsCode Inc.",
				URL:   "https://appscode.com",
				Email: "hello@appscode.com",
			},
			License: &spec.License{
				Name: "Apache 2.0",
				URL:  "https://www.apache.org/licenses/LICENSE-2.0.html",
			},
		},
		OpenAPIDefinitions: []common.GetOpenAPIDefinitions{
			v1alpha1.GetOpenAPIDefinitions,
		},
		Resources: []openapi.TypeInfo{
			{v1alpha1.SchemeGroupVersion, v1alpha1.ResourceRepositories, v1alpha1.ResourceKindRepository, true},
			{v1alpha1.SchemeGroupVersion, v1alpha1.ResourceBranches, v1alpha1.ResourceKindBranch, true},
			{v1alpha1.SchemeGroupVersion, v1alpha1.ResourceTags, v1alpha1.ResourceKindTag, true},
			{v1alpha1.SchemeGroupVersion, v1alpha1.ResourcePullRequests, v1alpha1.ResourceKindPullRequest, true},
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	filename := gort.GOPath() + "/src/kube.ci/git-apiserver/api/openapi-spec/swagger.json"
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		log.Fatal(err)
	}
	err = ioutil.WriteFile(filename, []byte(apispec), 0644)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	generateCRDDefinitions()
	generateSwaggerJson()
}
