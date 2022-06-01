/*
Copyright © 2020-2021 The k3d Author(s)

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/
package k3s

import (
	// "bytes"
	// "context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	// "io"
	// "os"
	// "path/filepath"
	"time"

	// l "github.com/rancher/k3d/v5/pkg/logger"
	// "github.com/rancher/k3d/v5/pkg/runtimes"
	// k3d "github.com/rancher/k3d/v5/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

// WriteKubeConfigOptions provide a set of options for writing a KubeConfig file
type WriteKubeConfigOptions struct {
	UpdateExisting       bool
	UpdateCurrentContext bool
	OverwriteExisting    bool
}

// KubeconfigGetWrite ...
// 1. fetches the KubeConfig from the first server node retrieved for a given cluster
// 2. modifies it by updating some fields with cluster-specific information
// 3. writes it to the specified output
// func KubeconfigGetWrite(ctx context.Context, runtime runtimes.Runtime, cluster *k3d.Cluster, output string, writeKubeConfigOptions *WriteKubeConfigOptions) (string, error) {

// 	// get kubeconfig from cluster node
// 	kubeconfig, err := KubeconfigGet(ctx, runtime, cluster)
// 	if err != nil {
// 		return output, fmt.Errorf("failed to get kubeconfig for cluster '%s': %w", cluster.Name, err)
// 	}

// 	// empty output parameter = write to default
// 	if output == "" {
// 		output, err = KubeconfigGetDefaultPath()
// 		if err != nil {
// 			return output, fmt.Errorf("failed to get default kubeconfig path: %w", err)
// 		}
// 	}

// 	// simply write to the output, ignoring existing contents
// 	if writeKubeConfigOptions.OverwriteExisting || output == "-" {
// 		return output, KubeconfigWriteToPath(ctx, kubeconfig, output)
// 	}

// 	// load config from existing file or fail if it has non-kubeconfig contents
// 	var existingKubeConfig *clientcmdapi.Config
// 	firstRun := true
// 	for {
// 		existingKubeConfig, err = clientcmd.LoadFromFile(output) // will return an empty config if file is empty
// 		if err != nil {

// 			// the output file does not exist: try to create it and try again
// 			if os.IsNotExist(err) && firstRun {
// 				log.Debugf("Output path '%s' doesn't exist, trying to create it...", output)

// 				// create directory path
// 				if err := os.MkdirAll(filepath.Dir(output), 0755); err != nil {
// 					return output, fmt.Errorf("failed to create output directory '%s': %w", filepath.Dir(output), err)
// 				}

// 				// try create output file
// 				f, err := os.Create(output)
// 				if err != nil {
// 					return output, fmt.Errorf("failed to create output file '%s': %w", output, err)
// 				}
// 				f.Close()

// 				// try again, but do not try to create the file this time
// 				firstRun = false
// 				continue
// 			}
// 			return output, fmt.Errorf("failed to open output file '%s' or it's not a kubeconfig: %w", output, err)
// 		}
// 		break
// 	}

// 	// update existing kubeconfig, but error out if there are conflicting fields but we don't want to update them
// 	return output, KubeconfigMerge(ctx, kubeconfig, existingKubeConfig, output, writeKubeConfigOptions.UpdateExisting, writeKubeConfigOptions.UpdateCurrentContext)

// }

// // KubeconfigGet grabs the kubeconfig file from /output from a server node container,
// // modifies it by updating some fields with cluster-specific information
// // and returns a Config object for further processing
// func KubeconfigGet(ctx context.Context, runtime runtimes.Runtime, cluster *k3d.Cluster) (*clientcmdapi.Config, error) {
// 	// get all server nodes for the selected cluster
// 	// TODO: getKubeconfig: we should make sure, that the server node we're trying to fetch from is actually running
// 	serverNodes, err := runtime.GetNodesByLabel(ctx, map[string]string{k3d.LabelClusterName: cluster.Name, k3d.LabelRole: string(k3d.ServerRole)})
// 	if err != nil {
// 		return nil, fmt.Errorf("runtime failed to get server nodes for cluster '%s': %w", cluster.Name, err)
// 	}
// 	if len(serverNodes) == 0 {
// 		return nil, fmt.Errorf("didn't find any server node for cluster '%s'", cluster.Name)
// 	}

// 	// prefer a server node, which actually has the port exposed
// 	var chosenServer *k3d.Node
// 	chosenServer = nil
// 	APIPort := k3d.DefaultAPIPort
// 	APIHost := k3d.DefaultAPIHost

// 	for _, server := range serverNodes {
// 		if _, ok := server.RuntimeLabels[k3d.LabelServerAPIPort]; ok {
// 			chosenServer = server
// 			APIPort = server.RuntimeLabels[k3d.LabelServerAPIPort]
// 			if _, ok := server.RuntimeLabels[k3d.LabelServerAPIHost]; ok {
// 				APIHost = server.RuntimeLabels[k3d.LabelServerAPIHost]
// 			}
// 			break
// 		}
// 	}

// 	if chosenServer == nil {
// 		chosenServer = serverNodes[0]
// 	}
// 	// get the kubeconfig from the first server node
// 	reader, err := runtime.GetKubeconfig(ctx, chosenServer)
// 	if err != nil {
// 		return nil, fmt.Errorf("runtime failed to pull kubeconfig from node '%s': %w", chosenServer.Name, err)
// 	}
// 	defer reader.Close()

// 	readBytes, err := io.ReadAll(reader)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to read kubeconfig file: %w", err)
// 	}

// 	// drop the first 512 bytes which contain file metadata/control characters
// 	// and trim any NULL characters
// 	trimBytes := bytes.Trim(readBytes[512:], "\x00")

// 	/*
// 	 * Modify the kubeconfig
// 	 */
// 	kc, err := clientcmd.Load(trimBytes)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to parse kubeconfig: %w", err)
// 	}

// 	// update the server URL
// 	kc.Clusters["default"].Server = fmt.Sprintf("https://%s:%s", APIHost, APIPort)

// 	// rename user from default to admin
// 	newAuthInfoName := fmt.Sprintf("admin@%s-%s", k3d.DefaultObjectNamePrefix, cluster.Name)
// 	kc.AuthInfos[newAuthInfoName] = kc.AuthInfos["default"]
// 	delete(kc.AuthInfos, "default")

// 	// rename cluster from default to clustername
// 	newClusterName := fmt.Sprintf("%s-%s", k3d.DefaultObjectNamePrefix, cluster.Name)
// 	kc.Clusters[newClusterName] = kc.Clusters["default"]
// 	delete(kc.Clusters, "default")

// 	// rename context from default to clustername
// 	newContextName := fmt.Sprintf("%s-%s", k3d.DefaultObjectNamePrefix, cluster.Name)
// 	kc.Contexts[newContextName] = kc.Contexts["default"]
// 	delete(kc.Contexts, "default")

// 	// update context with new values for cluster and user
// 	kc.Contexts[newContextName].AuthInfo = newAuthInfoName
// 	kc.Contexts[newContextName].Cluster = newClusterName

// 	// set current-context to new context name
// 	kc.CurrentContext = newContextName

// 	log.Tracef("Modified Kubeconfig: %+v", kc)

// 	return kc, nil
// }

// // KubeconfigWriteToPath takes a kubeconfig and writes it to some path, which can be '-' for os.Stdout
// func KubeconfigWriteToPath(ctx context.Context, kubeconfig *clientcmdapi.Config, path string) error {
// 	var output *os.File
// 	defer output.Close()
// 	var err error

// 	if path == "-" {
// 		output = os.Stdout
// 	} else {
// 		output, err = os.Create(path)
// 		if err != nil {
// 			return fmt.Errorf("failed to create file '%s': %w", path, err)
// 		}
// 		defer output.Close()
// 	}

// 	kubeconfigBytes, err := clientcmd.Write(*kubeconfig)
// 	if err != nil {
// 		return fmt.Errorf("failed to write kubeconfig: %w", err)
// 	}

// 	_, err = output.Write(kubeconfigBytes)
// 	if err != nil {
// 		return fmt.Errorf("failed to write file '%s': %w", output.Name(), err)
// 	}

// 	log.Debugf("Wrote kubeconfig to '%s'", output.Name())

// 	return nil

// }

// func KubeconfigMerge(srcKubeConfig string) (newKubeConfig *clientcmdapi.Config) {
// 	log.Infof("%v", clientcmd.RecommendedConfigPathEnvVar)
// 	newKubeConfig = clientcmdapi.NewConfig()

// 	// k := clientcmd.NewDefaultClientConfigLoadingRules()
// 	// newKubeConfig = clientcmdapi.NewCmdKubeconfig()
// 	return nil
// }

// KubeconfigMerge merges a new kubeconfig into an existing kubeconfig and returns the result
func KubeconfigMerge(newKubeConfig *clientcmdapi.Config, existingKubeConfig *clientcmdapi.Config, outPath string, opts WriteKubeConfigOptions) error {

	log.Tracef("Merging new Kubeconfig:\n%+v\n>>> into existing Kubeconfig:\n%+v", newKubeConfig, existingKubeConfig)

	// Overwrite values in existing kubeconfig
	for k, v := range newKubeConfig.Clusters {
		if _, ok := existingKubeConfig.Clusters[k]; ok {
			if !opts.OverwriteExisting {
				return fmt.Errorf("Cluster '%s' already exists in target KubeConfig", k)
			}
		}
		existingKubeConfig.Clusters[k] = v
	}

	for k, v := range newKubeConfig.AuthInfos {
		if _, ok := existingKubeConfig.AuthInfos[k]; ok {
			if !opts.OverwriteExisting {
				return fmt.Errorf("AuthInfo '%s' already exists in target KubeConfig", k)
			}
		}
		existingKubeConfig.AuthInfos[k] = v
	}

	for k, v := range newKubeConfig.Contexts {
		if _, ok := existingKubeConfig.Contexts[k]; ok && !opts.OverwriteExisting {
			return fmt.Errorf("Context '%s' already exists in target KubeConfig", k)
		}
		existingKubeConfig.Contexts[k] = v
	}

	// Set current context if it's
	// a) empty
	// b) not empty, but we want to update it
	if existingKubeConfig.CurrentContext == "" {
		opts.UpdateCurrentContext = true
	}
	if opts.UpdateCurrentContext {
		log.Debugf("Setting new current-context '%s'", newKubeConfig.CurrentContext)
		existingKubeConfig.CurrentContext = newKubeConfig.CurrentContext
	}

	return KubeconfigWrite(existingKubeConfig, outPath)
}

// KubeconfigWrite writes a kubeconfig to a path atomically
func KubeconfigWrite(kubeconfig *clientcmdapi.Config, path string) error {
	tempPath := fmt.Sprintf("%s.k3s_%s", path, time.Now().Format("20060102_150405.000000"))
	if err := clientcmd.WriteToFile(*kubeconfig, tempPath); err != nil {
		return fmt.Errorf("failed to write merged kubeconfig to temporary file '%s': %w", tempPath, err)
	}

	// Move temporary file over existing KubeConfig
	if err := os.Rename(tempPath, path); err != nil {
		return fmt.Errorf("failed to overwrite existing KubeConfig '%s' with new kubeconfig '%s': %w", path, tempPath, err)
	}

	log.Debugf("Wrote kubeconfig to '%s'", path)

	return nil
}

// KubeconfigTmpWrite write kubeconfig to tmp file
func KubeconfigTmpWrite(kubeconfig *clientcmdapi.Config) (path string, err error) {
	// path = fmt.Sprintf("%s.k3s_%s", path, time.Now().Format("20060102_150405.000000"))
	tmpPath, err := ioutil.TempFile("", "k3s_*.yaml")
	if err != nil {
		log.Fatal(err)
	}
	path = tmpPath.Name()
	if err := clientcmd.WriteToFile(*kubeconfig, path); err != nil {
		return path, fmt.Errorf("failed to write merged kubeconfig to temporary file '%s': %w", path, err)
	}
	// log.Debugf("Wrote kubeconfig tmp to '%s'", path)
	return path, nil
}

// LoadKubeconfig
func LoadKubeconfig(cfg, ip, context string, opts WriteKubeConfigOptions) (kubeconfig *clientcmdapi.Config, err error) {
	// log.Debugf("cfg: %v", cfg)
	if len(cfg) == 0 {
		return nil, fmt.Errorf("[LoadKubeconfig] %s", "Is not kubeconfig")
	}
	// kubeconfig, err = clientcmd.LoadFromFile(absPath)
	kubeconfig, err = clientcmd.Load([]byte(cfg))
	if err != nil {
		// log.Errorf("[LoadKubeconfig] %v", err.Error())
		return nil, err
	}

	// update the server URL
	kubeconfig.Clusters["default"].Server = ip

	// rename user from default to admin
	newAuthInfoName := fmt.Sprintf("admin@%s", context)
	kubeconfig.AuthInfos[newAuthInfoName] = kubeconfig.AuthInfos["default"]
	delete(kubeconfig.AuthInfos, "default")

	// rename cluster from default to clustername
	kubeconfig.Clusters[context] = kubeconfig.Clusters["default"]
	delete(kubeconfig.Clusters, "default")

	// rename context from default to clustername
	kubeconfig.Contexts[context] = kubeconfig.Contexts["default"]
	delete(kubeconfig.Contexts, "default")

	// update context with new values for cluster and user
	kubeconfig.Contexts[context].AuthInfo = newAuthInfoName
	kubeconfig.Contexts[context].Cluster = context

	// set current-context to new context name
	kubeconfig.CurrentContext = context

	return kubeconfig, err
}

func SaveKubeconfig(kubeconfig *clientcmdapi.Config, opts WriteKubeConfigOptions) (pathKubeConfig string, err error) {

	pathKubeConfig, existingKubeConfig, err := KubeconfigGetDefaultFile()
	if err != nil {
		// log.Errorf("-------[ERROR] defaultPathKubeconfig: %v", err.Error())
		existingKubeConfig = clientcmdapi.NewConfig()
	}
	// log.Warnf("------- defaultPathKubeconfig: %s", existingKubeConfig)
	err = KubeconfigMerge(kubeconfig, existingKubeConfig, pathKubeConfig, opts)
	// path, err := os.Stat(defaultPathKubeconfig)
	// if err !=nil {
	// 	log.Errorf("-------[ERROR] File exits: %v", err.Error())
	// }
	// log.Warnf("------- file: %s", path.Name())
	return pathKubeConfig, err
}

// RemoveCfg removes kubectl config file.
func RemoveCfg(context string) (err error) {
	pathKubeConfig, existingKubeConfig, err := KubeconfigGetDefaultFile()
	if err != nil {
		// log.Errorf("-------[ERROR] defaultPathKubeconfig: %v", err.Error())
		existingKubeConfig = clientcmdapi.NewConfig()
	}
	if existingKubeConfig.CurrentContext == context {
		existingKubeConfig.CurrentContext = ""
	}

	delete(existingKubeConfig.Contexts, context)
	delete(existingKubeConfig.Clusters, context)
	delete(existingKubeConfig.AuthInfos, context)

	// the auth info entry associated with the context needs to be deleted.
	for key := range existingKubeConfig.AuthInfos {
		if strings.Contains(key, fmt.Sprintf("@%s", context)) {
			delete(existingKubeConfig.AuthInfos, key)
		}
	}

	// p.Log.Debugf("Deleteted kubeconfig: %s", pathKubeConfig)
	// v, _ := yaml.Marshal(existingKubeConfig)
	// p.Log.Debugf("-----------------------\n%s", v)
	if err := clientcmd.WriteToFile(*existingKubeConfig, pathKubeConfig); err != nil {
		return fmt.Errorf("failed to write merged kubeconfig to temporary file '%s': %w", pathKubeConfig, err)
	}
	return nil
}

// / SaveCfg save kube config file.
func SaveCfg(cfg, ip, context string, opts WriteKubeConfigOptions) (pathKubeConfig string, err error) {
	kubeconfig, err := LoadKubeconfig(cfg, ip, context, opts)
	if err != nil {
		return pathKubeConfig, err
	}

	pathKubeConfig, existingKubeConfig, err := KubeconfigGetDefaultFile()
	if err != nil {
		existingKubeConfig = clientcmdapi.NewConfig()
	}
	err = KubeconfigMerge(kubeconfig, existingKubeConfig, pathKubeConfig, opts)
	return pathKubeConfig, err
}

func SwitchContext(ctx, kubeconfig string) (err error) {
	// config := LoadConfigFromPath(kubeconfig)
	// rawConfig, err := config.RawConfig()
	// if err != nil {
	// 	return err
	// }
	// if rawConfig.Contexts[ctx] == nil {
	// 	return fmt.Errorf("context %s doesn't exists", ctx)
	// }
	// rawConfig.CurrentContext = ctx
	// err = clientcmd.ModifyConfig(clientcmd.NewDefaultPathOptions(), rawConfig, true)
	return
}

func BuildKubeConfigFromFlags(context string) (*rest.Config, error) {
	var err error
	var path string

	pathKubeConfig := viper.GetString("kubeconfig")
	if len(pathKubeConfig) > 0 {
		path, _ = filepath.Abs(util.ExpandPath(pathKubeConfig))
	} else {
		path, err = KubeconfigGetDefaultPath()
		if err != nil {
			return nil, fmt.Errorf("failed to get default kubeconfig path: %w", err)
		}
	}
	return clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: path},
		&clientcmd.ConfigOverrides{
			CurrentContext: context,
		}).ClientConfig()
}

// KubeconfigGetDefaultFile loads the default KubeConfig file
func KubeconfigGetDefaultFile() (string, *clientcmdapi.Config, error) {
	var err error
	var path string

	pathKubeConfig := viper.GetString("kubeconfig")
	if len(pathKubeConfig) > 0 {
		path, _ = filepath.Abs(util.ExpandPath(pathKubeConfig))
	} else {
		path, err = KubeconfigGetDefaultPath()
		if err != nil {
			return "", nil, fmt.Errorf("failed to get default kubeconfig path: %w", err)
		}
	}
	log.Debugf("Using default kubeconfig '%s'", path)
	newKubeConfig, err := clientcmd.LoadFromFile(path)
	return path, newKubeConfig, err
}

// KubeconfigGetDefaultPath returns the path of the default kubeconfig, but errors if the KUBECONFIG env var specifies more than one file
func KubeconfigGetDefaultPath() (string, error) {
	defaultKubeConfigLoadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	if len(defaultKubeConfigLoadingRules.GetLoadingPrecedence()) > 1 {
		return "", fmt.Errorf("multiple kubeconfigs specified via KUBECONFIG env var: Please reduce to one entry, unset KUBECONFIG or explicitly choose an output")
	}
	return defaultKubeConfigLoadingRules.GetDefaultFilename(), nil
}

// // KubeconfigRemoveClusterFromDefaultConfig removes a cluster's details from the default kubeconfig
// func KubeconfigRemoveClusterFromDefaultConfig(ctx context.Context, cluster *k3d.Cluster) error {
// 	defaultKubeConfigPath, err := KubeconfigGetDefaultPath()
// 	if err != nil {
// 		return fmt.Errorf("failed to get default kubeconfig path: %w", err)
// 	}
// 	kubeconfig, err := KubeconfigGetDefaultFile()
// 	if err != nil {
// 		return fmt.Errorf("failed to get default kubeconfig file: %w", err)
// 	}
// 	kubeconfig = KubeconfigRemoveCluster(ctx, cluster, kubeconfig)
// 	return KubeconfigWrite(ctx, kubeconfig, defaultKubeConfigPath)
// }

// // KubeconfigRemoveCluster removes a cluster's details from a given kubeconfig
// func KubeconfigRemoveCluster(ctx context.Context, cluster *k3d.Cluster, kubeconfig *clientcmdapi.Config) *clientcmdapi.Config {
// 	clusterName := fmt.Sprintf("%s-%s", k3d.DefaultObjectNamePrefix, cluster.Name)
// 	contextName := fmt.Sprintf("%s-%s", k3d.DefaultObjectNamePrefix, cluster.Name)
// 	authInfoName := fmt.Sprintf("admin@%s-%s", k3d.DefaultObjectNamePrefix, cluster.Name)

// 	// delete elements from kubeconfig if they're present
// 	delete(kubeconfig.Contexts, contextName)
// 	delete(kubeconfig.Clusters, clusterName)
// 	delete(kubeconfig.AuthInfos, authInfoName)

// 	// set current-context to any other context, if it was set to the given cluster before
// 	if kubeconfig.CurrentContext == contextName {
// 		for k := range kubeconfig.Contexts {
// 			kubeconfig.CurrentContext = k
// 			break
// 		}
// 		// if current-context didn't change, unset it
// 		if kubeconfig.CurrentContext == contextName {
// 			kubeconfig.CurrentContext = ""
// 		}
// 	}
// 	return kubeconfig
// }
