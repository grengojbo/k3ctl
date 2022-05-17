package module

import (
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"
	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	log "github.com/sirupsen/logrus"
)

func MakeInstallCertManager(kubeConfigPath string, dryRun bool, certManager *k3sv1alpha1.CertManager, args *k3sv1alpha1.HelmRelease) (err error) {
	namespace := "cert-manager"
	wait := true
	updateRepo := true
	installed := false
	deleted := false
	// kubeConfigPath := ""

	overrides := map[string]string{}

	if args.UpdateStrategy == "none" {
		installed = false
	} else {
		overrides["installCRDs"] = "true"
		// ok, release := k3sv1alpha1.GetHelmRelease(ingress.Name, args.Releases)
		// if ok {
		// 	installed = false
		// 	if ingress.Disabled {
		// 		deleted = true
		// 	} else {
		// 		if len(ingress.Version) > 0 && ingress.Version != release.AppVersion {
		// 			installed = true
		// 		}
		// 	}
		// 	// log.Infof("Install Nginx Ingress controller %v", release.Revision)
		// }
	}

	if deleted {
		log.Infoln("TODO: Deleted Cert Manager...")
		// command := fmt.Sprintf(t.HelmDeleteCommand, ingress.Name, ingress.Namespace, kubeConfigPath)
		// _, _, err := k3s.RunLocalCommand(command, false, dryRun)
		// if err != nil {
		// 	log.Errorf("[RunLocalCommand] %s\n%v", err.Error())
		// }
	} else if installed {
		certmanagerOptions := types.DefaultInstallOptions().
			WithNamespace(namespace).
			WithHelmRepo("jetstack/cert-manager").
			WithHelmURL("https://charts.jetstack.io").
			WithOverrides(overrides).
			WithWait(wait).
			WithHelmUpdateRepo(updateRepo).
			WithKubeconfigPath(kubeConfigPath)

		_, err = apps.MakeInstallChart(certmanagerOptions)
		if err != nil {
			return err
		}

		log.Infof(certManagerInstallMsg)
	}
	return err
}

const CertManagerInfoMsg = `# Get started with cert-manager here:
# https://docs.cert-manager.io/en/latest/tutorials/acme/http-validation.html`

const certManagerInstallMsg = `=======================================================================
= cert-manager  has been installed.                                   =
=======================================================================` +
	"\n\n" + CertManagerInfoMsg + "\n\n"
