package module

import (
	"fmt"

	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"
	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/k3s"
	t "github.com/grengojbo/k3ctl/pkg/types"
	log "github.com/sirupsen/logrus"
)

// MakeInstallNginx
func MakeInstallNginx(ingress *k3sv1alpha1.Ingress, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallNginx"
	description := "Ingress Nginx"
	update := false

	release, ok := k3sv1alpha1.FindRelease(args.Releases, ingress.Name)
	if !ok {
		return fmt.Errorf("[%s] is not release...", name)
	}

	// log.Debugf("[%s] name: %s disabled: %v status: %v", name, ingress.Name, ingress.Disabled, release.Status)

	if ingress.Disabled {
		log.Warnf("%s disabled...", description)
		return nil
	} else if len(release.Status) > 0 {
		if args.UpdateStrategy == "first" {
			log.Warnln("addons.options.UpdateStrategy IS SET first")
			return nil
		}
		log.Infof("Update %s...", description)
		update = true
	} else {
		log.Infof("Install %s...", description)
	}

	overrides := map[string]string{}

	if !update {
		overrides["installCRDs"] = "true"
	}

	if ingress.HostMode {
		log.Infof("Running in host networking mode")
		overrides["controller.hostNetwork"] = "true"
		overrides["controller.hostPort.enabled"] = "true"
		overrides["controller.service.type"] = "NodePort"
		overrides["dnsPolicy"] = "ClusterFirstWithHostNet"
		overrides["controller.kind"] = "DaemonSet"
	} else {
		// overrides["controller.service.externalTrafficPolicy"] = "Cluster"
		overrides["controller.service.externalTrafficPolicy"] = "Local"
		overrides["controller.config.use-proxy-protocol"] = "false"
		overrides["defaultBackend.enabled"] = "true"
		// TODO: добавить пользовательский defalt backend
		// overrides["defaultBackend.image.registry"] = "k8s.gcr.io"
		// overrides["defaultBackend.image.image"] = "defaultbackend-amd64"
		// overrides["defaultBackend.image.tag"] = "1.5"
	}

	option := k3sv1alpha1.HelmOptions{
		CreateNamespace: false,
		KubeconfigPath:  kubeConfigPath,
		Overrides:       overrides,
		Helm:            &release,
		Wait:            args.Wait,
		Verbose:         false,
		DryRun:          dryRun,
	}
	err = Helm3Upgrade(&option)

	return err
}

func MakeInstallNginxOld(kubeConfigPath string, dryRun bool, ingress *k3sv1alpha1.Ingress, args *k3sv1alpha1.HelmRelease) (err error) {
	updateRepo := true
	installed := true
	deleted := false
	overrides := map[string]string{}

	if len(ingress.Name) == 0 {
		ingress.Name = t.NginxDefaultName
	}
	if len(ingress.Namespace) == 0 {
		ingress.Namespace = t.NginxDefaultNamespace
	}

	log.Debugf("ingress: hostMode: %v updateSt %s", ingress.HostMode, args.UpdateStrategy)
	if args.UpdateStrategy != "none" {
		installed = true
	} else {
		ok, release := k3sv1alpha1.GetHelmRelease(ingress.Name, args.Releases)
		if ok {
			installed = false
			if ingress.Disabled {
				deleted = true
			} else {
				if len(ingress.Version) > 0 && ingress.Version != release.AppVersion {
					installed = true
				}
			}
			// log.Infof("Install Nginx Ingress controller %v", release.Revision)
		}
	}
	if deleted {
		log.Infoln("Deleted Nginx Ingress controller...")
		command := fmt.Sprintf(t.HelmDeleteCommand, ingress.Name, ingress.Namespace, kubeConfigPath)
		_, _, err := k3s.RunLocalCommand(command, false, dryRun)
		if err != nil {
			log.Errorf("[RunLocalCommand] %s\n%v", err.Error())
		}
	} else if installed {
		log.Infoln("Install Nginx Ingress controller...")
		if ingress.HostMode {
			log.Infof("Running in host networking mode")
			overrides["controller.hostNetwork"] = "true"
			overrides["controller.hostPort.enabled"] = "true"
			overrides["controller.service.type"] = "NodePort"
			overrides["dnsPolicy"] = "ClusterFirstWithHostNet"
			overrides["controller.kind"] = "DaemonSet"
		} else {
			// overrides["controller.service.externalTrafficPolicy"] = "Cluster"
			overrides["controller.service.externalTrafficPolicy"] = "Local"
			overrides["controller.config.use-proxy-protocol"] = "false"
			overrides["defaultBackend.enabled"] = "true"
			// TODO: добавить пользовательский defalt backend
			// overrides["defaultBackend.image.registry"] = "k8s.gcr.io"
			// overrides["defaultBackend.image.image"] = "defaultbackend-amd64"
			// overrides["defaultBackend.image.tag"] = "1.5"

		}

		// customFlags, _ := command.Flags().GetStringArray("set")

		// if err := mergeFlags(overrides, customFlags); err != nil {
		// 	return err
		// }
		nginxOptions := types.DefaultInstallOptions().
			WithNamespace(ingress.Namespace).
			WithHelmRepo(t.NginxHelmRepo).
			WithHelmURL(t.NginxHelmURL).
			WithOverrides(overrides).
			WithWait(args.Wait).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(nginxOptions)
		if err != nil {
			return err
		}

		log.Infof(nginxIngressInstallMsg)
	}
	return nil
}

const NginxIngressInfoMsg = `# If you're using a local environment such as "minikube" or "KinD",
# then try the inlets operator with "arkade install inlets-operator"

# If you're using a managed Kubernetes service, then you'll find
# your LoadBalancer's IP under "EXTERNAL-IP" via:

kubectl get svc ingress-nginx-controller

# Find out more at:
# https://github.com/kubernetes/ingress-nginx/tree/master/charts/ingress-nginx`

const nginxIngressInstallMsg = "ingress-nginx has been installed."

// const nginxIngressInstallMsg = `=======================================================================
// = ingress-nginx has been installed.                                   =
// =======================================================================` +
// 	"\n\n" + NginxIngressInfoMsg + "\n\n"
