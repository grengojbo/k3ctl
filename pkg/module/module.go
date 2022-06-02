package module

import (
	"fmt"
	"strings"

	"github.com/grengojbo/k3ctl/pkg/k3s"
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

func CreateNamespace(ns string, KubeconfigPath string, dryRun bool) {
	var namespaces = []string{}

	// command := fmt.Sprintf("kubectl get ns -o=custom-columns='NAME:.metadata.name' --no-headers --kubeconfig %s --cluster=cloud", KubeconfigPath)
	command := fmt.Sprintf("kubectl get ns -o=custom-columns='NAME:.metadata.name' --no-headers --kubeconfig %s", KubeconfigPath)
	stdOut, _, err := k3s.RunLocalCommand(command, false, dryRun)
	if err != nil {
		log.Errorf("[RunLocalCommand] %v\n", err.Error())
	}
	lines := strings.Split(string(stdOut), "\n")
	for _, line := range lines {
		row := strings.TrimSpace(line)
		if len(row) > 0 {
			namespaces = append(namespaces, row)
		}
	}

	// lines := bufio.NewScanner(strings.NewReader(res.Stdout))
	// for lines.Scan() {
	// 	log.Debugf("ns: %s", lines)
	// 	// caps[lines.Text()] = true
	// }
}
