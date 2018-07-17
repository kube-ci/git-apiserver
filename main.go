package main

import (
	"os"
	"runtime"

	"github.com/appscode/go/log"
	logs "github.com/appscode/go/log/golog"
	_ "k8s.io/client-go/kubernetes/fake"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	_ "kube.ci/git-apiserver/client/clientset/versioned/fake"
	"kube.ci/git-apiserver/pkg/cmds"
)

func main() {
	logs.InitLogs()
	defer logs.FlushLogs()

	if len(os.Getenv("GOMAXPROCS")) == 0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}

	if err := cmds.NewRootCmd().Execute(); err != nil {
		log.Fatalln("Error in git-apiserver Main:", err)
	}
	log.Infoln("Exiting git-apiserver Main")
	os.Exit(0)
}
