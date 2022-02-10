package k3s

import (
	"fmt"
	// "time"

	// "context"

	// "github.com/cnrancher/autok3s/pkg/types"

	// yamlv3 "gopkg.in/yaml.v3"
	// v1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	// "k8s.io/client-go/kubernetes"
	// "k8s.io/client-go/rest"
	// "k8s.io/client-go/tools/clientcmd"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	oper "github.com/grengojbo/k3ctl/pkg/operator"

	// "github.com/grengojbo/k3ctl/pkg/types"

	// "github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

// ExecuteMaster TODO: delete
func ExecuteMaster(runCommand string, node *k3sv1alpha1.ContrelPlanNodes, dryRun bool) (result string, err error)  {
	if node.Bastion.Name == "local" {
		log.Infoln("Run command in localhost........")
		// stdOut, stdErr, err := RunLocalCommand(installK3scommand, true, dryRun)
		_, stdErr, err := RunLocalCommand(runCommand, true, dryRun)
		if err != nil {
			log.Fatalln(err.Error())
		} else if len(stdErr) > 0 {
			log.Errorf("stderr: %q", stdErr)
		}
		// log.Infof("stdout: %q", stdOut)
	} else {
		if node.Node.User != "root" {
			runCommand = fmt.Sprintf("sudo %s", runCommand)
		} 
		if dryRun {
			log.Warnf("Dry RUN: ssh %s@%s -p %d \"%s\"", node.Bastion.User, node.Bastion.Address, node.Bastion.SshPort, runCommand)
			return result, err
		}
		ssh := oper.SSHOperator{}
		ssh.NewSSHOperator(node.Bastion)
		stdOut, stdErr, err := ssh.Execute(runCommand)
		if err != nil {
			if len(stdErr) > 0 {
				log.Errorln(stdErr)
			}
			// log.Fatalln(err.Error())
			return "", err
		} else {
			return stdOut, err
		}
		// log.Warnln(stdOut)
		// RunExampleCommand2()
	}
	return result, err
}