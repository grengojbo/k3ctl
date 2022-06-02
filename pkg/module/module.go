package module

import (
	"fmt"
	"strings"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/k3s"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
)

func mergeFlags(existingMap map[string]string, setOverrides []string) error {
	for _, setOverride := range setOverrides {
		flag := strings.Split(setOverride, "=")
		if len(flag) != 2 {
			return fmt.Errorf("incorrect format for custom flag `%s`", setOverride)
		}
		existingMap[flag[0]] = flag[1]
	}
	return nil
}

// DeleteHelmReleases - Delete Helm Releases
func DeleteHelmReleases(releases []k3sv1alpha1.HelmInterfaces, kubeconfigPath string, dryRun bool) {
	for _, release := range releases {
		command := fmt.Sprintf("helm delete %s -n %s --kubeconfig %s", release.Name, release.Namespace, kubeconfigPath)
		log.Infof("Delete Helm Release: %s ", release.Name)
		stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
		if err != nil {
			log.Errorf("[RunLocalCommand] %v\n", err.Error())
		} else {
			log.Infof("[DeleteHelmReleases] %s", stdOut)
		}
	}
}

// AddHelmRepo - Add helm repository and update
func AddHelmRepo(repos []k3sv1alpha1.HelmInterfaces, kubeconfigPath string, updateRepo bool, dryRun bool) {
	command := fmt.Sprintf("helm repo list --kubeconfig %s -o json | jq -r '.[].name'", kubeconfigPath)
	stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
	if err != nil {
		log.Errorf("[RunLocalCommand] %v\n", err.Error())
	}
	lines := strings.Split(string(stdOut), "\n")
	for _, repo := range repos {
		if _, ok := util.Find(lines, repo.Repo); !ok {
			log.Infof("Add Helm Repo: %s for %s", repo.Repo, repo.Name)
			command := fmt.Sprintf("helm repo add %s %s --kubeconfig %s", repo.Repo, repo.Url, kubeconfigPath)
			// log.Warnf("command: %s", command)
			stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
			if err != nil {
				log.Errorf("[RunLocalCommand] %v\n", err.Error())
			} else {
				log.Infof("[AddHelmRepo] %s", stdOut)
			}
		}
	}
	if updateRepo {
		command = fmt.Sprintf("helm repo update --kubeconfig %s", kubeconfigPath)
		stdOut, _, err = k3s.RunLocalCommand(command, false, dryRun)
		if err != nil {
			log.Errorf("[RunLocalCommand] %v\n", err.Error())
		} else {
			log.Infof("[AddHelmRepo] %s", stdOut)
		}
	}
}

//  CreateNamespace - Create namespace is not exits
func CreateNamespace(ns []string, kubeconfigPath string, dryRun bool) {
	// var namespaces = []string{}

	// command := fmt.Sprintf("kubectl get ns -o=custom-columns='NAME:.metadata.name' --no-headers --kubeconfig %s --cluster=cloud", kubeconfigPath)
	command := fmt.Sprintf("kubectl get ns -o=custom-columns='NAME:.metadata.name' --no-headers --kubeconfig %s", kubeconfigPath)
	stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
	if err != nil {
		log.Errorf("[RunLocalCommand] %v\n", err.Error())
	}
	lines := strings.Split(string(stdOut), "\n")
	for _, line := range ns {
		if _, ok := util.Find(lines, line); !ok {
			// log.Debugf("create namespace: %s", line)
			command = fmt.Sprintf("kubectl create namespace %s --kubeconfig %s", line, kubeconfigPath)
			stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
			if err != nil {
				log.Errorf("[RunLocalCommand] %v\n", err.Error())
			} else {
				log.Infof("[CreateNamespace] %s", stdOut)
			}
		}
	}
	// for _, line := range lines {
	// 	row := strings.TrimSpace(line)
	// 	if len(row) > 0 {
	// 		log.Debugf("create namespace: %s" ,row)
	// 	}
	// }

	// lines := bufio.NewScanner(strings.NewReader(res.Stdout))
	// for lines.Scan() {
	// 	log.Debugf("ns: %s", lines)
	// 	// caps[lines.Text()] = true
	// }
}
