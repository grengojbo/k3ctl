package k3s

// https://github.com/alexellis/k3sup/blob/master/cmd/install.go

var kubeconfig []byte

type K3sExecOptions struct {
	Datastore    string
	ExtraArgs    string
	FlannelIPSec bool
	NoExtras     bool
}

func MakeInstallExec(cluster bool, host, tlsSAN string, options K3sExecOptions) string {
	// extraArgs := []string{}
	// if len(options.Datastore) > 0 {
	// 	extraArgs = append(extraArgs, fmt.Sprintf("--datastore-endpoint %s", options.Datastore))
	// }
	// if options.FlannelIPSec {
	// 	extraArgs = append(extraArgs, "--flannel-backend ipsec")
	// }

	// if options.NoExtras {
	// 	extraArgs = append(extraArgs, "--no-deploy servicelb")
	// 	extraArgs = append(extraArgs, "--no-deploy traefik")
	// }

	// extraArgs = append(extraArgs, options.ExtraArgs)
	// extraArgsCmdline := ""
	// for _, a := range extraArgs {
	// 	extraArgsCmdline += a + " "
	// }

	installExec := "INSTALL_K3S_EXEC='server"
	if cluster {
		installExec += " --cluster-init"
	}

	// san := host
	// if len(tlsSAN) > 0 {
	// 	san = tlsSAN
	// }
	// installExec += fmt.Sprintf(" --tls-san %s", san)

	// if trimmed := strings.TrimSpace(extraArgsCmdline); len(trimmed) > 0 {
	// 	installExec += fmt.Sprintf(" %s", trimmed)
	// }

	installExec += "'"

	return installExec
}
