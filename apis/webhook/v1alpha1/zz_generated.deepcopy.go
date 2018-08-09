// +build !ignore_autogenerated

/*
Copyright 2018 The KubeCI Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by deepcopy-gen. DO NOT EDIT.

package v1alpha1

import (
	github "github.com/google/go-github/github"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *GithubEvent) DeepCopyInto(out *GithubEvent) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	if in.Action != nil {
		in, out := &in.Action, &out.Action
		if *in == nil {
			*out = nil
		} else {
			*out = new(string)
			**out = **in
		}
	}
	if in.Repo != nil {
		in, out := &in.Repo, &out.Repo
		if *in == nil {
			*out = nil
		} else {
			*out = new(github.Repository)
			**out = **in
		}
	}
	if in.Sender != nil {
		in, out := &in.Sender, &out.Sender
		if *in == nil {
			*out = nil
		} else {
			*out = new(github.User)
			**out = **in
		}
	}
	if in.Issue != nil {
		in, out := &in.Issue, &out.Issue
		if *in == nil {
			*out = nil
		} else {
			*out = new(github.Issue)
			**out = **in
		}
	}
	if in.PullRequest != nil {
		in, out := &in.PullRequest, &out.PullRequest
		if *in == nil {
			*out = nil
		} else {
			*out = new(github.PullRequest)
			**out = **in
		}
	}
	return
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new GithubEvent.
func (in *GithubEvent) DeepCopy() *GithubEvent {
	if in == nil {
		return nil
	}
	out := new(GithubEvent)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *GithubEvent) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}