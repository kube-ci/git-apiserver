package git_repo

import (
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/config"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"gopkg.in/src-d/go-git.v4/storage/memory"
)

const (
	RemoteOrigin    = "origin"
	BranchRefPrefix = "refs/heads/"
	TagRefPrefix    = "refs/tags/"
)

type Repository struct {
	Branches []Reference
	Tags     []Reference
}

type Reference struct {
	Name string
	Hash string
}

func Fetch(url, token string) (Repository, error) {
	// repository auth, nil if token is empty
	// token-auth not working, use basic-auth with token as password
	// https://github.com/src-d/go-git/issues/730
	var auth *http.BasicAuth
	if token != "" {
		auth = &http.BasicAuth{
			Username: "token", // any string
			Password: token,
		}
	}

	// get remote refs (cmd: $ git ls-remote)
	refList, err := getRefs(url, auth)
	if err != nil {
		return Repository{}, err
	}

	var repo Repository

	for _, ref := range refList {
		if strings.HasPrefix(ref.Name().String(), BranchRefPrefix) {
			repo.Branches = append(repo.Branches, Reference{
				Name: strings.TrimPrefix(ref.Name().String(), BranchRefPrefix),
				Hash: ref.Hash().String(),
			})
		} else if strings.HasPrefix(ref.Name().String(), TagRefPrefix) {
			repo.Tags = append(repo.Tags, Reference{
				Name: strings.TrimPrefix(ref.Name().String(), TagRefPrefix),
				Hash: ref.Hash().String(),
			})
		}
	}

	return repo, nil
}

func getRefs(url string, auth transport.AuthMethod) ([]*plumbing.Reference, error) {
	repo := &git.Repository{
		Storer: memory.NewStorage(),
	}
	remote, err := repo.CreateRemote(&config.RemoteConfig{
		Name: RemoteOrigin,
		URLs: []string{url},
	})
	if err != nil {
		return nil, err
	}
	return remote.List(&git.ListOptions{Auth: auth})
}
