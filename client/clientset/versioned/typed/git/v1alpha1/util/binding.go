package util

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/go/log"
	"github.com/appscode/kutil"
	"github.com/evanphx/json-patch"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
	cs "kube.ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1"
)

func CreateOrPatchBinding(c cs.GitV1alpha1Interface, meta metav1.ObjectMeta, transform func(binding *api.Binding) *api.Binding) (*api.Binding, kutil.VerbType, error) {
	cur, err := c.Bindings(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		log.Infof("Creating Binding %s/%s.", meta.Namespace, meta.Name)
		out, err := c.Bindings(meta.Namespace).Create(transform(&api.Binding{
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
	return PatchBinding(c, cur, transform)
}

func PatchBinding(c cs.GitV1alpha1Interface, cur *api.Binding, transform func(*api.Binding) *api.Binding) (*api.Binding, kutil.VerbType, error) {
	return PatchBindingObject(c, cur, transform(cur.DeepCopy()))
}

func PatchBindingObject(c cs.GitV1alpha1Interface, cur, mod *api.Binding) (*api.Binding, kutil.VerbType, error) {
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
	out, err := c.Bindings(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateBinding(c cs.GitV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.Binding) *api.Binding) (result *api.Binding, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Bindings(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Bindings(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		log.Errorf("Attempt %d failed to update Binding %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update Binding %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}

func TryUpdateBindingStatus(c cs.GitV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.Binding) *api.Binding) (result *api.Binding, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Bindings(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Bindings(cur.Namespace).UpdateStatus(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		log.Errorf("Attempt %d failed to update status of Binding %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update status of Binding %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
