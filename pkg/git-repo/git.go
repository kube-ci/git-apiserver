package git_repo

import (
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

const (
	RemoteOrigin    = "origin"
	BranchRefPrefix = "refs/heads/"
	TagRefPrefix    = "refs/tags/"
)

type Repository struct {
	url  string
	path string
	auth *http.BasicAuth
	*git.Repository
}

type Reference struct {
	Name string
	Hash string
}

func New(url, path, token string) *Repository {
	repo := &Repository{
		url:        url,
		path:       path,
		auth:       nil,
		Repository: nil, // will be assigned in CloneOrFetch
	}

	// repository auth, nil if token is empty
	// token-auth not working, use basic-auth with token as password
	// https://github.com/src-d/go-git/issues/730
	if token != "" {
		repo.auth = &http.BasicAuth{
			Username: "token", // any string
			Password: token,
		}
	}

	return repo
}

// forceClone if repository crd changes
func (repo *Repository) CloneOrFetch(forceClone bool) error {
	var err error

	if forceClone {
		log.Println("Deleting old repository if exists...")
		if err = os.RemoveAll(repo.path); err != nil {
			return err
		}
	}

	// try to open repository from given path
	repo.Repository, err = git.PlainOpen(repo.path)
	if err != nil && err != git.ErrRepositoryNotExists {
		return err
	}

	if err == git.ErrRepositoryNotExists { // repository not exists, clone it
		log.Println("Cloning repo...")
		repo.Repository, err = git.PlainClone(repo.path, false, &git.CloneOptions{
			URL:  repo.url,
			Auth: repo.auth,
		})
		if err != nil {
			return err
		}
	} else { // repository exists, just fetch it
		log.Println("Fetching repo...")
		err = repo.Fetch(&git.FetchOptions{Auth: repo.auth})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return err
		}
	}

	return nil
}

// get origin branches
func (repo *Repository) GetBranches() ([]Reference, error) {
	var branches []Reference

	remote, err := repo.Remote(RemoteOrigin)
	if err != nil {
		return nil, err
	}

	refList, err := remote.List(&git.ListOptions{Auth: repo.auth})
	if err != nil {
		return nil, err
	}

	for _, ref := range refList {
		if strings.HasPrefix(ref.Name().String(), BranchRefPrefix) {
			branches = append(branches, Reference{
				Name: strings.TrimPrefix(ref.Name().String(), BranchRefPrefix),
				Hash: ref.Hash().String(),
			})
		}
	}

	return branches, nil
}

func (repo *Repository) GetTags() ([]Reference, error) {
	var tags []Reference

	refList, err := repo.Tags()
	if err != nil {
		return nil, err
	}

	refList.ForEach(func(ref *plumbing.Reference) error {
		tags = append(tags, Reference{
			Name: strings.TrimPrefix(ref.Name().String(), TagRefPrefix),
			Hash: ref.Hash().String(),
		})
		return nil
	})

	return tags, nil
}
