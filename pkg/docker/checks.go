package docker

const (
	ACRegistry        = "kubeci"
	ImageGitAPIServer = "git-apiserver"
)

type Docker struct {
	Registry, Image, Tag string
}

func (docker Docker) ToContainerImage() string {
	return docker.Registry + "/" + docker.Image + ":" + docker.Tag
}
