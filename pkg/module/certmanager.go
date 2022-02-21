package module

import (
	"github.com/alexellis/arkade/pkg/apps"
	"github.com/alexellis/arkade/pkg/types"
	log "github.com/sirupsen/logrus"
)

func MakeInstallCertManager(kubeConfigPath string) (err error) {
	namespace := "cert-manager"
	wait := true
	updateRepo := true
	// kubeConfigPath := ""

	overrides := map[string]string{}
	overrides["installCRDs"] = "true"

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
	return err
}

const CertManagerInfoMsg = `# Get started with cert-manager here:
# https://docs.cert-manager.io/en/latest/tutorials/acme/http-validation.html`

const certManagerInstallMsg = `=======================================================================
= cert-manager  has been installed.                                   =
=======================================================================` +
	"\n\n" + CertManagerInfoMsg + "\n\n"