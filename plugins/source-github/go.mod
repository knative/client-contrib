module github.com/maximilien/kn-source-github

go 1.13

require (
	github.com/maximilien/kn-source-pkg v0.5.0
	knative.dev/eventing-contrib v0.14.0
	knative.dev/test-infra v0.0.0-20200825022047-cb4bb218c5e5
)

replace (
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
)
