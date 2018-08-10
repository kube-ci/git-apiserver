package controller

import (
	"log"
	"os"
	"testing"
)

func TestListGithubPRs(t *testing.T) {
	owner := "diptadas"
	repo := "kubeci-gpig"
	token := ""

	prs, err := listGithubPRs(owner, repo, token)
	if err != nil {
		t.Error(err)
	}

	for _, pr := range prs {
		log.Println("PullRequest", *pr.Title)
	}
}

func TestListGithubPRsWithAuth(t *testing.T) {
	owner := "tamalsaha"
	repo := "private-test-repo"
	token := os.Getenv("github-access-token")

	prs, err := listGithubPRs(owner, repo, token)
	if err != nil {
		t.Error(err)
	}

	for _, pr := range prs {
		log.Println("PullRequest", *pr.Title)
	}
}
