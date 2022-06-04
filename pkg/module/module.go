package module

import (
	"fmt"
	"strings"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/k3s"
	"github.com/grengojbo/k3ctl/pkg/types"
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

// Helm3Upgrade - Install or update HELM Chart
func Helm3Upgrade(options *k3sv1alpha1.HelmOptions) (err error) {
	if options.CreateNamespace {
		log.Warnln("[Helm3Upgrade] TODO: CreateNamespace")
	}

	for _, secret := range options.Secrets {
		if err := CreateSecret(secret); err != nil {
			return err
		}
	}

	chart := fmt.Sprintf("%s/%s", options.Helm.Repo, options.Helm.Name)
	args := []string{"upgrade", "--install", options.Helm.Name, chart, "--namespace", options.Helm.Namespace, "--kubeconfig", options.KubeconfigPath}
	if len(options.Helm.Version) > 0 {
		args = append(args, "--version", options.Helm.Version)
	}
	if options.Wait {
		args = append(args, "--wait")
	}
	if options.DryRun {
		args = append(args, "--dry-run")
	}

	if len(options.Helm.ValuesFile) > 0 {
		args = append(args, "--values")
		args = append(args, options.Helm.ValuesFile)

	}
	// fmt.Println("VALUES", values)
	// if len(values) > 0 {
	// 	args = append(args, "--values")
	// 	if !strings.HasPrefix(values, "/") {
	// 		args = append(args, path.Join(basePath, values))
	// 	} else {
	// 		args = append(args, values)
	// 	}
	// }

	for k, v := range options.Overrides {
		args = append(args, "--set")
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	for k, v := range options.Helm.Values {
		args = append(args, "--set")
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}

	if len(args) > 0 {
		argsSt := strings.Join(args, " ")
		command := fmt.Sprintf("helm %s", argsSt)
		// log.Debugf("Command: %s\n", command)
		stdOut, _, err := k3s.RunLocalCommand(command, false, options.DryRun)
		if err != nil {
			log.Errorf("[Helm3Upgrade:RunLocalCommand] %v\n", err.Error())
		} else {
			log.Debug(stdOut)
			log.Infof("Release \"%s\" has been upgraded.", options.Helm.Name)
		}
	}
	return nil
}

// CreateSecret kubectl create secret
func CreateSecret(secret k3sv1alpha1.K8sSecret) error {
	log.Warnln("TODO: CreateSecret :)")
	// secretData, err := flattenSecretData(secret.SecretData)
	// if err != nil {
	// 	return err
	// }

	// args := []string{"-n", secret.Namespace, "create", "secret", secret.Type, secret.Name}
	// args = append(args, secretData...)

	// res, secretErr := KubectlTask(args...)

	// if secretErr != nil {
	// 	return secretErr
	// }
	// if res.ExitCode != 0 {
	// 	fmt.Printf("[Warning] unable to create secret %s, may already exist: %s", secret.Name, res.Stderr)
	// }

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

func flattenSecretData(data []k3sv1alpha1.SecretsData) ([]string, error) {
	var output []string

	for _, value := range data {
		switch value.Type {
		case types.StringLiteralSecret:
			output = append(output, fmt.Sprintf("--from-literal=%s=%s", value.Key, value.Value))

		case types.FromFileSecret:
			output = append(output, fmt.Sprintf("--from-file=%s=%s", value.Key, value.Value))
		default:

			return nil, fmt.Errorf("could not create secret value of type %s. Please use one of [%s, %s]", value.Type, types.StringLiteralSecret, types.FromFileSecret)

		}
	}

	return output, nil
}
