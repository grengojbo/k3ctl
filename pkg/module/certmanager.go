package module

import (
	// "github.com/alexellis/arkade/pkg/apps"
	// "github.com/alexellis/arkade/pkg/types"
	"fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	log "github.com/sirupsen/logrus"
)

func MakeInstallCertManager(certManager *k3sv1alpha1.CertManager, args *k3sv1alpha1.HelmRelease, kubeConfigPath string, dryRun bool) (err error) {
	name := "MakeInstallCertManager"
	description := "Cert Manager"
	update := false

	release, ok := k3sv1alpha1.FindRelease(args.Releases, certManager.Name)
	if !ok {
		return fmt.Errorf("[%s] is not release...", name)
	}

	// log.Debugf("[%s] name: %s disabled: %v status: %v", name, certManager.Name, certManager.Disabled, release.Status)

	if certManager.Disabled {
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

	overrides := map[string]string{}

	if !update {
		overrides["installCRDs"] = "true"
	}

	options := k3sv1alpha1.HelmOptions{
		CreateNamespace: false,
		KubeconfigPath:  kubeConfigPath,
		Overrides:       overrides,
		Helm:            &release,
		Wait:            args.Wait,
		Verbose:         false,
		DryRun:          dryRun,
	}
	err = Helm3Upgrade(&options)

	// if deleted {
	// 	log.Infoln("TODO: Deleted Cert Manager...")
	// 	// command := fmt.Sprintf(t.HelmDeleteCommand, ingress.Name, ingress.Namespace, kubeConfigPath)
	// 	// _, _, err := k3s.RunLocalCommand(command, false, dryRun)
	// 	// if err != nil {
	// 	// 	log.Errorf("[RunLocalCommand] %s\n%v", err.Error())
	// 	// }
	// } else if installed {
	// 	certmanagerOptions := types.DefaultInstallOptions().
	// 		WithNamespace(namespace).
	// 		WithHelmRepo("jetstack/cert-manager").
	// 		WithHelmURL("https://charts.jetstack.io").
	// 		WithOverrides(overrides).
	// 		WithWait(wait).
	// 		WithHelmUpdateRepo(updateRepo).
	// 		WithKubeconfigPath(kubeConfigPath)

	// 	_, err = apps.MakeInstallChart(certmanagerOptions)
	// 	if err != nil {
	// 		return err
	// 	}

	// 	// log.Infof(certManagerInstallMsg)
	// }
	return err
}

// const CertManagerInfoMsg = `# Get started with cert-manager here:
// # https://docs.cert-manager.io/en/latest/tutorials/acme/http-validation.html`

// const certManagerInstallMsg = `=======================================================================
// = cert-manager  has been installed.                                   =
// =======================================================================` +
// 	"\n\n" + CertManagerInfoMsg + "\n\n"
