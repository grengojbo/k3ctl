package module

import (
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

// VeleroSettings
func VeleroSettings(addons *k3sv1alpha1.Backup, clusterName string) (release k3sv1alpha1.HelmInterfaces) {
	repo := k3sv1alpha1.HelmRepo{
		Name: types.VeleroHelmRepoName,
		Repo: types.VeleroHelmRepo,
		Url:  types.VeleroHelmURL,
	}
	if len(addons.Velero.Repo.Name) > 0 {
		repo.Name = addons.Velero.Repo.Name
	}
	if len(addons.Velero.Repo.Repo) > 0 {
		repo.Repo = addons.Velero.Repo.Repo
	}
	if len(addons.Velero.Repo.Url) > 0 {
		repo.Url = addons.Velero.Repo.Url
	}

	if addons.Disabled {
		release.Deleted = true
	}
	if len(addons.Name) == 0 {
		addons.Name = types.BackupDefaultName
	}
	if len(addons.Namespace) == 0 {
		addons.Namespace = types.VeleroDefaultNamespace
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

	// Settings for Velero
	// TODO: добавить проверку есть ли креды для провайдера
	if len(addons.Velero.Providers) == 0 {
		addons.Velero.Providers = append(addons.Velero.Providers, "aws")
	}

	if len(addons.Schedules) == 0 {
		includedNamespaces := []string{"kube-system", "kube-public", "kube-node-lease"}
		schedule := k3sv1alpha1.SchedulesBackup{
			Name:               "systembackup",
			Schedule:           "0 0 * * *",
			IncludedNamespaces: includedNamespaces,
		}
		addons.Schedules = append(addons.Schedules, schedule)
	}
	// Set AWS Sorage
	storageAws := k3sv1alpha1.VeleroStorage{
		Name:   "aws",
		Image:  types.VeleroPluginAwsImage,
		Region: types.DefaultAwsRegion,
		Bucket: fmt.Sprintf("velero-%s", clusterName),
	}

	for _, provider := range addons.Velero.Providers {
		if provider == "aws" {
			addons.Velero.Storages = append(addons.Velero.Storages, storageAws)
		}
	}

	//  All Settings
	release.Name = addons.Name
	release.Namespace = addons.Namespace
	release.Repo = repo.Repo

	addons.Repo = repo
	return release
}

// MakeInstallVelero
// https://github.com/vmware-tanzu/helm-charts/blob/main/charts/velero/values.yaml
func MakeInstallVelero(addons *k3sv1alpha1.Backup, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallVelero"
	description := "Backup Velero"
	// update := false

	release, ok := k3sv1alpha1.FindRelease(args.Releases, addons.Name)
	if !ok {
		return fmt.Errorf("[%s] is not release...", name)
	}

	// log.Debugf("[%s] name: %s disabled: %v status: %v (cluster: %s)", name, addons.Name, addons.Disabled, release.Status, args.ClusterName)

	if len(addons.ValuesFile) > 0 {
		if err = util.CheckExitFile(addons.ValuesFile); err != nil {
			log.Errorf("IS NOT file: addons.backup.valuesFile=%s", addons.ValuesFile)
			return nil
		}
		release.ValuesFile = addons.ValuesFile
	} else {
		valuesFile, err := util.CheckExitValueFile(args.ClusterName, release.Name)
		if err == nil {
			release.ValuesFile = valuesFile
		}
	}
	// log.Warnf("ClusterName: %s release.Name: %s addons.ValuesFile: %s", args.ClusterName, release.Name, release.ValuesFile)
	// release.DependencyUpdate = true
	existingSecret := fmt.Sprintf("velero-%s-%s-creds", args.ClusterName, addons.Velero.Providers[0])
	secretFile := fmt.Sprintf("./variables/%s/secret-velero.ini", args.ClusterName)

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
		if err := util.CheckExitFile(secretFile); err != nil {
			log.Warnf("IS NOT Secret: %s", secretFile)
			log.Warnf("Is NOT install \"%s\" :(", description)
			return nil
		}

		log.Infof("Install %s...", description)

		secretsFile := k3sv1alpha1.SecretsData{
			Key:   "cloud",
			Value: secretFile,
			Type:  types.FromFileSecret,
		}
		secret := k3sv1alpha1.K8sSecret{
			Name:      existingSecret,
			Type:      "generic",
			Namespace: release.Namespace,
		}
		secret.SecretsData = append(secret.SecretsData, secretsFile)
		if err := CreateSecret(secret, kubeConfigPath, args.ClusterName, dryRun); err != nil {
			log.Errorf("Is NOT install \"%s\" (%s)", description, err.Error())
			return nil
		}
	}

	overrides := map[string]string{}

	// 	helm install velero vmware-tanzu/velero \
	// --namespace <YOUR NAMESPACE> \
	// --create-namespace \
	// --set-file credentials.secretContents.cloud=<FULL PATH TO FILE> \
	// --set configuration.provider=<PROVIDER NAME> \
	// --set configuration.backupStorageLocation.name=<BACKUP STORAGE LOCATION NAME> \
	// --set configuration.backupStorageLocation.bucket=<BUCKET NAME> \
	// --set configuration.backupStorageLocation.config.region=<REGION> \
	// --set configuration.volumeSnapshotLocation.name=<VOLUME SNAPSHOT LOCATION NAME> \
	// --set configuration.volumeSnapshotLocation.config.region=<REGION> \
	// --set initContainers[0].name=velero-plugin-for-<PROVIDER NAME> \
	// --set initContainers[0].image=velero/velero-plugin-for-<PROVIDER NAME>:<PROVIDER PLUGIN TAG> \
	// --set initContainers[0].volumeMounts[0].mountPath=/target \
	// --set initContainers[0].volumeMounts[0].name=plugins

	overrides["credentials.useSecret"] = "true"
	overrides["credentials.existingSecret"] = existingSecret
	overrides["configuration.provider"] = addons.Velero.Providers[0]

	// Is Enabled monitoring
	// if args.ServiceMonitor {
	// 	overrides["serviceMonitor.enabled"] = "true"
	// }

	for i, v := range addons.Velero.Storages {
		name := fmt.Sprintf("%s-%s", v.Name, args.ClusterName)
		overrides["configuration.backupStorageLocation.name"] = name
		if i == 0 {
			overrides["configuration.backupStorageLocation.default"] = "true"
		}
		overrides["configuration.backupStorageLocation.bucket"] = v.Bucket
		overrides["configuration.backupStorageLocation.config.region"] = v.Region
		overrides["configuration.volumeSnapshotLocation.name"] = name
		overrides["configuration.volumeSnapshotLocation.config.region"] = v.Region

		overrides[fmt.Sprintf("initContainers[%d].name", i)] = v.Name
		overrides[fmt.Sprintf("initContainers[%d].image", i)] = v.Image
		overrides[fmt.Sprintf("initContainers[%d].volumeMounts[%d].mountPath", i, i)] = "target"
		overrides[fmt.Sprintf("initContainers[%d].volumeMounts[%d].name", i, i)] = "plugins"
	}

	for _, v := range addons.Schedules {
		if v.Disabled {
			overrides[fmt.Sprintf("schedules.%s.disabled", v.Name)] = "true"
		} else {
			overrides[fmt.Sprintf("schedules.%s.disabled", v.Name)] = "false"
		}
		for key, value := range v.Labels {
			overrides[fmt.Sprintf("schedules.%s.labels.%s", v.Name, key)] = value
		}
		for key, value := range v.Annotations {
			overrides[fmt.Sprintf("schedules.%s.annotations.%s", v.Name, key)] = value
		}
		overrides[fmt.Sprintf("schedules.%s.schedule", v.Name)] = v.Schedule
		// overrides[fmt.Sprintf("schedules.%s.schedule", v.Name)] = fmt.Sprintf("\"%s\"", v.Schedule)
		if v.UseOwnerReferencesInBackup {
			overrides[fmt.Sprintf("schedules.%s.useOwnerReferencesInBackup", v.Name)] = "true"
		} else {
			overrides[fmt.Sprintf("schedules.%s.useOwnerReferencesInBackup", v.Name)] = "false"
		}
		if len(v.Ttl) > 0 {
			overrides[fmt.Sprintf("schedules.%s.template.ttl", v.Name)] = v.Ttl
		} else {
			overrides[fmt.Sprintf("schedules.%s.template.ttl", v.Name)] = "240h"
		}
		for i, val := range v.IncludedNamespaces {
			overrides[fmt.Sprintf("schedules.%s.template.includedNamespaces[%d]", v.Name, i)] = val
		}
		// overrides[fmt.Sprintf("schedules.%s.", v.Name)] =
	}

	// overrides["resources.requests.cpu"] = "500m"
	// overrides["resources.requests.memory"] = "512Mi"
	// overrides["resources.limits.cpu"] = "1000m"
	// overrides["resources.limits.memory"] = "1024Mi"

	// overrides["restic.resources.requests.cpu"] = "500m"
	// overrides["restic.resources.requests.memory"] = "512Mi"
	// overrides["restic.resources.limits.cpu"] = "1000m"
	// overrides["restic.resources.limits.memory"] = "1024Mi"
	// overrides[""] = ""

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
