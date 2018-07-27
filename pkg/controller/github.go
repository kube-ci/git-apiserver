package controller

import (
	"encoding/json"
	"log"

	"github.com/TamalSaha/go-oneliners"
	"github.com/emicklei/go-restful"
	"github.com/google/go-github/github"
)

// https://github.com/cloud-ark/kubediscovery/blob/master/pkg/apiserver/apiserver.go
func (c *Controller) GetWebService(path string) *restful.WebService {
	log.Println("WS PATH:", path)

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
	log.Println("Event:", eventType)

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
		log.Println(err)
	}
	oneliners.PrettyJson(prEvent, "PullRequest")
}
