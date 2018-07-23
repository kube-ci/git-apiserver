package git_repo

import (
	"log"
	"testing"
)

func TestGetGitRepository(t *testing.T) {
	url := "https://github.com/appscode/voyager.git"
	path := "/tmp/my-repo"

	gitRepo, err := GetGitRepository("my-repo", url, path)
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
