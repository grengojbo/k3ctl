package k3s

import (
	"errors"
	"fmt"
	"strings"

	execute "github.com/alexellis/go-execute/pkg/v1"
	operator "github.com/alexellis/k3sup/pkg/operator"

	oper "github.com/grengojbo/k3ctl/pkg/operator"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
	// "k8s.io/apimachinery/pkg/util/errors"
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
	K3sVersion   string
	K3sChannel   string
	Node         *k3sv1alpha1.Node
}

func MakeInstallExec(cluster bool, tlsSAN []string, options K3sExecOptions) K3sIstallOptions {
	extraArgs := []string{}
	k3sIstallOptions := K3sIstallOptions{}

	if len(options.Datastore) > 0 {
		extraArgs = append(extraArgs, fmt.Sprintf("--datastore-endpoint %s", options.Datastore))
	}

	if options.DisableLoadbalancer {
		extraArgs = append(extraArgs, "--no-deploy servicelb")
	} else {
		if len(options.LoadBalancer.MetalLb) > 0 {
			// TODO: #3 добавить проверку на ip adress
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

	if len(tlsSAN) > 0 {
		for _, san := range tlsSAN {
			installExec += fmt.Sprintf(" --tls-san %s", san)
		}
	}

	if trimmed := strings.TrimSpace(extraArgsCmdline); len(trimmed) > 0 {
		installExec += fmt.Sprintf(" %s", trimmed)
	}

	installExec += "'"

	k3sIstallOptions.ExecString = installExec

	if len(k3sIstallOptions.LoadBalancer) == 0 {
		k3sIstallOptions.LoadBalancer = types.ServiceLb
	}
	// --tls-san developer.cluster --node-taint CriticalAddonsOnly=true:NoExecute
	return k3sIstallOptions
}

// RunK3sCommand Выполняем команды по SSH или локально
// TODO: translate
func RunK3sCommand(bastion *k3sv1alpha1.BastionNode, installk3sExec *K3sIstallOptions, dryRun bool) error {
	installStr := util.CreateVersionStr(installk3sExec.K3sVersion, installk3sExec.K3sChannel)
	installK3scommand := fmt.Sprintf("%s | %s %s sh -\n", types.K3sGetScript, installk3sExec.ExecString, installStr)

	if len(installk3sExec.K3sChannel) == 0 && len(installk3sExec.K3sVersion) == 0 {
		return errors.New("Set kubernetesVersion or channel (Release channel: stable, latest, or i.e. v1.19)")
	}

	log.Infof("KubernetesVersion: %s (K3sChannel: %s)", installk3sExec.K3sVersion, installk3sExec.K3sChannel)
	log.Infof("CNI: %s Backend: %s", installk3sExec.CNI, installk3sExec.Backend)
	log.Infof("LoadBalancer: %s", installk3sExec.LoadBalancer)
	if len(installk3sExec.Ingress) > 0 {
		log.Infof("Ingress Controllers: %s", installk3sExec.Ingress)
	} else {
		log.Infoln("Ingress Controllers: default k3s Traefik")
	}
	log.Warnf("Bastion %s host: %s (ssh port: %d key: %s)", bastion.Name, bastion.Address, bastion.SshPort, bastion.SSHAuthorizedKey)
	log.Debugln("--------------------------------------------")

	// sudoPrefix := ""
	// if useSudo {
	// 	sudoPrefix = "sudo "
	// }
	// getConfigcommand := fmt.Sprintf(sudoPrefix + "cat /etc/rancher/k3s/k3s.yaml\n")

	if bastion.Name == "local" {
		log.Infoln("Run command in localhost........")
		// stdOut, stdErr, err := RunLocalCommand(installK3scommand, true, dryRun)
		_, stdErr, err := RunLocalCommand(installK3scommand, true, dryRun)
		if err != nil {
			log.Fatalln(err.Error())
		} else if len(stdErr) > 0 {
			log.Errorf("stderr: %q", stdErr)
		}
		// log.Infof("stdout: %q", stdOut)
	} else {
		if err := RunSshCommand(installK3scommand, bastion, true, dryRun); err != nil {
			log.Fatalln(err.Error())
		}
	}

	// RunExampleCommand()
	// RunExampleCommand2()
	return nil
}

// RunLocalCommand выполнение комманд на локальном хосте TODO: tranclate
func RunLocalCommand(myCommand string, saveKubeconfig bool, dryRun bool) (stdOut []byte, stdErr []byte, err error) {

	sudoPrefix := ""
	uid := execute.ExecTask{
		Command:     "echo",
		Args:        []string{"${UID}"},
		Shell:       true,
		StreamStdio: false, // если true то выводит вконсоль и в Stdout
	}
	resUid, err := uid.Execute()
	if err != nil {
		return stdOut, stdErr, err
	}

	if val, err := oper.ParseInt64Output(resUid.Stdout); err != nil {
		log.Fatalln(err.Error())
	} else if val != 0 {
		sudoPrefix = "sudo "
		// log.Debugf("Result: %v sudoPrefix: %s", val, sudoPrefix)
	}
	command := fmt.Sprintf("%s%s", sudoPrefix, myCommand)
	operator := operator.ExecOperator{}
	log.Infof("Executing: %s\n", command)
	log.Warningln("TODO: сейчас заглушка RunLocalCommand")
	// res, err := operator.Execute(command)
	res, err := operator.Execute("pwd")
	if err != nil {
		return stdOut, stdErr, err
	}

	if len(res.StdErr) > 0 {
		stdErr = res.StdErr
	}
	if len(res.StdOut) > 0 {
		stdOut = res.StdOut
	}

	if saveKubeconfig {
		log.Warningln("TODO: доделать сохранение Kubeconfig (obtainKubeconfig)")
		// if err = obtainKubeconfig(operator, getConfigcommand, host, context, localKubeconfig, merge, printConfig); err != nil {
		// 	return err
		// }
	}

	return stdOut, stdErr, nil
}

// RunSshCommand выполнение комманд на удаленном хосте по ssh TODO: tranclate
func RunSshCommand(myCommand string, bastion *k3sv1alpha1.BastionNode, saveKubeconfig bool, dryRun bool) error {
	// ssh := oper.NewSshConnection(bastion)
	ssh := oper.SSHOperator{}
	ssh.NewSSHOperator(bastion)
	// Проверяем запускаем комманду через бастион или напрямую
	if _, isset := util.Find(k3sv1alpha1.ConnectionHosts, bastion.Name); !isset {
		log.Errorln("TODO: Run command in from bastion host ........")
	} else {
		// log.Errorln("TODO: Run command in host ........")
		// log.Infof("Executing: %s\n", myCommand)
		// myCommand = "hostname -a"
		if dryRun {
			log.Infof("Executing: %s\n", myCommand)
		} else {
			log.Debugf("Executing: %s\n", myCommand)
			// Выполняем комманду по SSH
			if _, err := ssh.Run(myCommand); err != nil {
				log.Fatalln(err.Error())
			}
			// Demo stream
			// ssh.Stream("for i in {1..5}; do echo ${i}; sleep ; done; exit 2;", true)
			// ssh.Stream("apt update -y", false)
		}
	}
	// if dryRun {
	// 	log.Infof("Executing: %s\n", myCommand)
	// 	// ssh.Stream("for i in {1..5}; do echo ${i}; sleep ; done; exit 2;", false)
	// } else {
	// 	log.Errorln("TODO: add ssh run...")
	// }

	return nil
}

// func RunExampleCommand() {
// 	ls := execute.ExecTask{
// 		Command: "exit 1",
// 		Shell:   true,
// 	}
// 	res, err := ls.Execute()
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Warnf("==> stdout: %q, stderr: %q, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
// }

func RunExampleCommand2() {
	ls := execute.ExecTask{
		// Command: "df",
		// Args:    []string{"-P"},
		Command:     "echo",
		Args:        []string{"${UID}"},
		Shell:       true,
		StreamStdio: false, // если true то выводит вконсоль и в Stdout
	}
	res, err := ls.Execute()
	if err != nil {
		panic(err)
	}

	log.Warnf("stdout: %q, stderr: %q, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
	if val, err := oper.ParseInt64Output(res.Stdout); err != nil {
		log.Fatalln(err.Error())
	} else {
		log.Debugf("Result: %v", val)
	}
}

// func RunExampleCommand3() {
// 	cmd := execute.ExecTask{
// 		Command:     "docker",
// 		Args:        []string{"version"},
// 		StreamStdio: false,
// 	}

// 	res, err := cmd.Execute()
// 	if err != nil {
// 		panic(err)
// 	}

// 	if res.ExitCode != 0 {
// 		panic("Non-zero exit code: " + res.Stderr)
// 	}

// 	fmt.Printf("stdout: %s, stderr: %s, exit-code: %d\n", res.Stdout, res.Stderr, res.ExitCode)
// }
