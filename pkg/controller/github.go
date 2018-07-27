package controller

import (
	"encoding/json"

	"github.com/appscode/go/log"
	"github.com/emicklei/go-restful"
	"github.com/google/go-github/github"
	"k8s.io/apimachinery/pkg/labels"
)

// https://github.com/cloud-ark/kubediscovery/blob/master/pkg/apiserver/apiserver.go
func (c *Controller) GetWebService(path string) *restful.WebService {
	log.Infoln("WS PATH:", path)

	ws := new(restful.WebService).Path(path)
	//ws.Consumes("*/*")
	//ws.Produces(restful.MIME_JSON, restful.MIME_XML)
	//ws.ApiVersion("foocontroller.k8s.io/v1alpha1")

	githubPath := "/github"
	ws.Route(ws.POST(githubPath).To(c.githubEventHandler))

	return ws
}

func (c *Controller) githubEventHandler(request *restful.Request, response *restful.Response) {
	eventType := request.Request.Header.Get("X-GitHub-Event")
	log.Infoln("Event:", eventType)

	switch eventType {
	case "pull_request":
		c.githubPRHandler(request, response)
	default:
		return
	}

	response.Write([]byte("hello world"))
}

func (c *Controller) githubPRHandler(request *restful.Request, response *restful.Response) {
	var prEvent github.PullRequestEvent
	decoder := json.NewDecoder(request.Request.Body)
	if err := decoder.Decode(&prEvent); err != nil {
		log.Errorln(err)
	}
	// oneliners.PrettyJson(prEvent, "PullRequest")

	// find matching repository
	repositories, err := c.repoLister.List(labels.Everything())
	if err != nil {
		log.Errorln(err)
	}

	for _, repository := range repositories {
		if prEvent.Action != nil &&
			prEvent.Repo != nil &&
			prEvent.Repo.CloneURL != nil &&
			prEvent.PullRequest != nil &&
			repository.Spec.Url == *prEvent.Repo.CloneURL {
			log.Infof("PR event for repository %s/%s", repository.Namespace, repository.Name)
			if *prEvent.Action != "deleted" {
				// create or patch PullRequest CRD
			} else {
				// delete PullRequest CRD
			}
		}
	}
}
