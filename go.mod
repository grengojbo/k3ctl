module github.com/grengojbo/k3ctl

go 1.15

require (
	github.com/alexellis/go-execute v0.0.0-20200124154445-8697e4e28c5e
	github.com/alexellis/k3sup v0.0.0-20210413182206-de10fc701f46
	github.com/appleboy/easyssh-proxy v1.3.9
	github.com/docker/go-connections v0.4.0
	github.com/docker/go-units v0.4.0
	github.com/go-logr/logr v1.2.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/onsi/ginkgo v1.16.5
	github.com/onsi/gomega v1.18.1 // indirect
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/api v0.23.0 // indirect
	k8s.io/apimachinery v0.23.0
	k8s.io/client-go v0.23.0
	sigs.k8s.io/cluster-api v0.3.15
	sigs.k8s.io/controller-runtime v0.11.0
)
