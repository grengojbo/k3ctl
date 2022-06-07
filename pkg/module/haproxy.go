package module

import (
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	log "github.com/sirupsen/logrus"
)

// MakeInstallHaproxy
func MakeInstallHaproxy(ingress *k3sv1alpha1.Ingress, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallHaproxy"
	description := "Ingress Haproxy"
	// update := false

	release, ok := k3sv1alpha1.FindRelease(args.Releases, types.HaproxyDefaultName)
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
		// update = true
	} else {
		log.Infof("Install %s...", description)
	}

	// https://github.com/haproxytech/helm-charts/blob/main/kubernetes-ingress/values.yaml
	overrides := map[string]string{}

	// if !update {
	// 	overrides["installCRDs"] = "true"
	// }

	if ingress.HostMode {
		log.Infof("Running in host networking mode")
		overrides["controller.service.enabled"] = "false"
		overrides["dnsPolicy"] = "ClusterFirstWithHostNet"
		overrides["controller.kind"] = "DaemonSet"
		overrides["controller.daemonset.useHostPort"] = "true"
		overrides["controller.daemonset.useHostNetwork"] = "true"
		// overrides["controller."] = ""
		// overrides["controller."] = ""
	} else {
		overrides["controller.service.type"] = "LoadBalancer"
		// overrides["controller.service.externalIPs[]"] = ""
		// LoadBalancer IP
		// ref: https://kubernetes.io/docs/concepts/services-networking/service/#loadbalancer
		// overrides["controller.service.loadBalancerIP"] = ""
		// overrides["controller.service.externalTrafficPolicy"] = "Cluster"
		overrides["controller.service.externalTrafficPolicy"] = "Local"
		// overrides["controller.config.use-proxy-protocol"] = "false"
	}

	overrides["controller.replicaCount"] = "1"

	overrides["defaultBackend.enabled"] = "true"
	overrides["defaultBackend.replicaCount"] = "1"

	// overrides["defaultBackend.image.registry"] = "k8s.gcr.io"
	// overrides["defaultBackend.image.image"] = "defaultbackend-amd64"
	// overrides["defaultBackend.image.tag"] = "1.5"
	if len(ingress.DefaultBackend.Registry) > 0 {
		overrides["defaultBackend.image.registry"] = ingress.DefaultBackend.Registry
	}
	if len(ingress.DefaultBackend.Image) > 0 {
		overrides["defaultBackend.image.image"] = ingress.DefaultBackend.Image
	}
	if len(ingress.DefaultBackend.Tag) > 0 {
		overrides["defaultBackend.image.tag"] = ingress.DefaultBackend.Tag
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
