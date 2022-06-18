package module

import (
	// "github.com/alexellis/arkade/pkg/apps"
	// "github.com/alexellis/arkade/pkg/types"
	// "fmt"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	// "github.com/grengojbo/k3ctl/pkg/util"
	// log "github.com/sirupsen/logrus"
)

// LoadBalancerSettings
func LoadBalancerSettings(lb *k3sv1alpha1.LoadBalancer, clusterName string) {
	if len(lb.Name) == 0 {
		lb.Name = types.ServiceLb
	}
	// log.Warnln("TODO: LoadBalancerSettings")
	// repo := k3sv1alpha1.HelmRepo{
	// 	Name: types.CertManagerHelmRepoName,
	// 	Repo: types.CertManagerHelmRepo,
	// 	Url:  types.CertManagerHelmURL,
	// }
	// if len(addons.Repo.Name) > 0 {
	// 	repo.Name = addons.Repo.Name
	// }
	// if len(addons.Repo.Repo) > 0 {
	// 	repo.Repo = addons.Repo.Repo
	// }
	// if len(addons.Repo.Url) > 0 {
	// 	repo.Url = addons.Repo.Url
	// }

	// if addons.Disabled {
	// 	release.Deleted = true
	// }
	// if len(addons.Name) == 0 {
	// 	addons.Name = types.CertManagerDefaultName
	// }
	// if len(addons.Namespace) == 0 {
	// 	addons.Namespace = types.CertManagerDefaultNamespace
	// }
	// if len(addons.Version) > 0 {
	// 	release.Version = addons.Version
	// }
	// if len(addons.Values) > 0 {
	// 	release.Values = addons.Values
	// }
	// if len(addons.ValuesFile) > 0 {
	// 	release.ValuesFile = addons.ValuesFile
	// }

	// //  All Settings
	// release.Name = addons.Name
	// release.Namespace = addons.Namespace
	// release.Repo = repo.Repo

	// addons.Repo = repo
	// return release
}
