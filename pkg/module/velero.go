package module

import (
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

// MakeInstallVelero
func MakeInstallVelero(addons *k3sv1alpha1.Backup, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallVelero"
	description := "Backup Velero"
	// update := false

	release, ok := k3sv1alpha1.FindRelease(args.Releases, addons.Name)
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
	// release.DependencyUpdate = true
	existingSecret := fmt.Sprintf("velero-%s-%s-creds", args.ClusterName, addons.Provider)

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
		secret := k3sv1alpha1.K8sSecret{
			Name: existingSecret,
		}
		if err := CreateSecret(secret); err != nil {
			log.Errorf("Is NOT install \"%s\":\n%s", description, err.Error())
			return nil
		}
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

	// overrides["clusterName"] = args.ClusterName

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
	overrides["configuration.provider"] = addons.Provider

	overrides["configuration.backupStorageLocation.name"] = fmt.Sprintf("%s-%s-backup", addons.Provider, args.ClusterName)
	overrides["configuration.backupStorageLocation.bucket"] = addons.Bucket
	overrides["configuration.backupStorageLocation.config.region"] = addons.Region
	overrides["configuration.volumeSnapshotLocation.name"] = fmt.Sprintf("%s-%s-snapshot", addons.Provider, args.ClusterName)
	overrides["configuration.volumeSnapshotLocation.config.region"] = addons.Region

	if addons.Provider == "aws" {
		overrides["initContainers[0].name"] = "aws"
		overrides["initContainers[0].image"] = "velero/velero-plugin-for-aws:v1.4.1"
		overrides["initContainers[0].volumeMounts[0].mountPath"] = "target"
		overrides["initContainers[0].volumeMounts[0].name"] = "plugins"
	}
	// overrides[""] = ""
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

	// ./variables/iwisops/secret-velero.ini
	// kubectl create secret generic velero-iwisops-aws-creds --from-file=cloud=./variables/iwisops/secret-velero.ini
	// options := k3sv1alpha1.HelmOptions{
	// 	CreateNamespace: false,
	// 	KubeconfigPath:  kubeConfigPath,
	// 	Overrides:       overrides,
	// 	Helm:            &release,
	// 	Wait:            args.Wait,
	// 	Verbose:         false,
	// 	DryRun:          dryRun,
	// }
	// err = Helm3Upgrade(&options)

	return err
}
