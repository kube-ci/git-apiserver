package git_repo

import (
	"log"
	"testing"

	"gopkg.in/src-d/go-git.v4/plumbing/transport/http"
)

func TestGetGitRepository(t *testing.T) {
	url := "https://github.com/diptadas/kubeci-gpig.git"
	path := "/tmp/my-repo"

	gitRepo, err := GetGitRepository(url, path, nil)
	if err != nil {
		t.Error(err)
	}

	for _, branch := range gitRepo.Branches {
		log.Println("Branches", branch)
	}
	for _, tag := range gitRepo.Tags {
		log.Println("Tags", tag)
	}
}

func TestGetGitRepositoryWithAuth(t *testing.T) {
	url := "https://github.com/diptadas/kubeci-gpig.git"
	path := "/tmp/my-repo"

	gitRepo, err := GetGitRepository(url, path, &http.TokenAuth{
		Token: "...",
	})
	if err != nil {
		t.Error(err.Error())
	}

	for _, branch := range gitRepo.Branches {
		log.Println("Branches", branch)
	}
	for _, tag := range gitRepo.Tags {
		log.Println("Tags", tag)
	}
}
