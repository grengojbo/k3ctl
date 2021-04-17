package k3s

import (
	"fmt"
	"strings"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
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
	Networking          *k3sv1alpha1.Networking
	Ingress             string
	DisableLoadbalancer bool
	DisableIngress      bool
	SELinux             bool
	Rootless            bool
	SecretsEncryption   bool
}

type K3sIstallOptions struct {
	ExecString   string
	LoadBalancer string
	Ingress      string
	CNI          string
	Backend      string
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

	if options.DisableIngress || len(options.Ingress) > 0 {
		if ingress, isset := util.Find(types.IngressControllers, options.Ingress); isset {
			k3sIstallOptions.Ingress = ingress
			extraArgs = append(extraArgs, "--no-deploy traefik")
		} else {
			log.Fatalf("Ingress Controllers %s not support :(", options.Ingress)
		}
	}

	if len(options.Networking.ServiceSubnet) > 0 {
		log.Debugln("ServiceSubnet: ", options.Networking.ServiceSubnet)
		extraArgs = append(extraArgs, fmt.Sprintf("--service-cidr %s", options.Networking.ServiceSubnet))
	}

	if len(options.Networking.PodSubnet) > 0 {
		log.Debugln("PodSubnet: ", options.Networking.PodSubnet)
		extraArgs = append(extraArgs, fmt.Sprintf("--cluster-cidr %s", options.Networking.PodSubnet))
	}

	if len(options.Networking.DNSDomain) > 0 {
		log.Debugln("DNSDomain: ", options.Networking.DNSDomain)
		extraArgs = append(extraArgs, fmt.Sprintf("--cluster-domain %s", options.Networking.DNSDomain))
	}

	if len(options.Networking.ClusterDns) > 0 {
		log.Debugln("ClusterDns: ", options.Networking.ClusterDns)
		extraArgs = append(extraArgs, fmt.Sprintf("--cluster-dns %s", options.Networking.ClusterDns))
	}

	k3sIstallOptions.Backend = types.Vxlan
	k3sIstallOptions.CNI = types.Flannel
	if len(options.Networking.CNI) > 0 {
		if cni, isset := util.Find(types.CNIplugins, options.Networking.CNI); isset {
			k3sIstallOptions.CNI = cni
		} else {
			log.Fatalf("CNI plugins %s not support :(", options.Networking.CNI)
		}
	}
	if len(options.Networking.Backend) > 0 {
		if k3sIstallOptions.CNI == types.Flannel {
			if backend, isset := util.Find(types.FlannelBackends, options.Networking.Backend); isset {
				k3sIstallOptions.Backend = backend
			} else {
				log.Fatalf("CNI plugins %s backend %s not support :(", options.Networking.CNI, options.Networking.Backend)
			}
		} else if k3sIstallOptions.CNI == types.Calico {
			if backend, isset := util.Find(types.CalicoBackends, options.Networking.Backend); isset {
				k3sIstallOptions.Backend = backend
			} else {
				log.Fatalf("CNI plugins %s backend %s not support :(", options.Networking.CNI, options.Networking.Backend)
			}
		} else if k3sIstallOptions.CNI == types.Cilium {
			if backend, isset := util.Find(types.CiliumBackends, options.Networking.Backend); isset {
				k3sIstallOptions.Backend = backend
			} else {
				log.Fatalf("CNI plugins %s backend %s not support :(", options.Networking.CNI, options.Networking.Backend)
			}
		}
	}
	if k3sIstallOptions.CNI == types.Flannel {
		extraArgs = append(extraArgs, fmt.Sprintf("--flannel-backend=%s", k3sIstallOptions.Backend))
	} else {
		extraArgs = append(extraArgs, "--flannel-backend=none")
	}

	if options.SecretsEncryption {
		extraArgs = append(extraArgs, "--secrets-encryption")
	}

	if options.SELinux {
		extraArgs = append(extraArgs, "--selinux")
	}

	if options.Rootless {
		extraArgs = append(extraArgs, "--rootless")
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
	// --tls-san developer.iwis.io --secrets-encryption --node-taint CriticalAddonsOnly=true:NoExecute
	return k3sIstallOptions
}
