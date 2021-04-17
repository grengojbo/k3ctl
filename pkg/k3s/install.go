package k3s

import (
	"fmt"
	"strings"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	log "github.com/sirupsen/logrus"
)

// https://github.com/alexellis/k3sup/blob/master/cmd/install.go

var kubeconfig []byte

type K3sExecOptions struct {
	Datastore           string
	ExtraArgs           []string
	FlannelIPSec        bool
	NoExtras            bool
	LoadBalancer        *k3sv1alpha1.LoadBalancer
	Ingress             string
	DisableLoadbalancer bool
	DisableIngress      bool
}

type K3sIstallOptions struct {
	ExecString   string
	LoadBalancer string
}

func MakeInstallExec(cluster bool, host, tlsSAN string, options K3sExecOptions) K3sIstallOptions {
	extraArgs := []string{}
	k3sIstallOptions := K3sIstallOptions{}
	// if len(options.Datastore) > 0 {
	// 	extraArgs = append(extraArgs, fmt.Sprintf("--datastore-endpoint %s", options.Datastore))
	// }
	// if options.FlannelIPSec {
	// 	extraArgs = append(extraArgs, "--flannel-backend ipsec")
	// }

	if options.DisableLoadbalancer {
		extraArgs = append(extraArgs, "--no-deploy servicelb")
	} else {
		if len(options.LoadBalancer.MetalLb) > 0 {
			// TODO: добавить проверку на ip adress
			log.Debugln("LoadBalancer MetalLB: ", options.LoadBalancer.MetalLb)
			extraArgs = append(extraArgs, "--no-deploy servicelb")
			k3sIstallOptions.LoadBalancer = types.MetalLb
		} else if len(options.LoadBalancer.KubeVip) > 0 {
			// TODO: добавить проверку на ip adress
			log.Debugln("LoadBalancer kube-vip: ", options.LoadBalancer.KubeVip)
			extraArgs = append(extraArgs, "--no-deploy servicelb")
			k3sIstallOptions.LoadBalancer = types.KubeVip
		}
	}

	if options.DisableIngress || len(options.Ingress) != 0 {
		extraArgs = append(extraArgs, "--no-deploy traefik")
	}

	extraArgsCmdline := ""
	for _, a := range extraArgs {
		extraArgsCmdline += a + " "
	}

	for _, a := range options.ExtraArgs {
		if a != "[]" {
			extraArgsCmdline += a + " "
		}
	}

	installExec := "INSTALL_K3S_EXEC='server"
	if cluster {
		installExec += " --cluster-init"
	}

	// san := host
	// if len(tlsSAN) > 0 {
	// 	san = tlsSAN
	// }
	// installExec += fmt.Sprintf(" --tls-san %s", san)

	if trimmed := strings.TrimSpace(extraArgsCmdline); len(trimmed) > 0 {
		installExec += fmt.Sprintf(" %s", trimmed)
	}

	installExec += "'"

	k3sIstallOptions.ExecString = installExec

	if len(k3sIstallOptions.LoadBalancer) == 0 {
		k3sIstallOptions.LoadBalancer = types.ServiceLb
	}
	// --tls-san developer.iwis.io --cluster-cidr 10.42.0.0/19 --service-cidr 10.42.32.0/19 --cluster-dns 10.42.32.10 --flannel-backend=none --secrets-encryption --node-taint CriticalAddonsOnly=true:NoExecute
	return k3sIstallOptions
}
