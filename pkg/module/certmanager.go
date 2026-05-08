package module

import (
	// "github.com/alexellis/arkade/pkg/apps"
	// "github.com/alexellis/arkade/pkg/types"
	"fmt"
	"os"
	"os/exec"
	"strings"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/k3s"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

// CertManagerSettings
func CertManagerSettings(addons *k3sv1alpha1.CertManager, clusterName string) (release k3sv1alpha1.HelmInterfaces) {
	repo := k3sv1alpha1.HelmRepo{
		Name: types.CertManagerHelmRepoName,
		Repo: types.CertManagerHelmRepo,
		Url:  types.CertManagerHelmURL,
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
		addons.Name = types.CertManagerDefaultName
	}
	if len(addons.Namespace) == 0 {
		addons.Namespace = types.CertManagerDefaultNamespace
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

// MakeInstallCertManager
func MakeInstallCertManager(addons *k3sv1alpha1.CertManager, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallCertManager"
	description := "Cert Manager"
	update := false

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
		update = true
	} else {
		log.Infof("Install %s...", description)
	}

	if len(addons.ValuesFile) > 0 {
		if err = util.CheckExitFile(addons.ValuesFile); err != nil {
			log.Errorf("IS NOT file: addons.certManager.valuesFile=%s", addons.ValuesFile)
			return nil
		}
		release.ValuesFile = addons.ValuesFile
	} else {
		valuesFile, err := util.CheckExitValueFile(args.ClusterName, release.Name)
		if err == nil {
			release.ValuesFile = valuesFile
		}
	}

	// Load .env: cluster-level first (variables/<cluster>/.env), then fall back to .env
	clusterEnvFile := fmt.Sprintf("variables/%s/.env", args.ClusterName)
	if err := util.LoadDotEnv(clusterEnvFile); err != nil {
		log.Warnf("[%s] LoadDotEnv(%s): %v", name, clusterEnvFile, err)
	}
	if err := util.LoadDotEnv(".env"); err != nil {
		log.Warnf("[%s] LoadDotEnv(.env): %v", name, err)
	}
	if token := os.Getenv("CF_API_TOKEN"); len(token) > 0 {
		command := fmt.Sprintf(
			"kubectl create ns cert-manager --dry-run=client -o yaml | kubectl apply -f - --kubeconfig %s --context %s",
			kubeConfigPath, args.ClusterName,
		)
		if _, _, cerr := k3s.RunLocalCommand(command, false, dryRun); cerr != nil {
			log.Warnf("[%s] ensure namespace cert-manager: %v", name, cerr)
		}
		command = fmt.Sprintf(
			"kubectl -n cert-manager create secret generic cloudflare-api-token --from-literal=api-token=%s --dry-run=client -o yaml | kubectl apply -f - --kubeconfig %s --context %s",
			token, kubeConfigPath, args.ClusterName,
		)
		log.Infof("[%s] creating cloudflare-api-token secret from CF_API_TOKEN", name)
		if _, errOut, cerr := k3s.RunLocalCommand(command, false, dryRun); cerr != nil {
			log.Errorf("[%s] create cloudflare-api-token secret: %v %s", name, cerr, errOut)
		}
	} else {
		log.Debugf("[%s] CF_API_TOKEN not set — skipping cloudflare secret creation", name)
	}

	overrides := map[string]string{}

	if !update {
		overrides["installCRDs"] = "true"
	}

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

	if err == nil {
		if len(addons.Manifests) > 0 {
			ApplyManifests(addons.Manifests, kubeConfigPath, args.ClusterName, dryRun)
		} else if len(addons.Provider) > 0 && addons.Provider != k3sv1alpha1.CertManagerProviderHTTP {
			applyClusterIssuer(addons, kubeConfigPath, args.ClusterName, dryRun)
		}
	}

	return err
}

// applyClusterIssuer generates and applies a ClusterIssuer manifest based on provider.
func applyClusterIssuer(addons *k3sv1alpha1.CertManager, kubeConfigPath string, clusterName string, dryRun bool) {
	name := "applyClusterIssuer"
	email := addons.Email
	if len(email) == 0 {
		email = os.Getenv("CERT_MANAGER_EMAIL")
	}
	if len(email) == 0 {
		log.Warnf("[%s] email is not set — set certManager.email or CERT_MANAGER_EMAIL", name)
		return
	}

	if addons.Provider == k3sv1alpha1.CertManagerProviderCloudflare {
		if token := os.Getenv("CF_API_TOKEN"); len(token) == 0 {
			log.Errorf("[%s] provider=cloudflare but CF_API_TOKEN is not set — ClusterIssuer will NOT be applied. Set it in variables/%s/.env or export CF_API_TOKEN=...", name, clusterName)
			return
		}
	}

	var yaml string
	switch addons.Provider {
	case k3sv1alpha1.CertManagerProviderCloudflare:
		yaml = fmt.Sprintf(`---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: %s
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
      - dns01:
          cloudflare:
            apiTokenSecretRef:
              name: cloudflare-api-token
              key: api-token
`, email)
	case k3sv1alpha1.CertManagerProviderRoute53:
		yaml = fmt.Sprintf(`---
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: %s
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
      - dns01:
          route53:
            region: %s
`, email, os.Getenv("AWS_REGION"))
	default:
		log.Warnf("[%s] unknown provider: %s", name, addons.Provider)
		return
	}

	log.Infof("[%s] applying ClusterIssuer (provider=%s, email=%s)", name, addons.Provider, email)
	if dryRun {
		log.Infof("[%s] dry-run: kubectl apply -f -\n%s", name, yaml)
		return
	}

	args := []string{"apply", "-f", "-",
		"--kubeconfig", kubeConfigPath,
		"--context", clusterName,
	}
	cmd := exec.Command("kubectl", args...)
	cmd.Stdin = strings.NewReader(yaml)
	out, cerr := cmd.CombinedOutput()
	if cerr != nil {
		log.Errorf("[%s] kubectl apply: %v\n%s", name, cerr, string(out))
	} else {
		log.Infof("[%s] %s", name, strings.TrimSpace(string(out)))
	}
}

// const CertManagerInfoMsg = `# Get started with cert-manager here:
// # https://docs.cert-manager.io/en/latest/tutorials/acme/http-validation.html`

// const certManagerInstallMsg = `=======================================================================
// = cert-manager  has been installed.                                   =
// =======================================================================` +
// 	"\n\n" + CertManagerInfoMsg + "\n\n"
