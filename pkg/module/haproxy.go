package module

import (
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

// HaproxySettings
func HaproxySettings(addons *k3sv1alpha1.Ingress, lb *k3sv1alpha1.LoadBalancer, clusterName string) (release k3sv1alpha1.HelmInterfaces) {
	repo := k3sv1alpha1.HelmRepo{
		Name: types.HaproxyHelmRepoName,
		Repo: types.HaproxyHelmRepo,
		Url:  types.HaproxyHelmURL,
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
		addons.Name = types.HaproxyDefaultName
	}
	if len(addons.Namespace) == 0 {
		addons.Namespace = types.HaproxyDefaultNamespace
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

	// Install CRD manifests
	// if len(addons.Manifests) == 0 {
	// 	addons.Manifests = append(addons.Manifests, types.HaproxyCrdBackend)
	// 	addons.Manifests = append(addons.Manifests, types.HaproxyCrdDefaults)
	// 	addons.Manifests = append(addons.Manifests, types.HaproxyCrdGlobal)
	// }

	//  All Settings
	release.Name = addons.Name
	release.Namespace = addons.Namespace
	release.Repo = repo.Repo
	release.Manifests = addons.Manifests

	addons.Repo = repo
	addons.LoadBalancer = lb
	return release
}

// MakeInstallHaproxy
func MakeInstallHaproxy(addons *k3sv1alpha1.Ingress, args *k3sv1alpha1.HelmRelease, monitoring *k3sv1alpha1.Monitoring, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallHaproxy"
	description := "Ingress Haproxy"
	// update := false

	release, ok := k3sv1alpha1.FindRelease(args.Releases, types.HaproxyDefaultName)
	if !ok {
		return fmt.Errorf("[%s] is not release...", name)
	}

	// log.Debugf("[%s] name: %s disabled: %v status: %v", name, addons.Name, addons.Disabled, release.Status)

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

	// https://github.com/haproxytech/helm-charts/blob/main/kubernetes-ingress/values.yaml
	overrides := map[string]string{}

	// if !update {
	// 	overrides["installCRDs"] = "true"
	// }

	//  -- List of IP addresses at which the controller services are available
	//  Ref: https://kubernetes.io/docs/user-guide/services/#external-ips
	if len(addons.LoadBalancer.ExternalIP) > 0 {
		overrides["controller.service.externalIPs[0]"] = addons.LoadBalancer.ExternalIP
	}

	// Is Enabled monitoring
	// if args.ServiceMonitor {
	// ServiceMonitor
	// ref: https://github.com/prometheus-operator/prometheus-operator/blob/master/Documentation/user-guides/getting-started.md
	// Note: requires Prometheus Operator to be able to work, for example:
	// helm install prometheus prometheus-community/kube-prometheus-stack \
	//   --set prometheus.prometheusSpec.podMonitorSelectorNilUsesHelmValues=false \
	//   --set prometheus.prometheusSpec.serviceMonitorSelectorNilUsesHelmValues=false
	// overrides["controller.serviceMonitor.enabled"] = "true"
	// // overrides["controller.metrics."] = ""
	// }

	if addons.HostMode {
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
	if len(addons.DefaultBackend.Registry) > 0 {
		overrides["defaultBackend.image.registry"] = addons.DefaultBackend.Registry
	}
	if len(addons.DefaultBackend.Image) > 0 {
		overrides["defaultBackend.image.image"] = addons.DefaultBackend.Image
	}
	if len(addons.DefaultBackend.Tag) > 0 {
		overrides["defaultBackend.image.tag"] = addons.DefaultBackend.Tag
	}

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
