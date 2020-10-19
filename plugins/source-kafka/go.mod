module knative.dev/client-contrib/plugins/source-kafka

go 1.14

require (
	github.com/maximilien/kn-source-pkg v0.5.0
	github.com/spf13/cobra v1.0.1-0.20200715031239-b95db644ed1c
	github.com/spf13/pflag v1.0.5
	gotest.tools v2.2.0+incompatible
	k8s.io/apimachinery v0.18.7-rc.0
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	knative.dev/client v0.17.0
	knative.dev/eventing-contrib v0.14.0
	knative.dev/pkg v0.0.0-20200828200807-2335e4d84a05
	knative.dev/test-infra v0.0.0-20200825022047-cb4bb218c5e5
)

// Temporary pinning certain libraries. Please check periodically, whether these are still needed
// ----------------------------------------------------------------------------------------------
replace (
	github.com/spf13/cobra => github.com/chmouel/cobra v0.0.0-20191021105835-a78788917390

	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
)
