package k3s

import (
	"fmt"

	execute "github.com/alexellis/go-execute/pkg/v1"
	operator "github.com/alexellis/k3sup/pkg/operator"

	oper "github.com/grengojbo/k3ctl/pkg/operator"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
	// "k8s.io/apimachinery/pkg/util/errors"
)

// https://github.com/alexellis/k3sup/blob/master/cmd/install.go

var kubeconfig []byte

// type K3sExecOptions struct {
// 	Datastore           string
// 	ExtraArgs           []string
// 	FlannelIPSec        bool
// 	NoExtras            bool
// 	LoadBalancer        *k3sv1alpha1.LoadBalancer
// 	Networking          *k3sv1alpha1.Networking
// 	Ingress             string
// 	DisableLoadbalancer bool
// 	DisableIngress      bool
// 	SELinux             bool
// 	Rootless            bool
// 	SecretsEncryption   bool
// }

// type K3sIstallOptions struct {
// 	ExecString   string
// 	LoadBalancer string
// 	Ingress      string
// 	CNI          string
// 	Backend      string
// 	K3sVersion   string
// 	K3sChannel   string
// 	Node         *k3sv1alpha1.Node
// }

// func MakeAgentInstallExec(apiServerAddres string, token string, options K3sExecOptions) K3sIstallOptions {
// 	// extraArgs := []string{}
// 	// curl -sfL https://get.k3s.io | K3S_URL='https://<IP>6443' K3S_TOKEN='<TOKEN>' INSTALL_K3S_CHANNEL='stable' sh -s - --node-taint key=value:NoExecute
// 	k3sIstallOptions := K3sIstallOptions{}

// 	installExec := fmt.Sprintf(types.JoinAgentCommand, apiServerAddres, options.Networking.APIServerPort, token)
// 	k3sIstallOptions.ExecString = installExec
// 	return k3sIstallOptions
// }

// GetAgentToken TODO: delete подключаемся к мастеру и получает токен для подключения агента
// func GetAgentToken(masters []k3sv1alpha1.ContrelPlanNodes, dryRun bool) (token string, err error) {
// 	if len(masters) == 0 {
// 		return "", fmt.Errorf("Is NOT set control plane nodes")
// 	}
// 	// runCommand := "cat /var/lib/rancher/k3s/server/node-token"
// 	runCommand := types.CatTokenCommand
// 	for _, item := range masters {
// 		if item.Bastion.Name == "local" {
// 			log.Infoln("Run command in localhost........")
// 			// stdOut, stdErr, err := RunLocalCommand(installK3scommand, true, dryRun)
// 			_, stdErr, err := RunLocalCommand(runCommand, true, dryRun)
// 			if err != nil {
// 				log.Fatalln(err.Error())
// 			} else if len(stdErr) > 0 {
// 				log.Errorf("stderr: %q", stdErr)
// 			}
// 			// log.Infof("stdout: %q", stdOut)
// 		} else {
// 			if item.Node.User != "root" {
// 				runCommand = fmt.Sprintf("sudo %s", runCommand)
// 			}
// 			ssh := oper.SSHOperator{}
// 			ssh.NewSSHOperator(item.Bastion)
// 			stdOut, stdErr, err := ssh.Execute(runCommand)
// 			if err != nil {
// 				log.Errorln(stdErr)
// 				// log.Fatalln(err.Error())
// 			} else {
// 				return stdOut, err
// 			}
// 			// log.Warnln(stdOut)
// 			// RunExampleCommand2()
// 		}
// 	}
// 	return token, err
// }

// RunK3sCommand Выполняем команды по SSH или локально
// TODO: translate
// func RunK3sCommand(bastion *k3sv1alpha1.BastionNode, installk3sExec *K3sIstallOptions, dryRun bool) error {
// 	installStr := util.CreateVersionStr(installk3sExec.K3sVersion, installk3sExec.K3sChannel)
// 	installK3scommand := fmt.Sprintf("%s | %s %s sh -\n", types.K3sGetScript, installk3sExec.ExecString, installStr)

// 	if len(installk3sExec.K3sChannel) == 0 && len(installk3sExec.K3sVersion) == 0 {
// 		return errors.New("Set kubernetesVersion or channel (Release channel: stable, latest, or i.e. v1.19)")
// 	}

// 	log.Infof("KubernetesVersion: %s (K3sChannel: %s)", installk3sExec.K3sVersion, installk3sExec.K3sChannel)
// 	if installk3sExec.Node.Role == "master" {
// 		log.Infof("CNI: %s Backend: %s", installk3sExec.CNI, installk3sExec.Backend)
// 		log.Infof("LoadBalancer: %s", installk3sExec.LoadBalancer)
// 		if len(installk3sExec.Ingress) > 0 {
// 			log.Infof("Ingress Controllers: %s", installk3sExec.Ingress)
// 		} else {
// 			log.Infoln("Ingress Controllers: default k3s Traefik")
// 		}
// 	}
// 	log.Warnf("Bastion %s host: %s (ssh port: %d key: %s)", bastion.Name, bastion.Address, bastion.SshPort, bastion.SSHAuthorizedKey)
// 	log.Debugln("--------------------------------------------")

// 	// sudoPrefix := ""
// 	// if useSudo {
// 	// 	sudoPrefix = "sudo "
// 	// }
// 	// getConfigcommand := fmt.Sprintf(sudoPrefix + "cat /etc/rancher/k3s/k3s.yaml\n")

// 	if bastion.Name == "local" {
// 		log.Infoln("Run command in localhost........")
// 		// stdOut, stdErr, err := RunLocalCommand(installK3scommand, true, dryRun)
// 		_, stdErr, err := RunLocalCommand(installK3scommand, true, dryRun)
// 		if err != nil {
// 			log.Fatalln(err.Error())
// 		} else if len(stdErr) > 0 {
// 			log.Errorf("stderr: %q", stdErr)
// 		}
// 		// log.Infof("stdout: %q", stdOut)
// 	} else {
// 		if err := RunSshCommand(installK3scommand, bastion, true, dryRun); err != nil {
// 			log.Fatalln(err.Error())
// 		}
// 	}
//
// 	// RunExampleCommand()
// 	// RunExampleCommand2()
// 	return nil
// }

// RunLocalCommand выполнение комманд на локальном хосте TODO: tranclate
func RunLocalCommand(myCommand string, sudo bool, dryRun bool) (stdOut []byte, stdErr []byte, err error) {

	sudoPrefix := ""
	if sudo {
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
	}

	command := fmt.Sprintf("%s%s", sudoPrefix, myCommand)
	operator := operator.ExecOperator{}
	if dryRun {
		log.Infof("DRY-RUN %s", command)
	} else {
		log.Debugf("Executing: %s\n", command)

		res, err := operator.ExecuteStdio(command, false)
		// res, err := operator.Execute("pwd")
		if err != nil {
			return stdOut, stdErr, err
		}

		if len(res.StdErr) > 0 {
			stdErr = res.StdErr
		}
		if len(res.StdOut) > 0 {
			stdOut = res.StdOut
		}
	}

	return stdOut, stdErr, nil
}

func RunLocalCommandOld(myCommand string, saveKubeconfig bool, dryRun bool) (stdOut []byte, stdErr []byte, err error) {

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
		log.Fatalln("TODO: Run command in from bastion host ........")
		// log.Errorln("TODO: Run command in from bastion host ........")
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
