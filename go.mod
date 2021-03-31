module github.com/grengojbo/k3s-operator

go 1.15

require (
	github.com/docker/go-connections v0.4.0
	github.com/go-logr/logr v0.4.0
	github.com/grengojbo/k3ctl v0.0.0-20210331171136-f146cd5cba9d
	github.com/liggitt/tabwriter v0.0.0-20181228230101-89fcab3d43de
	github.com/onsi/ginkgo v1.14.1
	github.com/onsi/gomega v1.10.2
	github.com/sirupsen/logrus v1.8.1
	github.com/spf13/cobra v1.1.3
	gopkg.in/yaml.v2 v2.4.0
	k8s.io/apimachinery v0.20.5
	k8s.io/client-go v1.5.2
	sigs.k8s.io/controller-runtime v0.8.3
)
