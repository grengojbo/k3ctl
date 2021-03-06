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

	// for _, secret := range options.Secrets {
	// 	if err := CreateSecret(secret); err != nil {
	// 		return err
	// 	}
	// }

	// chart := fmt.Sprintf("%s/%s", options.Helm.Repo, options.Helm.Name)
	args := []string{"upgrade", "--install", options.Helm.Name, options.Helm.Repo, "--namespace", options.Helm.Namespace, "--kubeconfig", options.KubeconfigPath, "--kube-context", options.ClusterName}
	if len(options.Helm.Version) > 0 {
		args = append(args, "--version", options.Helm.Version)
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
		if len(options.Helm.Values[k]) == 0 {
			args = append(args, "--set")
			args = append(args, fmt.Sprintf("'%s=%s'", k, v))
		}
	}

	for k, v := range options.Helm.Values {
		args = append(args, "--set")
		args = append(args, fmt.Sprintf("'%s=%s'", k, v))
	}

	if len(args) > 0 {
		if options.Wait {
			args = append(args, "--wait")
		}
		if options.DryRun {
			args = append(args, "--dry-run")
		}

		// Dependency Update
		if options.Helm.DependencyUpdate {
			command := fmt.Sprintf("helm dependency update %s/%s", options.Helm.Repo, options.Helm.Name)
			stdOut, errOut, err := k3s.RunLocalCommand(command, false, options.DryRun)
			if err != nil {
				log.Errorf("[Helm3Upgrade:RunLocalCommand] %v\n", err.Error())
				log.Errorf("errOut: %s", errOut)
			} else {
				// log.Debug(stdOut)
				log.Debugf("%s", stdOut)
			}
		}

		argsSt := strings.Join(args, " ")
		command := fmt.Sprintf("helm %s", argsSt)
		// log.Debugf("Command: %s\n", command)
		stdOut, errOut, err := k3s.RunLocalCommand(command, false, options.DryRun)
		if err != nil {
			log.Errorf("[Helm3Upgrade:RunLocalCommand] %v\n", err.Error())
		} else {
			// log.Debug(stdOut)
			log.Debugf("%s", stdOut)
			if len(errOut) > 0 {
				log.Errorf("errOut: %s", errOut)
			}
			log.Infof("Release \"%s\" has been upgraded.", options.Helm.Name)
		}
	}
	return nil
}

// CreateSecret kubectl create secret
func CreateSecret(secret k3sv1alpha1.K8sSecret, kubeConfigPath string, clusterName string, dryRun bool) error {

	args := []string{"-n", secret.Namespace}
	secretsData, err := flattenSecretData(secret.SecretsData)
	if err != nil {
		return err
	}
	args = append(args, secretsData...)

	argsSt := strings.Join(args, " ")
	command := fmt.Sprintf(types.SecretCreateCommand, secret.Type, secret.Name, argsSt, kubeConfigPath, clusterName)
	log.Infof("Create Secret: %s...", secret.Name)
	stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
	if err != nil {
		log.Errorf("unable to create secret  %v", err.Error())
	} else {
		log.Infof("[CreateSecret] %s", stdOut)
	}
	return nil
}

// DeleteSecret
func DeleteSecret(secret k3sv1alpha1.K8sSecret, kubeConfigPath string, clusterName string, dryRun bool) error {
	command := fmt.Sprintf(types.SecretListCommand, secret.Namespace, kubeConfigPath, clusterName)
	stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
	if err != nil {
		log.Errorf("%s", err.Error())
	} else {
		lines := strings.Split(string(stdOut), "\n")
		if _, ok := util.Find(lines, secret.Name); ok {
			command = fmt.Sprintf(types.SecretDeleteCommand, secret.Name, secret.Namespace, kubeConfigPath, clusterName)
			stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
			if err != nil {
				log.Errorf("unable to delete secret  %v", err.Error())
			} else {
				log.Infof("[DeleteSecret] %s", stdOut)
			}
		}
	}
	return nil
}

// DeleteHelmReleases - Delete Helm Releases
func DeleteHelmReleases(releases []k3sv1alpha1.HelmInterfaces, kubeconfigPath string, clusterName string, dryRun bool) {
	for _, release := range releases {
		command := fmt.Sprintf(types.HelmDeleteCommand, release.Name, release.Namespace, kubeconfigPath, clusterName)
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
func AddHelmRepo(repos []k3sv1alpha1.HelmRepo, kubeconfigPath string, clusterName string, updateRepo bool, dryRun bool) {
	command := fmt.Sprintf(types.HelmRepoListCommand, kubeconfigPath, clusterName)
	stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
	if err != nil {
		log.Errorf("[RunLocalCommand] %v\n", err.Error())
	}
	lines := strings.Split(string(stdOut), "\n")
	for _, repo := range repos {
		// log.Warnf("repo: %v", repo)
		if _, ok := util.Find(lines, repo.Name); !ok {
			log.Infof("Add Helm Repo: %s for %s", repo.Name, repo.Repo)
			command := fmt.Sprintf(types.HelmRepoAddCommand, repo.Name, repo.Url, kubeconfigPath, clusterName)
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
		command = fmt.Sprintf(types.HelmRepoUpdateCommand, kubeconfigPath, clusterName)
		stdOut, _, err = k3s.RunLocalCommand(command, false, dryRun)
		if err != nil {
			log.Errorf("[RunLocalCommand] %v\n", err.Error())
		} else {
			log.Infoln("[AddHelmRepo] start updated...")
			log.Debugf("%s", stdOut)
			log.Infoln("[AddHelmRepo] finish updated")
		}
	}
}

//  CreateNamespace - Create namespace is not exits
func CreateNamespace(ns []string, kubeconfigPath string, clusterName string, dryRun bool) {
	command := fmt.Sprintf(types.NamespaceGetCommand, kubeconfigPath, clusterName)
	stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
	if err != nil {
		log.Errorf("[RunLocalCommand] %v\n", err.Error())
	}
	lines := strings.Split(string(stdOut), "\n")
	for _, line := range ns {
		if _, ok := util.Find(lines, line); !ok {
			// log.Debugf("create namespace: %s", line)
			if line != "default" && line != "kube-system" {
				command = fmt.Sprintf(types.NamespaceCreateCommand, line, kubeconfigPath, clusterName)
				stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
				if err != nil {
					log.Errorf("[RunLocalCommand] %v\n", err.Error())
				} else {
					log.Infof("[CreateNamespace] %s", stdOut)
				}
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
			if err := util.CheckExitFile(value.Value); err != nil {
				return nil, fmt.Errorf("IS NOT file: %s", value.Value)
			}
			output = append(output, fmt.Sprintf("--from-literal=%s=%s", value.Key, value.Value))

		case types.FromFileSecret:
			output = append(output, fmt.Sprintf("--from-file=%s=%s", value.Key, value.Value))
		default:

			return nil, fmt.Errorf("could not create secret value of type %s. Please use one of [%s, %s]", value.Type, types.StringLiteralSecret, types.FromFileSecret)

		}
	}

	return output, nil
}
