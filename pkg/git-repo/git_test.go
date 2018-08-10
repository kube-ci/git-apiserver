package git_repo

import (
	"log"
	"os"
	"testing"
)

func TestGetBranches(t *testing.T) {
	url := "https://github.com/diptadas/kubeci-gpig.git"
	path := "/tmp/my-repo"
	token := ""

	repo := New(url, path, token)
	if err := repo.CloneOrFetch(true); err != nil {
		t.Error(err)
	}
	if err := repo.CloneOrFetch(false); err != nil { // should fetch instead of cloning
		t.Error(err)
	}

	branches, err := repo.GetBranches()
	if err != nil {
		t.Error(err)
	}

	for _, branch := range branches {
		log.Println(branch)
	}
}

func TestGetTags(t *testing.T) {
	url := "https://github.com/diptadas/kubeci-gpig.git"
	path := "/tmp/my-repo"
	token := ""

	os.RemoveAll(path)

	repo := New(url, path, token)
	if err := repo.CloneOrFetch(); err != nil {
		t.Error(err)
	}

	tags, err := repo.GetTags()
	if err != nil {
		t.Error(err)
	}

	for _, tag := range tags {
		log.Println(tag)
	}
}

func TestGetBranchesWithAuth(t *testing.T) {
	url := "https://github.com/tamalsaha/private-test-repo.git"
	path := "/tmp/my-repo"
	token := os.Getenv("github-access-token")

	os.RemoveAll(path)

	repo := New(url, path, token)
	if err := repo.CloneOrFetch(); err != nil {
		t.Error(err)
	}

	branches, err := repo.GetBranches()
	if err != nil {
		t.Error(err)
	}

	for _, branch := range branches {
		log.Println(branch)
	}
}
