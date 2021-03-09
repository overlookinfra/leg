module github.com/puppetlabs/leg/k8sutil

go 1.14

require (
	github.com/google/uuid v1.1.2
	github.com/puppetlabs/leg/errmap v0.1.0
	github.com/puppetlabs/leg/lifecycle v0.2.0
	github.com/puppetlabs/leg/stringutil v0.1.0
	github.com/puppetlabs/leg/timeutil v0.3.0
	github.com/rancher/remotedialer v0.2.5
	github.com/spf13/viper v1.4.0
	github.com/stretchr/testify v1.6.1
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.20.1
	k8s.io/client-go v0.19.2
	k8s.io/klog/v2 v2.4.0
	sigs.k8s.io/controller-runtime v0.7.0
	sigs.k8s.io/controller-tools v0.4.1
)
