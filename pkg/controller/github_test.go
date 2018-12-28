package controller

import (
	"log"
	"os"
	"testing"
)

func testListGithubPRs(t *testing.T) {
	owner := "kube-ci"
	repo := "kubeci-gpig"
	token := ""

	prs, err := fetchGithubPRs(owner, repo, token)
	if err != nil {
		t.Error(err)
	}

	for _, pr := range prs {
		log.Println("PullRequest", *pr.Title)
	}
}

func testListGithubPRsWithAuth(t *testing.T) {
	owner := "tamalsaha"
	repo := "private-test-repo"
	token := os.Getenv("github-access-token")

	prs, err := fetchGithubPRs(owner, repo, token)
	if err != nil {
		t.Error(err)
	}

	for _, pr := range prs {
		log.Println("PullRequest", *pr.Title)
	}
}
