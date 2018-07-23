package git_repo

import (
	"log"
	"os"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
)

const (
	RemoteOrigin = "origin"
)

type GitRepository struct {
	Name string
	Url  string
	Path string

	Branches []Reference
	Tags     []Reference
}

type Reference struct {
	Name string
	Hash string
}

func GetGitRepository(name, url, path string) (GitRepository, error) {
	gitRepo := GitRepository{
		Name: name,
		Url:  url,
		Path: path,
	}

	// clone or fetch repo
	repo, err := getRepo(path, url)
	if err != nil {
		return GitRepository{}, err
	}

	// get origin branches
	originBranches, err := getOriginBranches(repo)
	if err != nil {
		return GitRepository{}, err
	}
	for _, reference := range originBranches {
		branch := getReference(reference)
		branch.Name = strings.TrimPrefix(branch.Name, "refs/heads/")
		gitRepo.Branches = append(gitRepo.Branches, branch)
	}

	// get tags
	tags, err := repo.Tags()
	if err != nil {
		return GitRepository{}, err
	}
	tags.ForEach(func(reference *plumbing.Reference) error {
		tag := getReference(reference)
		tag.Name = strings.TrimPrefix(tag.Name, "refs/tags/")
		gitRepo.Tags = append(gitRepo.Tags, tag)
		return nil
	})

	return gitRepo, nil
}

func getRepo(path, url string) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err != nil && err != git.ErrRepositoryNotExists {
		return nil, err
	}
	if err == git.ErrRepositoryNotExists {
		log.Println("Cloning repo...")
		repo, err = git.PlainClone(path, false, &git.CloneOptions{URL: url})
		if err != nil {
			return nil, err
		}
	} else {
		remote, err := repo.Remote(RemoteOrigin)
		if err != nil && err != git.ErrRemoteNotFound {
			return nil, err
		}
		if err == git.ErrRemoteNotFound || remote.Config().URLs[0] != url {
			log.Println("Remote changed, deleting old repo...")
			if err = os.RemoveAll(path); err != nil {
				return nil, err
			}
			log.Println("Cloning repo...")
			repo, err = git.PlainClone(path, false, &git.CloneOptions{URL: url})
			if err != nil {
				return nil, err
			}
		} else {
			log.Println("Fetching repo...")
			err = repo.Fetch(&git.FetchOptions{})
			if err != nil && err != git.NoErrAlreadyUpToDate {
				return nil, err
			}
		}
	}

	return repo, nil
}

func getOriginBranches(repo *git.Repository) ([]*plumbing.Reference, error) {
	var refBranches []*plumbing.Reference

	remote, err := repo.Remote(RemoteOrigin)
	if err != nil {
		return nil, err
	}

	refList, err := remote.List(&git.ListOptions{})
	if err != nil {
		return nil, err
	}

	refPrefix := "refs/heads/"
	for _, ref := range refList {
		refName := ref.Name().String()
		if !strings.HasPrefix(refName, refPrefix) {
			continue
		}
		refBranches = append(refBranches, ref)
	}

	return refBranches, nil
}

func getReference(reference *plumbing.Reference) Reference {
	return Reference{
		Name: reference.Name().String(),
		Hash: reference.Hash().String(),
	}
}
