package producer

import (
	"fmt"
	"log"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	git_v1alpha1 "kube.ci/git-apiserver/apis/git/v1alpha1"
)

type Producer struct {
	Repository string
	Url        string
	Secret     string
}

func (p *Producer) Run() {
	log.Println("Producer Run...")
	for i := 0; i < 100; i++ {
		branch := &git_v1alpha1.Branch{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("my-branch-%d", i),
				Namespace: "default",
				Labels: map[string]string{
					"repository": p.Repository,
				},
			},
			Spec: git_v1alpha1.BranchSpec{
				LastCommitHash: fmt.Sprintf("fake-hash-%d", i),
			},
		}
		log.Println(branch)
		// CreateOrUpdateBranch(branch)
		time.Sleep(time.Second * 10)
	}
}
