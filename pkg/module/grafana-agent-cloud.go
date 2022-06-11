package module

import (
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

// GrafanaAgentCloudSettings
func GrafanaAgentCloudSettings(addons *k3sv1alpha1.Monitoring, clusterName string) (release k3sv1alpha1.HelmInterfaces) {
	repo := k3sv1alpha1.HelmRepo{
		Name: types.GrafanaAgentCloudHelmRepoName,
		Repo: types.GrafanaAgentCloudHelmRepo,
		Url:  types.GrafanaAgentCloudHelmURL,
	}

	if addons.Disabled {
		release.Deleted = true
	}
	if len(addons.Name) == 0 {
		addons.Name = types.GrafanaAgentCloudDefaultName
	}
	if len(addons.Namespace) == 0 {
		addons.Namespace = types.MonitoringDefaultNamespace
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
			log.Warnf("IS NOT SET addons.monitoring.valuesFile OR %s", err.Error())
			log.Warnf("TODO: add link to documentation...")
			return nil
		}
	} else {
		if err = util.CheckExitFile(addons.ValuesFile); err != nil {
			log.Warnf("IS NOT file: addons.monitoring.valuesFile=%s", addons.ValuesFile)
			log.Warnf("TODO: add link to documentation...")
			return nil
		}
	}
	// log.Warnf("ValuesFile: %s", addons.ValuesFile)
	release.ValuesFile = addons.ValuesFile
	// release.DependencyUpdate = true

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

	overrides["clusterName"] = args.ClusterName

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
