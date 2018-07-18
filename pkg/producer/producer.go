package producer

import (
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	repo_v1alpha1 "kube.ci/git-apiserver/apis/repositories/v1alpha1"
	"kube.ci/git-apiserver/pkg/registry/branchThree"
)

type Producer struct {
	Repository string
	Url        string
	Secret     string

	BranchRegistry *branchThree.REST
}

func (p *Producer) Run() {
	log.Println("Producer Run...")
	for i := 0; i < 100; i++ {
		branch := &repo_v1alpha1.Branch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("my-branch-%d", i),
				Namespace: "default",
				Labels: map[string]string{
					"repository": p.Repository,
				},
			},
			Status: repo_v1alpha1.BranchStatus{
				LastCommitHash: fmt.Sprintf("fake-hash-%d", i),
			},
		}
		p.BranchRegistry.CreateOrUpdateBranch(branch)
		time.Sleep(time.Second * 10)
	}
}
