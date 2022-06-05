package module

import (
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	log "github.com/sirupsen/logrus"
)

// MakeInstallHaproxy
func MakeInstallHaproxy(ingress *k3sv1alpha1.Ingress, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallHaproxy"
	description := "Ingress Haproxy"
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
		overrides["defaultBackend.enabled"] = "true"
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

	options := k3sv1alpha1.HelmOptions{
		CreateNamespace: false,
		KubeconfigPath:  kubeConfigPath,
		Overrides:       overrides,
		Helm:            &release,
		Wait:            args.Wait,
		Verbose:         false,
		DryRun:          dryRun,
	}
	err = Helm3Upgrade(&options)

	return err
}
