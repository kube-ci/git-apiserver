package util

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/go/log"
	jsonpatch "github.com/evanphx/json-patch"
	api "github.com/kube-ci/git-apiserver/apis/git/v1alpha1"
	cs "github.com/kube-ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	kutil "kmodules.xyz/client-go"
)

func CreateOrPatchBranch(c cs.GitV1alpha1Interface, meta metav1.ObjectMeta, transform func(branch *api.Branch) *api.Branch) (*api.Branch, kutil.VerbType, error) {
	cur, err := c.Branches(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		log.Infof("Creating Branch %s/%s.", meta.Namespace, meta.Name)
		out, err := c.Branches(meta.Namespace).Create(transform(&api.Branch{
			TypeMeta: metav1.TypeMeta{
				Kind:       "Branch",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchBranch(c, cur, transform)
}

func PatchBranch(c cs.GitV1alpha1Interface, cur *api.Branch, transform func(*api.Branch) *api.Branch) (*api.Branch, kutil.VerbType, error) {
	return PatchBranchObject(c, cur, transform(cur.DeepCopy()))
}

func PatchBranchObject(c cs.GitV1alpha1Interface, cur, mod *api.Branch) (*api.Branch, kutil.VerbType, error) {
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
	log.Infof("Patching Branch %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.Branches(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdateBranch(c cs.GitV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.Branch) *api.Branch) (result *api.Branch, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.Branches(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.Branches(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		log.Errorf("Attempt %d failed to update Branch %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update Branch %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}
