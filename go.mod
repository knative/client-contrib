module github.com/knative/client-contrib

go 1.12

require (
	contrib.go.opencensus.io/exporter/ocagent v0.6.0 // indirect
	github.com/golang/glog v0.0.0-20160126235308-23def4e6c14b
	github.com/google/go-containerregistry v0.0.0-20200227193449-ba53fa10e72c // indirect
	github.com/knative/serving v0.12.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/openzipkin/zipkin-go v0.2.2 // indirect
	github.com/spf13/cobra v0.0.5
	github.com/spf13/viper v1.4.0
	github.com/tektoncd/pipeline v0.10.1
	k8s.io/api v0.17.0
	k8s.io/apiextensions-apiserver v0.17.0 // indirect
	k8s.io/apimachinery v0.17.1
	k8s.io/client-go v0.17.0
	knative.dev/pkg v0.0.0-20200206191432-9072f494c940
	knative.dev/serving v0.12.1-0.20200206201132-525b15d87dc1
)
