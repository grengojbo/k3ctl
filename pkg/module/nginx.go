package module

import (
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

// NginxSettings
func NginxSettings(addons *k3sv1alpha1.Ingress, clusterName string) (release k3sv1alpha1.HelmInterfaces) {
	repo := k3sv1alpha1.HelmRepo{
		Name: types.NginxHelmRepoNane,
		Repo: types.NginxHelmRepo,
		Url:  types.NginxHelmURL,
	}
	if len(addons.Repo.Name) > 0 {
		repo.Name = addons.Repo.Name
	}
	if len(addons.Repo.Repo) > 0 {
		repo.Repo = addons.Repo.Repo
	}
	if len(addons.Repo.Url) > 0 {
		repo.Url = addons.Repo.Url
	}

	if addons.Disabled {
		release.Deleted = true
	}
	if len(addons.Name) == 0 {
		addons.Name = types.NginxDefaultName
	}
	if len(addons.Namespace) == 0 {
		addons.Namespace = types.NginxDefaultNamespace
	}
	if len(addons.Version) > 0 {
		release.Version = addons.Version
	}
	if len(addons.Values) > 0 {
		release.Values = addons.Values
	}
	if len(addons.ValuesFile) > 0 {
		release.ValuesFile = addons.ValuesFile
	}

	//  All Settings
	release.Name = addons.Name
	release.Namespace = addons.Namespace
	release.Repo = repo.Repo

	addons.Repo = repo
	return release
}

// MakeInstallNginx
func MakeInstallNginx(addons *k3sv1alpha1.Ingress, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallNginx"
	description := "Ingress Nginx"
	update := false

	release, ok := k3sv1alpha1.FindRelease(args.Releases, types.NginxDefaultName)
	if !ok {
		return fmt.Errorf("[%s] is not release...", name)
	}

	// log.Debugf("[%s] name: %s disabled: %v status: %v (cluster: %s)", name, addons.Name, addons.Disabled, release.Status, args.ClusterName)

	if addons.Disabled {
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

	if len(addons.ValuesFile) > 0 {
		if err = util.CheckExitFile(addons.ValuesFile); err != nil {
			log.Errorf("IS NOT file: addons.ingress.valuesFile=%s", addons.ValuesFile)
			return nil
		}
		release.ValuesFile = addons.ValuesFile
	} else {
		valuesFile, err := util.CheckExitValueFile(args.ClusterName, release.Name)
		if err == nil {
			release.ValuesFile = valuesFile
		}
	}

	// Is Enabled monitoring
	// if args.ServiceMonitor {
	// 	overrides["prometheus.servicemonitor.enabled"] = "true"
	// }

	overrides := map[string]string{}

	if !update {
		overrides["installCRDs"] = "true"
	}

	if addons.HostMode {
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
	}

	overrides["defaultBackend.enabled"] = "true"

	if len(addons.DefaultBackend.Registry) > 0 {
		overrides["defaultBackend.image.registry"] = addons.DefaultBackend.Registry
	}
	if len(addons.DefaultBackend.Image) > 0 {
		overrides["defaultBackend.image.image"] = addons.DefaultBackend.Image
	}
	if len(addons.DefaultBackend.Tag) > 0 {
		overrides["defaultBackend.image.tag"] = addons.DefaultBackend.Tag
	}
	// overrides["defaultBackend.image.registry"] = "k8s.gcr.io"
	// overrides["defaultBackend.image.image"] = "defaultbackend-amd64"
	// overrides["defaultBackend.image.tag"] = "1.5"

	options := k3sv1alpha1.HelmOptions{
		ClusterName:     args.ClusterName,
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
