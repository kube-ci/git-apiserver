package v1alpha1

import (
	"encoding/json"

	"github.com/appscode/go/log"
	"github.com/appscode/kutil"
	"github.com/evanphx/json-patch"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
)

func (c *bindings) CreateOrPatch(meta metav1.ObjectMeta, transform func(binding *api.Binding) *api.Binding) (*api.Binding, kutil.VerbType, error) {
	cur, err := c.Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		log.Infof("Creating Binding %s/%s.", meta.Namespace, meta.Name)
		out, err := c.Create(transform(&api.Binding{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Binding",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return c.patchBindingObject(cur, transform(cur.DeepCopy()))
}

func (c *bindings) patchBindingObject(cur, mod *api.Binding) (*api.Binding, kutil.VerbType, error) {
	curJson, err := json.Marshal(cur)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	modJson, err := json.Marshal(mod)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}

	patch, err := jsonpatch.CreateMergePatch(curJson, modJson)
	if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	if len(patch) == 0 || string(patch) == "{}" {
		return cur, kutil.VerbUnchanged, nil
	}

	log.Infof("Patching Binding %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}
