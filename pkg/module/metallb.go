package module

// https://metallb.universe.tf/installation/
// https://github.com/metallb/metallb/tree/main/charts/metallb
// https://habr.com/ru/company/southbridge/blog/443110/

// MakeInstallMetalLB
func MakeInstallMetalLB(kubeConfigPath string, dryRun bool) (err error) {
	// TODO: add support metallb
	// ~/go/src/github.com/alexellis/arkade/cmd/apps/metallb_app.go
	return nil
}

// MetalLBShort     :=   "Install MetalLB in L2 (ARP) mode"
// MetalLBLong     :=    `Install a network load-balancer implementation for Kubernetes using standard routing protocols`
// MetalLBExample :=     `arkade install metallb-arp --address-range=<cidr>`

// configInline:
//   address-pools:
//    - name: default
//      protocol: layer2
//      addresses:
//      - 198.51.100.0/24
