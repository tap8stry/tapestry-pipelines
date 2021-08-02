module github.com/tapestry-pipelines

go 1.16

require (
	github.com/GoogleContainerTools/skaffold v1.28.1 // indirect
	github.com/alvaroloes/enumer v1.1.2 // indirect
	github.com/coreos/go-systemd v0.0.0-20190620071333-e64a0ec8b42a // indirect
	github.com/ghodss/yaml v1.0.0 // indirect
	github.com/golangci/go-tools v0.0.0-20190318055746-e32c54105b7c // indirect
	github.com/golangci/gosec v0.0.0-20190211064107-66fb7fc33547 // indirect
	github.com/google/go-containerregistry v0.5.1
	github.com/google/monologue v0.0.0-20190606152607-4b11a32b5934 // indirect
	github.com/google/trillian-examples v0.0.0-20190603134952-4e75ba15216c // indirect
	github.com/jinzhu/copier v0.3.2 // indirect
	github.com/letsencrypt/pkcs11key v2.0.1-0.20170608213348-396559074696+incompatible // indirect
	github.com/mattn/go-sqlite3 v1.10.0 // indirect
	github.com/oliveagle/jsonpath v0.0.0-20180606110733-2e52cf6e6852 // indirect
	github.com/onsi/ginkgo v1.15.0 // indirect
	github.com/onsi/gomega v1.11.0 // indirect
	github.com/peterbourgon/ff/v3 v3.1.0
	github.com/pkg/errors v0.9.1
	github.com/prometheus/common v0.29.0 // indirect
	github.com/r3labs/diff v1.1.0 // indirect
	github.com/sigstore/cosign v0.6.0
	github.com/sigstore/fulcio v0.1.1 // indirect
	github.com/sigstore/k8s-manifest-sigstore v0.0.0-20210802081253-989d5862c336
	github.com/sirupsen/logrus v1.8.1 // indirect
	github.com/spf13/cobra v1.2.1 // indirect
	github.com/tektoncd/pipeline v0.26.0
	github.com/xlab/treeprint v1.1.0
	go.etcd.io/etcd v3.3.13+incompatible // indirect
	go.hein.dev/go-version v0.1.0
	go.uber.org/zap v1.18.1
	golang.org/x/term v0.0.0-20210615171337-6886f2dfbf5b // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	k8s.io/cli-runtime v0.21.2 // indirect
	k8s.io/kube-openapi v0.0.0-20210305001622-591a79e4bda7 // indirect
	k8s.io/kubectl v0.19.4 // indirect
	sigs.k8s.io/controller-runtime v0.0.0-00010101000000-000000000000 // indirect
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
