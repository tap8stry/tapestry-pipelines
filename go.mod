module github.com/tapestry-pipelines

go 1.16

require (
	github.com/google/go-containerregistry v0.5.1
	github.com/nxadm/tail v1.4.8 // indirect
	github.com/peterbourgon/ff/v3 v3.1.0
	github.com/pkg/errors v0.9.1
	github.com/sigstore/cosign v0.6.0
	github.com/sigstore/k8s-manifest-sigstore v0.0.0-20210802081253-989d5862c336
	github.com/tektoncd/pipeline v0.28.1
	github.com/xlab/treeprint v1.1.0
	go.hein.dev/go-version v0.1.0
	go.uber.org/zap v1.19.0
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b
)

replace (
	github.com/sigstore/cosign => github.com/sigstore/cosign v1.0.1
	k8s.io/api => k8s.io/api v0.21.2
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.21.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.21.2
	k8s.io/apiserver => k8s.io/apiserver v0.21.2
	k8s.io/cli-runtime => k8s.io/cli-runtime v0.21.2
	k8s.io/client-go => k8s.io/client-go v0.21.2
	k8s.io/code-generator => k8s.io/code-generator v0.21.2
	k8s.io/kubectl => k8s.io/kubectl v0.21.2
	sigs.k8s.io/controller-runtime => sigs.k8s.io/controller-runtime v0.8.3
)
