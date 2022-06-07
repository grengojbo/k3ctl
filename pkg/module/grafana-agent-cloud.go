package module

import (
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

// MakeInstallGrafanaAgentCloud
func MakeInstallGrafanaAgentCloud(addons *k3sv1alpha1.Monitoring, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallGrafanaAgentCloud"
	description := "Grafana Agent Cloud"
	// update := false

	release, ok := k3sv1alpha1.FindRelease(args.Releases, types.GrafanaAgentCloudDefaultName)
	if !ok {
		return fmt.Errorf("[%s] is not release...", name)
	}

	log.Debugf("[%s] name: %s disabled: %v status: %v (cluster: %s)", name, addons.Name, addons.Disabled, release.Status, args.ClusterName)

	if len(addons.ValuesFile) == 0 {
		addons.ValuesFile, err = util.CheckExitValueFile(args.ClusterName, release.Name)
		if err != nil {
			log.Errorf("IS NOT SET addons.monitoring.valuesFile OR %s", err.Error())
			return nil
		}
	} else {
		if err = util.CheckExitFile(addons.ValuesFile); err != nil {
			log.Errorf("IS NOT file: addons.monitoring.valuesFile=%s", addons.ValuesFile)
			return nil
		}
	}
	release.ValuesFile = addons.ValuesFile
	release.DependencyUpdate = true

	if addons.Disabled {
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

	overrides := map[string]string{}

	// if !update {
	// 	overrides["installCRDs"] = "true"
	// }

	// if ingress.HostMode {
	// 	log.Infof("Running in host networking mode")
	// 	overrides["controller.hostNetwork"] = "true"
	// 	overrides["controller.hostPort.enabled"] = "true"
	// 	overrides["controller.service.type"] = "NodePort"
	// 	overrides["dnsPolicy"] = "ClusterFirstWithHostNet"
	// 	overrides["controller.kind"] = "DaemonSet"
	// } else {
	// 	// overrides["controller.service.externalTrafficPolicy"] = "Cluster"
	// 	overrides["controller.service.externalTrafficPolicy"] = "Local"
	// 	overrides["controller.config.use-proxy-protocol"] = "false"
	// }

	overrides["clusterName"] = args.ClusterName

	// if len(ingress.DefaultBackend.Registry) > 0 {
	// 	overrides["defaultBackend.image.registry"] = ingress.DefaultBackend.Registry
	// }
	// if len(ingress.DefaultBackend.Image) > 0 {
	// 	overrides["defaultBackend.image.image"] = ingress.DefaultBackend.Image
	// }
	// if len(ingress.DefaultBackend.Tag) > 0 {
	// 	overrides["defaultBackend.image.tag"] = ingress.DefaultBackend.Tag
	// }
	// // overrides["defaultBackend.image.registry"] = "k8s.gcr.io"
	// // overrides["defaultBackend.image.image"] = "defaultbackend-amd64"
	// // overrides["defaultBackend.image.tag"] = "1.5"

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
