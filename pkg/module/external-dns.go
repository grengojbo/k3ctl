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
// func ExternalDnsSettings(addons *k3sv1alpha1.ExternalDns, lb *k3sv1alpha1.LoadBalancer, clusterName string) (release k3sv1alpha1.HelmInterfaces) {
func ExternalDnsSettings(spec *k3sv1alpha1.ClusterSpec) (release k3sv1alpha1.HelmInterfaces) {
	addons := &spec.Addons.ExternalDns
	// lb := &spec.LoadBalancer

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
	}

	if len(addons.Provider) == 0 {
		if len(spec.Providers.Default) > 0 {
			addons.Provider = spec.Providers.Default
		}

	}
	// DNS provider where the DNS records will be created.
	// provider

	// aws.credentials.secretKey	When using the AWS provider, set aws_secret_access_key in the AWS credentials (optional)	""
	// aws.credentials.accessKey	When using the AWS provider, set aws_access_key_id in the AWS credentials (optional)
	// Use an existing secret with key "credentials" defined.
	// aws.credentials.secretName
	// When using the AWS provider, AWS_DEFAULT_REGION to set in the environment (optional)
	// aws.region

	// When using the Azure provider, set the secret containing the azure.json file
	// azure.secretName

	// cloudflare.apiToken	When using the Cloudflare provider, CF_API_TOKEN to set (optional)	""
	// cloudflare.apiKey	When using the Cloudflare provider, CF_API_KEY to set (optional)	""
	// cloudflare.secretName	When using the Cloudflare provider, it's the name of the secret containing cloudflare_api_token or cloudflare_api_key.	""
	// cloudflare.email	When using the Cloudflare provider, CF_API_EMAIL to set (optional). Needed when using CF_API_KEY	""
	// cloudflare.proxied	When using the Cloudflare provider, enable the proxy feature (DDOS protection, CDN...) (optional)	true

	// domainFilters	Limit possible target zones by domain suffixes (optional)	[]
	// excludeDomains	Exclude subdomains (optional)	[]
	// regexDomainFilter	Limit possible target zones by regex domain suffixes (optional)	""
	// regexDomainExclusion	Exclude subdomains by using regex pattern (optional)	""
	// zoneNameFilters	Filter target zones by zone domain (optional)	[]
	// zoneIdFilters	Limit possible target zones by zone id (optional)	[]
	// annotationFilter	Filter sources managed by external-dns via annotation using label selector (optional)	""
	// labelFilter	Select sources managed by external-dns using label selector (optional)	""
	// crd.create	Install and use the integrated DNSEndpoint CRD	false

	// txtOwnerId	A name that identifies this instance of ExternalDNS. Currently used by registry types: txt & aws-sd (optional)

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
func MakeInstallExternalDns(spec *k3sv1alpha1.ClusterSpec, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	addons := &spec.Addons.ExternalDns
	name := "MakeInstallExternalDns"
	description := "External DNS"
	// update := false

	overrides := map[string]string{}

	release, ok := k3sv1alpha1.FindRelease(args.Releases, addons.Name)
	if !ok {
		return fmt.Errorf("[%s] is not release...", name)
	}

	// log.Debugf("[%s] name: %s disabled: %v status: %v (cluster: %s)", name, addons.Name, addons.Disabled, release.Status, args.ClusterName)

	if len(spec.LoadBalancer.Domain) == 0 {
		log.Warnf("IS NOT Set \"spec.loadBalancer.domain\" %s disabled...", description)
		return nil
	}
	if len(addons.Provider) == 0 {
		log.Warnf("IS NOT Set \"spec.externalDns.provider\" %s disabled...", description)
		return nil
	}

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

	// if addons.Provider == types.ProviderAws {
	// 	if ok, secretFile := util.CheckСredentials(spec.ClusterName, types.ProviderAws); ok {

	// 	}
	// }

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

	if len(spec.LoadBalancer.ExternalIP) == 0 {
		log.Warnln("IS NOT Set \"spec.loadBalancer.externalIP\", is used private ip disabled...")
	}

	// if !update {
	// 	overrides["installCRDs"] = "true"
	// }

	//  -- List of IP addresses at which the controller services are available
	//  Ref: https://kubernetes.io/docs/user-guide/services/#external-ips
	// if len(lb.ExternalIP) > 0 {
	// 	overrides["controller.service.externalIPs[0]"] = lb.ExternalIP
	// }

	// Is Enabled monitoring
	// metrics.enabled	Enable prometheus to access external-dns metrics endpoint	false
	// metrics.serviceMonitor.enabled	Create ServiceMonitor object	false
	// if args.ServiceMonitor {
	// 	overrides["prometheus.servicemonitor.enabled"] = "true"
	// }

	// helm install my-release \
	// --set provider=aws \
	// --set aws.zoneType=public \
	// --set txtOwnerId=HOSTED_ZONE_IDENTIFIER \
	// --set domainFilters[0]=HOSTED_ZONE_NAME \
	// bitnami/external-dns

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
