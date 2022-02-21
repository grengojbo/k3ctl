package module

import (
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"
	log "github.com/sirupsen/logrus"
)

// MakeInstallNginx
func MakeInstallNginx(kubeConfigPath string) (err error) {
	namespace := "ingress-nginx"
	wait := true
	updateRepo := true
	hostMode := false

		overrides := map[string]string{}

	if hostMode {
		log.Infof("Running in host networking mode")
			overrides["controller.hostNetwork"] = "true"
			overrides["controller.hostPort.enabled"] = "true"
			overrides["controller.service.type"] = "NodePort"
			overrides["dnsPolicy"] = "ClusterFirstWithHostNet"
			overrides["controller.kind"] = "DaemonSet"
		}

	nginxOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("ingress-nginx/ingress-nginx").
			WithHelmURL("https://kubernetes.github.io/ingress-nginx").
			WithOverrides(overrides).
			WithWait(wait).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(nginxOptions)

		if err != nil {
			return err
		}

		log.Infof(nginxIngressInstallMsg)

	return nil
}

const NginxIngressInfoMsg = `# If you're using a local environment such as "minikube" or "KinD",
# then try the inlets operator with "arkade install inlets-operator"

# If you're using a managed Kubernetes service, then you'll find
# your LoadBalancer's IP under "EXTERNAL-IP" via:

kubectl get svc ingress-nginx-controller

# Find out more at:
# https://github.com/kubernetes/ingress-nginx/tree/master/charts/ingress-nginx`

const nginxIngressInstallMsg = `=======================================================================
= ingress-nginx has been installed.                                   =
=======================================================================` +
	"\n\n" + NginxIngressInfoMsg + "\n\n"