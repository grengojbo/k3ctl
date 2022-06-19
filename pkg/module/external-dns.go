package module

import (
	// "github.com/alexellis/arkade/pkg/apps"
	// "github.com/alexellis/arkade/pkg/types"
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

// ExternalDnsSettings
func ExternalDnsSettings(addons *k3sv1alpha1.ExternalDns, lb *k3sv1alpha1.LoadBalancer, clusterName string) (release k3sv1alpha1.HelmInterfaces) {
	repo := k3sv1alpha1.HelmRepo{
		Name: types.ExternalDnsHelmRepoName,
		Repo: types.ExternalDnsHelmRepo,
		Url:  types.ExternalDnsHelmURL,
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
	} else {
		if len(lb.Domain) == 0 {
			log.Warnf("IS NOT Set loadBalancer.domain")
			if len(addons.Domains) == 0 {
				log.Warnf("IS NOT Set addons.externalDns.domains")
			}
		}
	}

	if len(addons.Name) == 0 {
		addons.Name = types.ExternalDnsDefaultName
	}
	if len(addons.Namespace) == 0 {
		addons.Namespace = types.ExternalDnsDefaultNamespace
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

// MakeInstallExternalDns
func MakeInstallExternalDns(addons *k3sv1alpha1.ExternalDns, lb *k3sv1alpha1.LoadBalancer, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallExternalDns"
	description := "External DNS"
	// update := false

	release, ok := k3sv1alpha1.FindRelease(args.Releases, addons.Name)
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
		// update = true
	} else {
		log.Infof("Install %s...", description)
	}

	if len(addons.ValuesFile) > 0 {
		if err = util.CheckExitFile(addons.ValuesFile); err != nil {
			log.Errorf("IS NOT file: addons.externalDns.valuesFile=%s", addons.ValuesFile)
			return nil
		}
		release.ValuesFile = addons.ValuesFile
	} else {
		valuesFile, err := util.CheckExitValueFile(args.ClusterName, release.Name)
		if err == nil {
			release.ValuesFile = valuesFile
		}
	}

	overrides := map[string]string{}

	// if !update {
	// 	overrides["installCRDs"] = "true"
	// }

	//  -- List of IP addresses at which the controller services are available
	//  Ref: https://kubernetes.io/docs/user-guide/services/#external-ips
	// if len(lb.ExternalIP) > 0 {
	// 	overrides["controller.service.externalIPs[0]"] = lb.ExternalIP
	// }

	// Is Enabled monitoring
	// if args.ServiceMonitor {
	// 	overrides["prometheus.servicemonitor.enabled"] = "true"
	// }

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

// const CertManagerInfoMsg = `# Get started with cert-manager here:
// # https://docs.cert-manager.io/en/latest/tutorials/acme/http-validation.html`

// const certManagerInstallMsg = `=======================================================================
// = cert-manager  has been installed.                                   =
// =======================================================================` +
// 	"\n\n" + CertManagerInfoMsg + "\n\n"
