package git_repo

import (
	"log"
	"os"
	"testing"
)

func TestFetch(t *testing.T) {
	url := "https://github.com/kube-ci/kubeci-gpig.git"
	token := ""

	repo, err := Fetch(url, token)
	if err != nil {
		t.Error(err)
	}

	log.Println("Branches")
	for _, branch := range repo.Branches {
		log.Println(branch)
	}

	log.Println("Tags")
	for _, tag := range repo.Tags {
		log.Println(tag)
	}
}

func testFetchPrivate(t *testing.T) {
	url := "https://github.com/tamalsaha/private-test-repo.git"
	token := os.Getenv("github-access-token")

	repo, err := Fetch(url, token)
	if err != nil {
		t.Error(err)
	}

	log.Println("Branches")
	for _, branch := range repo.Branches {
		log.Println(branch)
	}

	log.Println("Tags")
	for _, tag := range repo.Tags {
		log.Println(tag)
	}
}
