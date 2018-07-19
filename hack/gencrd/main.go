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
	"github.com/golang/glog"
	crd_api "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/kube-openapi/pkg/common"
	git_install "kube.ci/git-apiserver/apis/git/install"
	git_v1alpha1 "kube.ci/git-apiserver/apis/git/v1alpha1"
)

func generateCRDDefinitions() {
	filename := gort.GOPath() + "/src/kube.ci/git-apiserver/apis/git/v1alpha1/crds.yaml"

	f, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	crds := []*crd_api.CustomResourceDefinition{
		git_v1alpha1.Repository{}.CustomResourceDefinition(),
	}
	for _, crd := range crds {
		err = crdutils.MarshallCrd(f, crd, "yaml")
		if err != nil {
			log.Fatal(err)
		}
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
			git_v1alpha1.GetOpenAPIDefinitions,
		},
		Resources: []openapi.TypeInfo{
			{git_v1alpha1.SchemeGroupVersion, git_v1alpha1.ResourceRepositories, git_v1alpha1.ResourceKindRepository, true},
			{git_v1alpha1.SchemeGroupVersion, git_v1alpha1.ResourceBranches, git_v1alpha1.ResourceKindBranch, true},
		},
	})
	if err != nil {
		glog.Fatal(err)
	}

	filename := gort.GOPath() + "/src/kube.ci/git-apiserver/api/openapi-spec/swagger.json"
	err = os.MkdirAll(filepath.Dir(filename), 0755)
	if err != nil {
		glog.Fatal(err)
	}
	err = ioutil.WriteFile(filename, []byte(apispec), 0644)
	if err != nil {
		glog.Fatal(err)
	}
}

func main() {
	generateCRDDefinitions()
	generateSwaggerJson()
}
