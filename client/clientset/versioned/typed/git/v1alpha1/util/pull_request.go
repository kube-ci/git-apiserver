package util

import (
	"encoding/json"
	"fmt"

	"github.com/appscode/kutil"
	"github.com/evanphx/json-patch"
	"github.com/golang/glog"
	kerr "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/wait"
	api "kube.ci/git-apiserver/apis/git/v1alpha1"
	cs "kube.ci/git-apiserver/client/clientset/versioned/typed/git/v1alpha1"
)

func CreateOrPatchPullRequest(c cs.GitV1alpha1Interface, meta metav1.ObjectMeta, transform func(pullRequest *api.PullRequest) *api.PullRequest) (*api.PullRequest, kutil.VerbType, error) {
	cur, err := c.PullRequests(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
	if kerr.IsNotFound(err) {
		glog.V(3).Infof("Creating PullRequest %s/%s.", meta.Namespace, meta.Name)
		out, err := c.PullRequests(meta.Namespace).Create(transform(&api.PullRequest{
			TypeMeta: metav1.TypeMeta{
				Kind:       "PullRequest",
				APIVersion: api.SchemeGroupVersion.String(),
			},
			ObjectMeta: meta,
		}))
		return out, kutil.VerbCreated, err
	} else if err != nil {
		return nil, kutil.VerbUnchanged, err
	}
	return PatchPullRequest(c, cur, transform)
}

func PatchPullRequest(c cs.GitV1alpha1Interface, cur *api.PullRequest, transform func(*api.PullRequest) *api.PullRequest) (*api.PullRequest, kutil.VerbType, error) {
	return PatchPullRequestObject(c, cur, transform(cur.DeepCopy()))
}

func PatchPullRequestObject(c cs.GitV1alpha1Interface, cur, mod *api.PullRequest) (*api.PullRequest, kutil.VerbType, error) {
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
	glog.V(3).Infof("Patching PullRequest %s/%s with %s.", cur.Namespace, cur.Name, string(patch))
	out, err := c.PullRequests(cur.Namespace).Patch(cur.Name, types.MergePatchType, patch)
	return out, kutil.VerbPatched, err
}

func TryUpdatePullRequest(c cs.GitV1alpha1Interface, meta metav1.ObjectMeta, transform func(*api.PullRequest) *api.PullRequest) (result *api.PullRequest, err error) {
	attempt := 0
	err = wait.PollImmediate(kutil.RetryInterval, kutil.RetryTimeout, func() (bool, error) {
		attempt++
		cur, e2 := c.PullRequests(meta.Namespace).Get(meta.Name, metav1.GetOptions{})
		if kerr.IsNotFound(e2) {
			return false, e2
		} else if e2 == nil {
			result, e2 = c.PullRequests(cur.Namespace).Update(transform(cur.DeepCopy()))
			return e2 == nil, nil
		}
		glog.Errorf("Attempt %d failed to update PullRequest %s/%s due to %v.", attempt, cur.Namespace, cur.Name, e2)
		return false, nil
	})

	if err != nil {
		err = fmt.Errorf("failed to update PullRequest %s/%s after %d attempts due to %v", meta.Namespace, meta.Name, attempt, err)
	}
	return
}