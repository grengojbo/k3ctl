package k3s

import (
	"fmt"
	"time"

	// "context"

	// "github.com/cnrancher/autok3s/pkg/types"

	// yamlv3 "gopkg.in/yaml.v3"
	// v1 "k8s.io/api/core/v1"
	// metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	// "github.com/grengojbo/k3ctl/pkg/util"
)

// https://github.com/cnrancher/autok3s/blob/a9468516b89009a0d5488cdaad4eb0eb5370cedc/pkg/cluster/cluster.go#L5

func buildConfigFromFlags(context, kubeconfigPath string) (*rest.Config, error) {
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfigPath},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

// GetClusterConfig generate kube config.
func GetClusterConfig(context string) (*kubernetes.Clientset, error) {
	// config, err := buildConfigFromFlags(name, kubeconfig)
	config, err := BuildKubeConfigFromFlags(context)
	if err != nil {
		return nil, err
	}
	config.Timeout = 15 * time.Second
	c, err := kubernetes.NewForConfig(config)
	return c, err
}

// GetKubeconfig подключаемся к мастеру и скачиваем kubeconfig
func GetKubeconfig(masters []k3sv1alpha1.ContrelPlanNodes, dryRun bool) (kubeconfig string, err error) {
	if dryRun {
		return "", nil
	}
	if len(masters) == 0 {
		return "", fmt.Errorf("Is NOT set control plane nodes")
	}
	// runCommand := "cat /var/lib/rancher/k3s/server/node-token"
	for _, item := range masters {
		kubeconfig, err = ExecuteMaster(types.CatCfgCommand, &item, dryRun)
		// if err != nil {
		// 	log.Errorln(err)
		// }
		return kubeconfig, err
	}
	return kubeconfig, err
}

// GetClusterStatus get cluster status using cluster's /readyz API.
// func GetClusterStatus(c *kubernetes.Clientset) string {
// 	_, err := c.RESTClient().Get().Timeout(15 * time.Second).RequestURI("/readyz").DoRaw(context.TODO())
// 	if err != nil {
// 		return types.ClusterStatusStopped
// 	}
// 	return types.ClusterStatusRunning
// }

// // GetClusterVersion get kube cluster version.
// func GetClusterVersion(c *kubernetes.Clientset) string {
// 	v, err := c.DiscoveryClient.ServerVersion()
// 	if err != nil {
// 		// return types.ClusterStatusUnknown
// 	}
// 	return v.GitVersion
// }
