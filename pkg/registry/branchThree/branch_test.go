package branchThree

import (
	"testing"

	"github.com/TamalSaha/go-oneliners"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apiserver/pkg/registry/generic"
)

func TestNewBranchRegistry(t *testing.T) {
	scheme := runtime.NewScheme()
	if ret, err := NewREST(scheme, generic.RESTOptions{}); err != nil {
		t.Errorf(err.Error())
	} else {
		oneliners.PrettyJson(ret)
	}
}
