module github.com/kube-ci/git-apiserver

go 1.12

require (
	github.com/Azure/azure-sdk-for-go v29.0.0+incompatible // indirect
	github.com/OneOfOne/xxhash v1.2.5 // indirect
	github.com/alcortesm/tgz v0.0.0-20161220082320-9c5fe88206d7 // indirect
	github.com/anmitsu/go-shlex v0.0.0-20161002113705-648efa622239 // indirect
	github.com/appscode/go v0.0.0-20190424183524-60025f1135c9
	github.com/aws/aws-sdk-go v1.19.32 // indirect
	github.com/codeskyblue/go-sh v0.0.0-20190412065543-76bd3d59ff27
	github.com/dgryski/go-sip13 v0.0.0-20190329191031-25c5027a8c7b // indirect
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/docker/libtrust v0.0.0-20160708172513-aabc10ec26b7 // indirect
	github.com/emirpasic/gods v1.9.0 // indirect
	github.com/evanphx/json-patch v4.2.0+incompatible
	github.com/flosch/pongo2 v0.0.0-20190505152737-8914e1cf9164 // indirect
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/gliderlabs/ssh v0.1.1 // indirect
	github.com/go-ini/ini v1.42.0 // indirect
	github.com/go-openapi/spec v0.19.0
	github.com/golang/mock v1.3.1 // indirect
	github.com/google/go-github/v25 v25.0.4
	github.com/google/pprof v0.0.0-20190515194954-54271f7e092f // indirect
	github.com/gopherjs/gopherjs v0.0.0-20190430165422-3e4dfb77656c // indirect
	github.com/jbenet/go-context v0.0.0-20150711004518-d14ea06fba99 // indirect
	github.com/juju/errors v0.0.0-20190207033735-e65537c515d7 // indirect
	github.com/juju/loggo v0.0.0-20190212223446-d976af380377 // indirect
	github.com/juju/testing v0.0.0-20190429233213-dfc56b8c09fc // indirect
	github.com/kevinburke/ssh_config v0.0.0-20180830205328-81db2a75821e // indirect
	github.com/kisielk/errcheck v1.2.0 // indirect
	github.com/kr/pty v1.1.4 // indirect
	github.com/kube-ci/engine v0.0.0-20190517161517-3ce0fe7de607 // indirect
	github.com/mitchellh/go-homedir v1.0.0 // indirect
	github.com/onsi/ginkgo v1.8.0
	github.com/onsi/gomega v1.5.0
	github.com/pelletier/go-buffruneio v0.2.0 // indirect
	github.com/pkg/errors v0.8.1
	github.com/prometheus/tsdb v0.8.0 // indirect
	github.com/rogpeppe/fastuuid v1.1.0 // indirect
	github.com/shurcooL/sanitized_anchor_name v1.0.0 // indirect
	github.com/smartystreets/assertions v0.0.0-20190401211740-f487f9de1cd3 // indirect
	github.com/spaolacci/murmur3 v1.1.0 // indirect
	github.com/spf13/afero v1.2.2
	github.com/spf13/cobra v0.0.3
	github.com/spf13/pflag v1.0.3
	github.com/src-d/gcfg v1.4.0 // indirect
	github.com/stretchr/objx v0.2.0 // indirect
	github.com/xanzy/ssh-agent v0.2.0 // indirect
	golang.org/x/exp v0.0.0-20190510132918-efd6b22b2522 // indirect
	golang.org/x/image v0.0.0-20190516052701-61b8692d9a5c // indirect
	golang.org/x/lint v0.0.0-20190409202823-959b441ac422 // indirect
	golang.org/x/mobile v0.0.0-20190509164839-32b2708ab171 // indirect
	golang.org/x/oauth2 v0.0.0-20190402181905-9f3314589c9a
	golang.org/x/tools v0.0.0-20190517162734-1da8801a9502 // indirect
	gomodules.xyz/cert v1.0.0
	gopkg.in/ini.v1 v1.42.0 // indirect
	gopkg.in/src-d/go-billy.v4 v4.2.0 // indirect
	gopkg.in/src-d/go-git-fixtures.v3 v3.1.1 // indirect
	gopkg.in/src-d/go-git.v4 v4.7.0
	gopkg.in/warnings.v0 v0.1.2 // indirect
	honnef.co/go/tools v0.0.0-20190418001031-e561f6794a2a // indirect
	k8s.io/api v0.0.0-20190515023547-db5a9d1c40eb
	k8s.io/apiextensions-apiserver v0.0.0-20190515024537-2fd0e9006049
	k8s.io/apimachinery v0.0.0-20190515023456-b74e4c97951f
	k8s.io/apiserver v0.0.0-20190515064100-fc28ef5782df
	k8s.io/client-go v11.0.1-0.20190409021438-1a26190bd76a+incompatible
	k8s.io/kube-aggregator v0.0.0-20190515024249-81a6edcf70be
	k8s.io/kube-openapi v0.0.0-20190510232812-a01b7d5d6c22
	k8s.io/kubernetes v1.14.2
	kmodules.xyz/client-go v0.0.0-20190515205239-a16030cc2e50
	kmodules.xyz/webhook-runtime v0.0.0-20190508094945-962d01212c5b
)

replace (
	github.com/graymeta/stow => github.com/appscode/stow v0.0.0-20190506085026-ca5baa008ea3
	gopkg.in/robfig/cron.v2 => github.com/appscode/cron v0.0.0-20170717094345-ca60c6d796d4
	k8s.io/api => k8s.io/api v0.0.0-20190313235455-40a48860b5ab
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.0.0-20190315093550-53c4693659ed
	k8s.io/apimachinery => github.com/kmodules/apimachinery v0.0.0-20190508045248-a52a97a7a2bf
	k8s.io/apiserver => github.com/kmodules/apiserver v0.0.0-20190508082252-8397d761d4b5
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.0.0-20190314001948-2899ed30580f
	k8s.io/cloud-provider => k8s.io/cloud-provider v0.0.0-20190314002645-c892ea32361a
	k8s.io/component-base => k8s.io/component-base v0.0.0-20190314000054-4a91899592f4
	k8s.io/klog => k8s.io/klog v0.3.0
	k8s.io/kube-aggregator => k8s.io/kube-aggregator v0.0.0-20190314000639-da8327669ac5
	k8s.io/kube-openapi => k8s.io/kube-openapi v0.0.0-20190228160746-b3a7cee44a30
	k8s.io/metrics => k8s.io/metrics v0.0.0-20190314001731-1bd6a4002213
	k8s.io/utils => k8s.io/utils v0.0.0-20190221042446-c2654d5206da
)
