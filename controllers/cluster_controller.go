/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/appleboy/easyssh-proxy"
	"github.com/avast/retry-go/v4"
	"github.com/go-logr/logr"
	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/k3s"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"golang.org/x/sync/syncmap"
	"gopkg.in/yaml.v2"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"sigs.k8s.io/cluster-api/api/v1alpha3"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

type ProviderBase struct {
	// types.Metadata `json:",inline"`
	// types.Status   `json:"status"`
	// types.SSH      `json:",inline"`
	CmdFlags	 	types.CmdFlags
	Cluster    	*k3sv1alpha1.Cluster
	Clientset  	*kubernetes.Clientset
	Config 			*clientcmdapi.Config
	SSH 				*easyssh.MakeConfig
	M          	*sync.Map
	Log        	*logrus.Logger
	Callbacks   map[string]*providerProcess
}

type providerProcess struct {
	ContextName string
	Event       string
	Fn          func(interface{})
}

// NewClusterFromViperConfig new base provider.
func NewClusterFromConfig(configViper *viper.Viper, cmdFlags types.CmdFlags) (providerBase *ProviderBase, err error) {
	providerBase = &ProviderBase{
		CmdFlags: cmdFlags,
		// Metadata: types.Metadata{
		// 	UI:            ui,
		// 	K3sVersion:    k3sVersion,
		// 	K3sChannel:    k3sChannel,
		// 	InstallScript: k3sInstallScript,
		// 	Cluster:       embedEtcd,
		// 	Master:        master,
		// 	Worker:        worker,
		// 	ClusterCidr:   defaultCidr,
		// 	DockerScript:  dockerInstallScript,
		// },
		// Status: types.Status{
		// 	MasterNodes: make([]types.Node, 0),
		// 	WorkerNodes: make([]types.Node, 0),
		// },
		// SSH: types.SSH{
		// 	SSHPort: "22",
		// },
		// SSH: &oper.SSH{},
		M: new(syncmap.Map),
	}

	providerBase.initLogging(&cmdFlags)
	err = providerBase.FromViperSimple(configViper)
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// добавляем в статус списое master, worker нод
	providerBase.setGroupNodes()
	// устанавливаем настройки для кластера в зависимости от количества и типов инстант
	providerBase.SetDefaulSettings()

	return providerBase, err
}

// FromViperSimple Load config from Viper
func (p *ProviderBase) FromViperSimple(config *viper.Viper) (error) {

	var cfg k3sv1alpha1.Cluster

	// determine config kind
	if config.GetString("kind") != "" && strings.ToLower(config.GetString("kind")) != "cluster" {
		return fmt.Errorf("Wrong `kind` '%s' != 'Cluster' in config file", config.GetString("kind"))
	}

	if err := config.Unmarshal(&cfg); err != nil {
		// log.Errorln("Failed to unmarshal File config")

		return err
	}
	cfg.TypeMeta.APIVersion = config.GetString("apiversion")
	cfg.TypeMeta.Kind = config.GetString("kind")

	cfg.ObjectMeta.Name = config.GetString("metadata.name")

	// if !cfg.Spec.KubeconfigOptions.SwitchCurrentContext {
		// cfg.Spec.KubeconfigOptions.SwitchCurrentContext = true
	// }

	if cfg.Spec.Networking.APIServerPort == 0 {
		cfg.Spec.Networking.APIServerPort = 6443
	}

	// CNI драйвер по умолчанию
	if len(cfg.Spec.Networking.CNI) == 0 {
		cfg.Spec.Networking.CNI = "flannel"
	}

	// host-gw default backend for flannel
	if cfg.Spec.Networking.CNI == "flannel" && len(cfg.Spec.Networking.Backend) == 0 {
		cfg.Spec.Networking.Backend = "host-gw"
	}

	if len(cfg.GetProvider()) == 0 {
		cfg.Spec.Provider = "native"
	}

	if len(cfg.Spec.ClusterToken) == 0 {
		// p.Log.Errorf("ClusterToken: %s", cfg.Spec.ClusterToken)
		cfg.Spec.ClusterToken = util.GenerateRandomString(32)
	}
	
	if len(cfg.Spec.AgentToken) == 0 {
		cfg.Spec.AgentToken = util.GenerateRandomString(32)
	}

	if len(cfg.Spec.KubeconfigOptions.ConnectType) == 0 {
		cfg.Spec.KubeconfigOptions.ConnectType = k3sv1alpha1.InternalIP
	}
	if len(cfg.Spec.KubeconfigOptions.Patch) == 0 {
		kcfg, err := k3s.KubeconfigGetDefaultPath()
		if err != nil {
			p.Log.Errorf(err.Error())
		}
		cfg.Spec.KubeconfigOptions.Patch = kcfg
	}

	p.Cluster = &cfg
	return nil
}

// SetDefaulSettings
func (p *ProviderBase) SetDefaulSettings() {
	if p.Cluster.Spec.Agents == 0 && p.Cluster.Spec.Servers == 1 {
		p.Log.Infoln("[SetDefaulSettings] TODO: Settings for single node cluster...")
	} else if p.Cluster.Spec.Agents > 0 && p.Cluster.Spec.Servers == 1 {
		p.Log.Infoln("[SetDefaulSettings] TODO: Settings for one master cluster...")
	} else if p.Cluster.Spec.Agents > 0 && p.Cluster.Spec.Servers > 1 {
		p.Log.Infoln("[SetDefaulSettings] TODO: Settings for multi master cluster...")
	}
}

// InitK3sCluster initial K3S cluster.
func (p *ProviderBase) InitK3sCluster() error {
	p.Log.Infof("[%s] executing init k3s cluster logic...", p.Cluster.GetProvider())

	// provider, err := providers.GetProvider(p.Cluster.GetProvider())
	// if err != nil {
	// 	return err
	// }

	// k3sScript := cluster.InstallScript
	// k3sMirror := cluster.Mirror
	// dockerMirror := cluster.DockerMirror

	// if len(cluster.MasterNodes) <= 0 || len(cluster.MasterNodes[0].InternalIPAddress) <= 0 {
	// 	return errors.New("[cluster] master node internal ip address can not be empty")
	// }

	// publicIP := cluster.IP
	// if cluster.IP == "" {
	// 	cluster.IP = cluster.MasterNodes[0].InternalIPAddress[0]
	// 	publicIP = cluster.MasterNodes[0].PublicIPAddress[0]
	// }
	
	// k3sOpt := k3sv1alpha1.K3sExecOptions{
	// 	// 	NoExtras:     k3sNoExtras,
	// 	ExtraArgs:           p.Cluster.Spec.K3sOptions.ExtraServerArgs,
	// 	Ingress:             p.Cluster.Spec.Addons.Ingress.Name,
	// 	// DisableLoadbalancer: p.Cluster.Spec.Options.DisableLoadbalancer,
	// 	// DisableIngress:      p.Cluster.Spec.Options.DisableIngress,
	// 	// SecretsEncryption:   p.Cluster.Spec.Options.SecretsEncryption,
	// 	// SELinux:             p.Cluster.Spec.Options.SELinux,
	// 	// Rootless:            p.Cluster.Spec.Options.Rootless,
		
	// 	Options: 						 &p.Cluster.Spec.Options,
	// 	LoadBalancer:        &p.Cluster.Spec.LoadBalancer,
	// 	Networking:          &p.Cluster.Spec.Networking,
	// 	K3sChannel: 				 p.Cluster.Spec.K3sChannel,
	// 	KubernetesVersion: 	 p.Cluster.Spec.KubernetesVersion,
	// }

	// isCluster := false
	// if len(p.Cluster.Spec.Datastore.Provider) > 0 {
	// 	if p.Cluster.Spec.Datastore.Provider == k3sv1alpha1.DatastoreEtcd {
	// 		isCluster = true
	// 	} else if datastore, err := p.Cluster.GetDatastore(); err != nil {
	// 		p.Log.Fatalln(err.Error())
	// 	} else {
	// 		k3sOpt.Datastore = datastore
	// 		p.Log.Infof("datastore connection string: %s", datastore)
	// 	}
	// }

	masters := p.GetMasterNodes()
	firstMaster := true
	for i, node := range masters {
		if node.State.Status == types.StatusMissing {
			if firstMaster {
				masters[i].State.Status = types.StatusCreating
				installk3sExec := p.MakeInstallExec()
				// installk3sExec.Node = node

				tlsSAN := p.Cluster.GetTlsSan(node, &p.Cluster.Spec.Networking)
				p.initAdditionalMaster(tlsSAN, node, &installk3sExec)
				p.LoadNodeStatus()
				firstMaster = false
			}
		} else {
			firstMaster = false
			p.Log.Warningln("[InitK3sCluster] TODO: присоединение master node к первому мастеру")
		}
	}

	p.LoadNodeStatus()

	// // append tls-sans to k3s install script:
	// // 1. appends from --tls-sans flags.
	// // 2. appends all master nodes' first public address.
	// var tlsSans string
	// p.TLSSans = append(p.TLSSans, publicIP)
	// for _, master := range cluster.MasterNodes {
	// 	if master.PublicIPAddress[0] != "" && master.PublicIPAddress[0] != publicIP {
	// 		p.TLSSans = append(p.TLSSans, master.PublicIPAddress[0])
	// 	}
	// }
	// for _, tlsSan := range p.TLSSans {
	// 	tlsSans = tlsSans + fmt.Sprintf(" --tls-san %s", tlsSan)
	// }
	// // save p.TlsSans to db.
	// cluster.TLSSans = p.TLSSans

	// masterExtraArgs := cluster.MasterExtraArgs
	// workerExtraArgs := cluster.WorkerExtraArgs

	// if cluster.DataStore != "" {
	// 	cluster.Cluster = false
	// 	masterExtraArgs += " --datastore-endpoint " + cluster.DataStore
	// }

	// if cluster.Network != "" {
	// 	masterExtraArgs += fmt.Sprintf(" --flannel-backend=%s", cluster.Network)
	// }

	// if cluster.ClusterCidr != "" {
	// 	masterExtraArgs += " --cluster-cidr " + cluster.ClusterCidr
	// }

	// p.Logger.Infof("[%s] creating k3s master-%d...", p.Provider, 1)
	// master0ExtraArgs := masterExtraArgs
	// providerExtraArgs := provider.GenerateMasterExtraArgs(cluster, cluster.MasterNodes[0])
	// if providerExtraArgs != "" {
	// 	master0ExtraArgs += providerExtraArgs
	// }
	// if cluster.Cluster {
	// 	master0ExtraArgs += " --cluster-init"
	// }

	// if err := p.initMaster(k3sScript, k3sMirror, dockerMirror, tlsSans, publicIP, master0ExtraArgs, cluster, cluster.MasterNodes[0]); err != nil {
	// 	return err
	// }
	// p.Logger.Infof("[%s] successfully created k3s master-%d", p.Provider, 1)

	// for i, master := range cluster.MasterNodes {
	// 	// skip first master nodes.
	// 	if i == 0 {
	// 		continue
	// 	}
	// 	p.Logger.Infof("[%s] creating k3s master-%d...", p.Provider, i+1)
	// 	masterNExtraArgs := masterExtraArgs
	// 	providerExtraArgs := provider.GenerateMasterExtraArgs(cluster, master)
	// 	if providerExtraArgs != "" {
	// 		masterNExtraArgs += providerExtraArgs
	// 	}
	// 	if err := p.initAdditionalMaster(k3sScript, k3sMirror, dockerMirror, tlsSans, publicIP, masterNExtraArgs, cluster, master); err != nil {
	// 		return err
	// 	}
	// 	p.Logger.Infof("[%s] successfully created k3s master-%d", p.Provider, i+1)
	// }

	// workerErrChan := make(chan error)
	// workerWaitGroupDone := make(chan bool)
	// workerWaitGroup := &sync.WaitGroup{}
	// workerWaitGroup.Add(len(cluster.WorkerNodes))

	// for i, worker := range cluster.WorkerNodes {
	// 	go func(i int, worker types.Node) {
	// 		p.Logger.Infof("[%s] creating k3s worker-%d...", p.Provider, i+1)
	// 		extraArgs := workerExtraArgs
	// 		providerExtraArgs := provider.GenerateWorkerExtraArgs(cluster, worker)
	// 		if providerExtraArgs != "" {
	// 			extraArgs += providerExtraArgs
	// 		}
	// 		p.initWorker(workerWaitGroup, workerErrChan, k3sScript, k3sMirror, dockerMirror, extraArgs, cluster, worker)
	// 		p.Logger.Infof("[%s] successfully created k3s worker-%d", p.Provider, i+1)
	// 	}(i, worker)
	// }

	// go func() {
	// 	workerWaitGroup.Wait()
	// 	close(workerWaitGroupDone)
	// }()

	// select {
	// case <-workerWaitGroupDone:
	// 	break
	// case err := <-workerErrChan:
	// 	return err
	// }

	// // get k3s cluster config.
	// cfg, err := p.execute(&cluster.MasterNodes[0], []string{catCfgCommand})
	// if err != nil {
	// 	return err
	// }

	// // merge current cluster to kube config.
	// if err := SaveCfg(cfg, publicIP, cluster.ContextName); err != nil {
	// 	return err
	// }
	// _ = os.Setenv(clientcmd.RecommendedConfigPathEnvVar, filepath.Join(common.CfgPath, common.KubeCfgFile))
	// cluster.Status.Status = common.StatusRunning

	// // write current cluster to state file.
	// // native provider no need to operate .state file.
	// if p.Provider != "native" {
	// 	if err := common.DefaultDB.SaveCluster(cluster); err != nil {
	// 		return err
	// 	}
	// }

	// p.Logger.Infof("[%s] deploying additional manifests", p.Provider)

	// // deploy additional UI manifests.
	// enabledPlugins := map[string]bool{}
	// if cluster.UI {
	// 	enabledPlugins["dashboard"] = true
	// }

	// // deploy plugin
	// if cluster.Enable != nil {
	// 	for _, comp := range cluster.Enable {
	// 		enabledPlugins[comp] = true
	// 	}
	// }

	// for plugin := range enabledPlugins {
	// 	if plugin == "dashboard" {
	// 		if _, err := p.execute(&cluster.MasterNodes[0], []string{fmt.Sprintf(deployUICommand,
	// 			base64.StdEncoding.EncodeToString([]byte(dashboardTmpl)), common.K3sManifestsDir)}); err != nil {
	// 			return err
	// 		}
	// 	} else if plugin == "explorer" {
	// 		// start kube-explorer
	// 		port, err := common.EnableExplorer(context.Background(), cluster.ContextName)
	// 		if err != nil {
	// 			p.Logger.Errorf("[%s] failed to start kube-explorer for cluster %s: %v", p.Provider, cluster.ContextName, err)
	// 		}
	// 		if port != 0 {
	// 			p.Logger.Infof("[%s] kube-explorer for cluster %s will listen on 127.0.0.1:%d...", p.Provider, cluster.Name, port)
	// 		}
	// 	}
	// }

	// p.Logger.Infof("[%s] successfully deployed additional manifests", p.Provider)
	// p.Logger.Infof("[%s] successfully executed init k3s cluster logic", p.Provider)
	return nil
}

// GetAPIServerUrl url для подключения к API серверу
// сперва проверяется InternalDNS, InternalIP, ExternalDNS, ExternalIP
// если isExternal=true то сперва ExternalDNS, ExternalIP, InternalDNS, InternalIP
// если retry > 0 то проверяется tcp ping на API сервер (время между повторами с каждым радом больше в 2 раза)
func (p *ProviderBase) GetAPIServerUrl(master *k3sv1alpha1.Node, retry int, isExternal bool) (apiServerUrl string, err error) {
	if isExternal {
		for _, item := range p.Cluster.Spec.Networking.APIServerAddresses {
			if item.Type == v1alpha3.MachineAddressType(k3sv1alpha1.ExternalDNS) {
				p.Log.Debugf("[GetAPIServerUrl] check %v", item.Type)
				if retry == 0 {
					return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
				}
				if err := util.PingRetry(&util.PingArgs{
					Host: item.Address,
					Port: int(p.Cluster.Spec.Networking.APIServerPort),
					Retry: retry,
				}); err == nil {
					return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
				}
			} 
		}
		for _, item := range p.Cluster.Spec.Networking.APIServerAddresses {
		 if item.Type == v1alpha3.MachineAddressType(k3sv1alpha1.ExternalIP) {
				p.Log.Debugf("[GetAPIServerUrl] check %v", item.Type)
				if retry == 0 {
					return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
				}
				if err := util.PingRetry(&util.PingArgs{
					Host: item.Address,
					Port: int(p.Cluster.Spec.Networking.APIServerPort),
					Retry: retry,
				}); err == nil {
					return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
				}
			}
		}
	} 
	
	for _, item := range p.Cluster.Spec.Networking.APIServerAddresses {
		if item.Type == v1alpha3.MachineAddressType(k3sv1alpha1.InternalDNS) {
			p.Log.Debugf("[GetAPIServerUrl] check %v", item.Type)
			if retry == 0 {
				return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
			}
			if err := util.PingRetry(&util.PingArgs{
				Host: item.Address,
				Port: int(p.Cluster.Spec.Networking.APIServerPort),
				Retry: retry,
			}); err == nil {
				return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
			}
		}
	}
	for _, item := range p.Cluster.Spec.Networking.APIServerAddresses {
		if item.Type == v1alpha3.MachineAddressType(k3sv1alpha1.InternalIP) {
			p.Log.Debugf("[GetAPIServerUrl] check %v", item.Type)
			if retry == 0 {
				return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
			}
			if err := util.PingRetry(&util.PingArgs{
				Host: item.Address,
				Port: int(p.Cluster.Spec.Networking.APIServerPort),
				Retry: retry,
			}); err == nil {
				return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
			}
		}
	}

	if !isExternal {
		for _, item := range p.Cluster.Spec.Networking.APIServerAddresses {
			if item.Type == v1alpha3.MachineAddressType(k3sv1alpha1.ExternalDNS) {
				p.Log.Debugf("[GetAPIServerUrl] check %v", item.Type)
				if retry == 0 {
					return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
				}
				if err := util.PingRetry(&util.PingArgs{
					Host: item.Address,
					Port: int(p.Cluster.Spec.Networking.APIServerPort),
					Retry: retry,
				}); err == nil {
					return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
				}
			}
		}
		for _, item := range p.Cluster.Spec.Networking.APIServerAddresses {
			if item.Type == v1alpha3.MachineAddressType(k3sv1alpha1.ExternalIP) {
				p.Log.Debugf("[GetAPIServerUrl] check %v", item.Type)
				if retry == 0 {
					return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
				}
				if err := util.PingRetry(&util.PingArgs{
					Host: item.Address,
					Port: int(p.Cluster.Spec.Networking.APIServerPort),
					Retry: retry,
				}); err == nil {
					return fmt.Sprintf("https://%s:%d", item.Address, p.Cluster.Spec.Networking.APIServerPort), nil
				}
			}
		}
	}

  if isExternal {
		nodeIP, ok := p.Cluster.GetNodeAddress(master, "external")
		if ok {
			p.Log.Debugf("[GetAPIServerUrl] check %s node ip", "external")
			if retry == 0 {
				return fmt.Sprintf("https://%s:%d", nodeIP, p.Cluster.Spec.Networking.APIServerPort), nil
			}
			if err := util.PingRetry(&util.PingArgs{
				Host: nodeIP,
				Port: int(p.Cluster.Spec.Networking.APIServerPort),
				Retry: retry,
			}); err == nil {
			return fmt.Sprintf("https://%s:%d", nodeIP, p.Cluster.Spec.Networking.APIServerPort), nil
			}
		}
	}
	
	nodeIP, ok := p.Cluster.GetNodeAddress(master, "internal")
	if ok {
		p.Log.Debugf("[GetAPIServerUrl] check %s node ip", "internal")
		if retry == 0 {
			return fmt.Sprintf("https://%s:%d", nodeIP, p.Cluster.Spec.Networking.APIServerPort), nil
		}
		if err := util.PingRetry(&util.PingArgs{
			Host: nodeIP,
			Port: int(p.Cluster.Spec.Networking.APIServerPort),
			Retry: retry,
		}); err == nil {
			return fmt.Sprintf("https://%s:%d", nodeIP, p.Cluster.Spec.Networking.APIServerPort), nil
		}
	}

	return "", fmt.Errorf("Is NOT set Api server URL")
}

// SetClientset setting Clientset for clusterName
func (p *ProviderBase) SetClientset(clusterName string) (err error) {
	
	// use the current context in kubeconfig
	config, err := k3s.BuildKubeConfigFromFlags(clusterName)
	if err != nil {
			return err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err !=nil {
		return err
	}
	p.Clientset=clientset

	return err
}

// SetClientsetFromConfig create Clientset from clientcmdapi.Config
func (p *ProviderBase) SetClientsetFromConfig() (err error) {
	
	// // use the current context in kubeconfig
	// config, err := k3s.BuildKubeConfigFromFlags(clusterName)
	// config, err := clientcmd.DefaultClientConfig.ClientConfig()
	// if err != nil {
	// 		return err
	// }

	// clientConfig, err := clientcmd.NewDefaultClientConfig(*kubeconfig, &clientcmd.ConfigOverrides{
	// 	ClusterDefaults: clientcmdapi.Cluster{Server: master},
	// }).ClientConfig()
	clientConfig, err := clientcmd.NewDefaultClientConfig(*p.Config, &clientcmd.ConfigOverrides{}).ClientConfig()

	clientset, err := kubernetes.NewForConfig(clientConfig)
	if err !=nil {
		return err
	}
	p.Clientset=clientset

	return err
}

// GetClusterStatus cluster status
func (p *ProviderBase) GetClusterStatus() string {
	// _, err := c.Clientset.RESTClient().Get().Timeout(15 * time.Second).RequestURI("/readyz").DoRaw(context.TODO())
	_, err := p.Clientset.RESTClient().Get().Timeout(types.NodeWaitForLogMessageRestartWarnTime).RequestURI("/readyz").DoRaw(context.TODO())
	if err != nil {
		return types.ClusterStatusStopped
	}
	return types.ClusterStatusRunning
}

// ListNodes Read more on Kubernets Nodes. In client-go, NodeInterface includes all the APIs to deal with Nodes.
func (p *ProviderBase) ListNodes() ([]v1.Node, error) {
	nodes, err := p.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
			return nil, err
	}

	return nodes.Items, nil
}

// // DescribeClusterNodes describe cluster nodes.
// func (p *ProviderBase) DescribeClusterNodes() (instanceNodes []k3sv1alpha1.ClusterNode, err error) {
// 	// list cluster nodes.
// 	timeout := int64(5 * time.Second)
// 	nodeList, err := p.Clientset.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{TimeoutSeconds: &timeout})
// 	if err != nil || nodeList == nil {
// 		return nil, err
// 	}
// 	for _, node := range nodeList.Items {
// 		var internalIP, hostName string
// 		addressList := node.Status.Addresses
// 		for _, address := range addressList {
// 			switch address.Type {
// 			case v1.NodeInternalIP:
// 				internalIP = address.Address
// 			case v1.NodeHostName:
// 				hostName = address.Address
// 			default:
// 				continue
// 			}
// 		}
// 		for index, n := range instanceNodes {
// 			isCurrentInstance := false
// 			for _, address := range n.InternalIP {
// 				if address == internalIP {
// 					isCurrentInstance = true
// 					break
// 				}
// 			}
// 			if !isCurrentInstance {
// 				if n.InstanceID == node.Name {
// 					isCurrentInstance = true
// 				}
// 			}
// 			if isCurrentInstance {
// 				n.HostName = hostName
// 				n.Version = node.Status.NodeInfo.KubeletVersion
// 				n.ContainerRuntimeVersion = node.Status.NodeInfo.ContainerRuntimeVersion
// 				// get roles.
// 				labels := node.Labels
// 				roles := make([]string, 0)
// 				for role := range labels {
// 					if strings.HasPrefix(role, "node-role.kubernetes.io") {
// 						roleArray := strings.Split(role, "/")
// 						if len(roleArray) > 1 {
// 							roles = append(roles, roleArray[1])
// 						}
// 					}
// 				}
// 				if len(roles) == 0 {
// 					roles = append(roles, "<none>")
// 				}
// 				sort.Strings(roles)
// 				n.Roles = strings.Join(roles, ",")
// 				// get status.
// 				// conditions := node.Status.Conditions
// 				// for _, c := range conditions {
// 				// 	if c.Type == v1.NodeReady {
// 				// 		if c.Status == v1.ConditionTrue {
// 				// 			n.Status = "Ready"
// 				// 		} else {
// 				// 			n.Status = "NotReady"
// 				// 		}
// 				// 		break
// 				// 	}
// 				// }
// 				n.Status = k3sClient.GetStatus(&node)
// 				instanceNodes[index] = n
// 				break
// 			}
// 		}
// 	}
// 	return instanceNodes, nil
// }

// LoadNodeStatus 
func (p *ProviderBase) LoadNodeStatus() {
	masters := p.Cluster.Status.MasterNodes
	for i, node := range masters {
		if len(node.State.Status) == 0 {
			masters[i].State.Status = types.ClusterStatusUnknown
			// проверяем есть ли у нас подключение к API кубера
			if p.Clientset == nil {
				// masters[i].State.Status = types.StatusMissing
				// проверяем что установили k3s server
				if ok := p.CheckExitFile(types.MasterUninstallCommand, node); !ok {
					masters[i].State.Status = types.StatusMissing
				} else {
					p.Log.Warningln("[LoadNodeStatus] TODO: [KCTL-11] добавить проверку состояния кластера")
					kubeconfig, err := p.GetKubeconfig(node)
					if err != nil {
						p.Log.Errorf("[LoadNodeStatus] GetKubeconfig: %v", err.Error())
					}
					p.Config = kubeconfig
					p.SetClientsetFromConfig()

					clusterStatus := p.GetClusterStatus()
					masters[i].State.Status = types.StatusMissing
					p.Log.Infof("[LoadNodeStatus] Cluster STATUS: %s", clusterStatus)
					// isExternal := true
					// apiServerUrl, err := p.Cluster.GetAPIServerUrl(masters, &p.Cluster.Spec.Networking, isExternal)
					// if err != nil {
					// 	log.Fatal(err)
					// }
					// log.Debugf("apiServerUrl: %s", apiServerUrl)
					// v, _ := yaml.Marshal(kubeconfig)
					// p.Log.Debugf("========== kubeconfig ==========\n%s\n==========================\n", v)
				}
			}
		} else if p.Clientset == nil {
			kubeconfig, err := p.GetKubeconfig(node)
			if err != nil {
				p.Log.Errorf("[LoadNodeStatus] GetKubeconfig: %v", err.Error())
			}
			p.Config = kubeconfig
			p.SetClientsetFromConfig()
		}
	}
}

// GetMasterNodes return master nodes
func (p *ProviderBase) GetMasterNodes() (masters []*k3sv1alpha1.Node) {
	masters = p.Cluster.Status.MasterNodes
	for _, node := range masters{
		p.Log.Warnf("master status: %v", node.State.Status)
	}
	return masters
}

// GetMasterNodes return workers nodes
func (p *ProviderBase) GetWorkerNodes() (workers []*k3sv1alpha1.Node) {
	return p.Cluster.Status.WorkerNodes
}

// Execute command local or ssh 
func (p *ProviderBase) Execute(command string, node *k3sv1alpha1.Node, stream bool) (stdOut string, err error) {
	bastion, err := p.Cluster.GetBastion(node.Bastion, node)
	if err != nil {
		p.Log.Fatalln(err.Error())
	} else {
		p.Log.Debugf("master node: %s bastion: %s", node.Name, bastion.Address)
	}
	
	res, err := yaml.Marshal(node)
	if err != nil {
		p.Log.Errorf(err.Error())
	}
	p.Log.Tracef("--------------------------------------\nNODE\n-------------------------------------- \n%s\n--------------------------------------", res)

	if node.Bastion == "local" {
		p.Log.Infoln("Run command in localhost........")
		// stdOut, stdErr, err := RunLocalCommand(installK3scommand, true, dryRun)
		stdOut, stdErr, err := k3s.RunLocalCommand(command, true, p.CmdFlags.DryRun)
		if err != nil {
			p.Log.Fatalln(err.Error())
		} else if len(stdErr) > 0 {
			p.Log.Errorf("stderr: %q", stdErr)
		}
		p.Log.Debugf("stdout: %q", stdOut)
		
		return string(stdOut), err
	} else {
		if node.User != "root" {
			command = fmt.Sprintf("sudo %s", command)
		}
		p.NewSSH(bastion)
		if stream {
			p.sshStream(command, false)
			return stdOut, err
		}
		stdOut, stdErr, err := p.sshExecute(command)
		if err != nil {
			p.Log.Errorf("--- sshExecute stdErr ---\n%v\n--------------------------\n%v\n--- END stdErr ---", stdErr, err.Error())
			// log.Fatalln(err.Error())
		}
		return stdOut, err
	}
	// RunExampleCommand2()

}

// ExecuteMaster execute command in master node
func (p *ProviderBase) ExecuteMaster(command string, first bool) (stdOuts []string, err error) {
	for _, node := range p.GetMasterNodes() {
		stdOut, err := p.Execute(command, node, false)
		if err == nil {
			stdOuts = append(stdOuts, stdOut)
			if first {
				return stdOuts, nil
			}
		}
	}
	
	if len(stdOuts) > 0 {
		return stdOuts, err
	}
	return stdOuts, fmt.Errorf("Is Not master node to run commmands")
}

// ShutDownWithDrain will cause q to ignore all new items added to it. As soon
// as the worker goroutines have "drained", i.e: finished processing and called
// Done on all existing items in the queue; they will be instructed to exit and
// ShutDownWithDrain will return. Hence: a strict requirement for using this is;
// your workers must ensure that Done is called on all items in the queue once
// the shut down has been initiated, if that is not the case: this will block
// indefinitely. It is, however, safe to call ShutDown after having called
// ShutDownWithDrain, as to force the queue shut down to terminate immediately
// without waiting for the drainage.
func (p *ProviderBase) ShutDownWithDrain(node *k3sv1alpha1.Node) {
	p.setDrain(node)
	p.shutdown(node)
	// q.setDrain(true)
	// q.shutdown()
	// for q.isProcessing() && q.shouldDrain() {
	// 	q.waitForProcessing()
	// }
}

// setDrain execute drain command in master node
func (p *ProviderBase) setDrain(node *k3sv1alpha1.Node) {
	command := fmt.Sprintf(types.DrainCommand, node.Name)
	for _, master := range p.GetMasterNodes() {
		stdOut, err := p.Execute(command, master, false)
		if err != nil {
			p.Log.Errorf(err.Error())
		} else {
			p.Log.Debugf("[setDrain] stdOut: %v", stdOut)
			break
		}
	}
}

// setDelete execute delete node command in master node
func (p *ProviderBase) setDelete(node *k3sv1alpha1.Node) {
	command := fmt.Sprintf(types.DeleteNodeCommand, node.Name)
	for _, master := range p.GetMasterNodes() {
		stdOut, err := p.Execute(command, master, false)
		if err != nil {
			p.Log.Errorf(err.Error())
		} else {
			p.Log.Debugf("[setDelete] stdOut: %v", stdOut)
			break
		}
	}
}

// shutdown uninstall k3s TODO: shutdown node command
func (p *ProviderBase) shutdown(node *k3sv1alpha1.Node) {
	command := fmt.Sprintf("sh %s", types.WorkerUninstallCommand)
	if node.Role == k3sv1alpha1.Role(types.ServerRole) {
		command = fmt.Sprintf("sh %s", types.MasterUninstallCommand)
	}
	_, _ = p.Execute(command, node, true)

	// stdOut, err := p.Execute(command, node, true)
	// if err != nil {
	// 	p.Log.Errorf(err.Error())
	// }
	// p.Log.Debugf("[shutdown] stdOut: %v", stdOut)
}

// MakeInstallExec установка сервера
func (p *ProviderBase) MakeInstallExec() (k3sIstallOptions k3sv1alpha1.K3sIstallOptions) {
	extraArgs := []string{}
	k3sIstallOptions = k3sv1alpha1.K3sIstallOptions{
		K3sVersion: p.Cluster.Spec.KubernetesVersion,
		K3sChannel: p.Cluster.Spec.K3sChannel,
	}

	if len(p.Cluster.Spec.Datastore.Provider) > 0 {
		if p.Cluster.Spec.Datastore.Provider == k3sv1alpha1.DatastoreEtcd {
			k3sIstallOptions.IsCluster = true
		} else if datastore, err := p.Cluster.GetDatastore(); err != nil {
			p.Log.Fatalln(err.Error())
		} else {
			extraArgs = append(extraArgs, fmt.Sprintf("--datastore-endpoint %s", datastore))
			p.Log.Infof("datastore connection string: %s", datastore)
		}
	}

	if p.Cluster.Spec.Options.DisableLoadbalancer {
		extraArgs = append(extraArgs, "--no-deploy servicelb")
	} else {
		if len(p.Cluster.Spec.LoadBalancer.MetalLb) > 0 {
			// TODO: #3 добавить проверку на ip adress
			p.Log.Debugln("LoadBalancer MetalLB: ", p.Cluster.Spec.LoadBalancer.MetalLb)
			extraArgs = append(extraArgs, "--no-deploy servicelb")
			k3sIstallOptions.LoadBalancer = types.MetalLb
		} else if len(p.Cluster.Spec.LoadBalancer.KubeVip) > 0 {
			// TODO: добавить проверку на ip adress
			p.Log.Debugln("LoadBalancer kube-vip: ", p.Cluster.Spec.LoadBalancer.KubeVip)
			extraArgs = append(extraArgs, "--no-deploy servicelb")
			k3sIstallOptions.LoadBalancer = types.KubeVip
		}
	}

	// if options.Options.DisableIngress || len(options.Ingress) > 0 {
	// 	if ingress, isset := util.Find(types.IngressControllers, options.Ingress); isset {
	// 		k3sIstallOptions.Ingress = ingress
	// 		extraArgs = append(extraArgs, "--no-deploy traefik")
	// 	} else {
	// 		p.Log.Fatalf("Ingress Controllers %s not support :(", options.Ingress)
	// 	}
	// }
	extraArgs = append(extraArgs, "--no-deploy traefik")

	if len(p.Cluster.Spec.Networking.ServiceSubnet) > 0 {
		p.Log.Debugln("ServiceSubnet: ", p.Cluster.Spec.Networking.ServiceSubnet)
		extraArgs = append(extraArgs, fmt.Sprintf("--service-cidr %s", p.Cluster.Spec.Networking.ServiceSubnet))
	}

	if len(p.Cluster.Spec.Networking.PodSubnet) > 0 {
		p.Log.Debugln("PodSubnet: ", p.Cluster.Spec.Networking.PodSubnet)
		extraArgs = append(extraArgs, fmt.Sprintf("--cluster-cidr %s", p.Cluster.Spec.Networking.PodSubnet))
	}

	if len(p.Cluster.Spec.Networking.DNSDomain) > 0 {
		p.Log.Debugln("DNSDomain: ", p.Cluster.Spec.Networking.DNSDomain)
		extraArgs = append(extraArgs, fmt.Sprintf("--cluster-domain %s", p.Cluster.Spec.Networking.DNSDomain))
	}

	if len(p.Cluster.Spec.Networking.ClusterDns) > 0 {
		p.Log.Debugln("ClusterDns: ", p.Cluster.Spec.Networking.ClusterDns)
		extraArgs = append(extraArgs, fmt.Sprintf("--cluster-dns %s", p.Cluster.Spec.Networking.ClusterDns))
	}

	k3sIstallOptions.Backend = types.Vxlan
	k3sIstallOptions.CNI = types.Flannel
	if len(p.Cluster.Spec.Networking.CNI) > 0 {
		if cni, isset := util.Find(types.CNIplugins, p.Cluster.Spec.Networking.CNI); isset {
			k3sIstallOptions.CNI = cni
		} else {
			p.Log.Fatalf("CNI plugins %s not support :(", p.Cluster.Spec.Networking.CNI)
		}
	}
	if len(p.Cluster.Spec.Networking.Backend) > 0 {
		if k3sIstallOptions.CNI == types.Flannel {
			if backend, isset := util.Find(types.FlannelBackends, p.Cluster.Spec.Networking.Backend); isset {
				k3sIstallOptions.Backend = backend
			} else {
				p.Log.Fatalf("CNI plugins %s backend %s not support :(", p.Cluster.Spec.Networking.CNI, p.Cluster.Spec.Networking.Backend)
			}
		} else if k3sIstallOptions.CNI == types.Calico {
			if backend, isset := util.Find(types.CalicoBackends, p.Cluster.Spec.Networking.Backend); isset {
				k3sIstallOptions.Backend = backend
			} else {
				p.Log.Fatalf("CNI plugins %s backend %s not support :(", p.Cluster.Spec.Networking.CNI, p.Cluster.Spec.Networking.Backend)
			}
		} else if k3sIstallOptions.CNI == types.Cilium {
			if backend, isset := util.Find(types.CiliumBackends, p.Cluster.Spec.Networking.Backend); isset {
				k3sIstallOptions.Backend = backend
			} else {
				p.Log.Fatalf("CNI plugins %s backend %s not support :(", p.Cluster.Spec.Networking.CNI, p.Cluster.Spec.Networking.Backend)
			}
		}
	}
	if k3sIstallOptions.CNI == types.Flannel {
		extraArgs = append(extraArgs, fmt.Sprintf("--flannel-backend=%s", k3sIstallOptions.Backend))
	} else {
		extraArgs = append(extraArgs, "--flannel-backend=none")
	}

	if p.Cluster.Spec.Options.SecretsEncryption {
		extraArgs = append(extraArgs, "--secrets-encryption")
	}

	if p.Cluster.Spec.Options.SELinux {
		extraArgs = append(extraArgs, "--selinux")
	}

	if p.Cluster.Spec.Options.Rootless {
		extraArgs = append(extraArgs, "--rootless")
	}

	extraArgsCmdline := ""
	for _, a := range extraArgs {
		extraArgsCmdline += a + " "
	}

	for _, a := range p.Cluster.Spec.K3sOptions.ExtraServerArgs {
		if a != "[]" {
			extraArgsCmdline += a + " "
		}
	}

	installExec := ""

	if trimmed := strings.TrimSpace(extraArgsCmdline); len(trimmed) > 0 {
		installExec += fmt.Sprintf(" %s", trimmed)
	}

	if len(k3sIstallOptions.LoadBalancer) == 0 {
		k3sIstallOptions.LoadBalancer = types.ServiceLb
	}

	k3sIstallOptions.ExecString = installExec

	// --tls-san developer.cluster --node-taint CriticalAddonsOnly=true:NoExecute
	return k3sIstallOptions
}

// MakeAgentInstallExec compile agent install string
func (p *ProviderBase) MakeAgentInstallExec(opts *k3sv1alpha1.K3sWorkerOptions) string {
	// curl -sfL https://get.k3s.io | K3S_URL='https://<IP>6443' K3S_TOKEN='<TOKEN>' INSTALL_K3S_CHANNEL='stable' sh -s - --node-taint key=value:NoExecute
	// p.Log.Debugf("K3sVersion=%v K3sChannel=%v, %v", opts.K3sVersion, opts.K3sChannel, util.CreateVersionStr(opts.K3sVersion, opts.K3sChannel))
	return fmt.Sprintf(opts.JoinAgentCommand, opts.ApiServerAddres, opts.ApiServerPort, opts.Token, util.CreateVersionStr(opts.K3sVersion, opts.K3sChannel))
}

// initAdditionalMaster add first master node
func (p *ProviderBase) initAdditionalMaster(tlsSAN []string, node *k3sv1alpha1.Node, opts *k3sv1alpha1.K3sIstallOptions) {
	// TODO: перевести на K3S_AGENT_TOKEN_FILE
	extraArgs := fmt.Sprintf("K3S_AGENT_TOKEN='%s'", p.Cluster.Spec.AgentToken)
	execArgs := ""
	
	// TODO: перевести на переменнын окружения
	// K3S_CLUSTER_INIT
	// K3S_CLUSTER_RESET
	if opts.IsCluster {
		// extraArgs = fmt.Sprintf("%s K3S_CLUSTER_INIT=true", extraArgs)
		execArgs += " --cluster-init"
	}
	
	if len(tlsSAN) > 0 {
		for _, san := range tlsSAN {
			execArgs += fmt.Sprintf(" --tls-san %s", san)
		}
	}
	
	for _, ip := range node.Addresses {
		if ip.Type == v1alpha3.MachineAddressType(k3sv1alpha1.ExternalIP) {
			execArgs = fmt.Sprintf(" %s --node-external-ip %s", execArgs, ip.Address)
		} else if ip.Type == v1alpha3.MachineAddressType(k3sv1alpha1.InternalIP) {
			execArgs = fmt.Sprintf(" %s --node-ip %s", execArgs, ip.Address)
		} 
	}
	command := fmt.Sprintf(types.InitMasterCommand, types.K3sGetScript, extraArgs, p.Cluster.Spec.ClusterToken, opts.ExecString, execArgs, util.CreateVersionStr(opts.K3sVersion, opts.K3sChannel))
	p.Log.Debugf("[initAdditionalMaster] RUN %s", command)

	err := retry.Do(
		func() error {
			result, err := p.Execute(command, node, true)

			if err == nil {
			// 	defer func() {
			// 		if err := resp.Body.Close(); err != nil {
			// 			panic(err)
			// 		}
				if len(result) > 0 {
					p.Log.Debugf("--- |%s| ---", strings.Trim(result, "\n"))
					return nil
				}
			// 	}()
			// 	body, err = ioutil.ReadAll(resp.Body)
			}

			return err
		},
		retry.Attempts(10), // количество попыток
		retry.Delay(10 * time.Second), // задержка в секундах
	)
	if err != nil {
		p.Log.Errorf(err.Error())	
	}

}

// joinWorker join worker node to cluster
func (p *ProviderBase) joinWorker(token string, node *k3sv1alpha1.Node) {
	// command := types.WorkerUninstallCommand
	apiServerAddres, err := p.Cluster.GetAPIServerAddress(node, &p.Cluster.Spec.Networking)
	p.Log.Debugf("apiServerAddresses: %s", apiServerAddres)
	if err != nil {
		p.Log.Fatal(err)
	}

	// TODO: add lavels to node
	cnt := p.Cluster.GetNodeLabels(node)
	p.Log.Warnf("TODO: add lavels to node =-> cnt: %d", cnt)

	opts := &k3sv1alpha1.K3sWorkerOptions{
		JoinAgentCommand: types.JoinAgentCommand,
		ApiServerAddres: apiServerAddres,
		ApiServerPort: p.Cluster.Spec.Networking.APIServerPort,
		Token: token,
		K3sVersion: p.Cluster.Spec.KubernetesVersion,
		K3sChannel: p.Cluster.Spec.K3sChannel,
	}
	command := p.MakeAgentInstallExec(opts)
	// p.Log.Debugf("Exec command: %s", command)
	// installk3sAgentExec := p.Cluster.MakeAgentInstallExec(opts)
			// 			installk3sAgentExec.K3sChannel = cfg.Spec.K3sChannel
			// 			installk3sAgentExec.K3sVersion = cfg.Spec.KubernetesVersion
			// 			installk3sAgentExec.Node = node

	// _, _ = p.Execute(command, node, true)
	// stdOut, err := p.Execute(command, node, true)
	// if err != nil {
	// 	p.Log.Errorf(err.Error())
	// }
	_, _ = p.Execute(command, node, true)
	// p.Log.Debugf("[joinWorker] stdOut: %v", stdOut)
}

// GetAgentToken возвращает токен агента
func (p *ProviderBase) GetAgentToken(master *k3sv1alpha1.Node) (token string, err error) {
	command := fmt.Sprintf("cat %s", types.FileClusterToken)
	token, err = p.Execute(command, master, false)
	return strings.Trim(token, "\n"), err
}

// CheckExitFile проверка на существование файла на сервере
func (p *ProviderBase) CheckExitFile(file string, node *k3sv1alpha1.Node) (ok bool) {
	command := fmt.Sprintf(types.TestExitFile, file)
	err := retry.Do(
		func() error {
			result, err := p.Execute(command, node, false)

			if err == nil {
			// 	defer func() {
			// 		if err := resp.Body.Close(); err != nil {
			// 			panic(err)
			// 		}
				if len(result) > 0 {
					p.Log.Debugf("--- file: %s |%s| ---", file, strings.Trim(result, "\n"))
					ok = false
					return nil
				}
				ok = true
			// 	}()
			// 	body, err = ioutil.ReadAll(resp.Body)
			}

			return err
		},
		retry.Attempts(10), // количество попыток
		retry.Delay(10 * time.Second), // задержка в секундах
	)
	if err != nil {
		p.Log.Errorf("---------- CheckExitFile ------------")
		p.Log.Errorf(err.Error())	
	}
	
	return ok
}

// CreateK3sCluster create K3S cluster.
func (p *ProviderBase) CreateK3sCluster() (err error) {
	// logFile, err := common.GetLogFile(p.Name)
	// if err != nil {
	// 	return err
	// }
	// c := &types.Cluster{
	// 	Metadata: p.Metadata,
	// 	Options:  p.Options,
	// 	Status:   p.Status,
	// }
	defer func() {
		if err != nil {
			p.Log.Errorf("[%s] failed to create cluster: %v", p.Cluster.GetName(), err)
			// TODO: сделать откат кластера при ошибке
			// p.RollbackCluster(p.rollbackInstance)
		}
		// if err == nil && len(p.Cluster.Status.MasterNodes) > 0 {
		// 	p.Log.Info(types.UsageInfoTitle)
		// 	p.Log.Infof(types.UsageContext, p.Cluster.GetName())
		// 	p.Log.Info(types.UsagePods)
		// }
		// _ = logFile.Close()
		if p.Callbacks != nil {
			if process, ok := p.Callbacks[p.Cluster.GetName()]; ok && process.Event == "create" {
				logEvent := &types.LogEvent{
					Name:        "create",
					ContextName: p.Cluster.GetName(),
				}
				process.Fn(logEvent)
			}
		}
	}()

	// p.Logger = common.NewLogger(common.Debug, logFile)
	// p.Logger.Infof("[%s] executing create logic...", p.GetProviderName())

	// // set ssh default value.
	// if p.SSHUser == "" {
	// 	p.SSHUser = defaultUser
	// }
	// if p.SSHPassword == "" && p.SSHKeyPath == "" {
	// 	p.SSHKeyPath = defaultSSHKeyPath
	// }
	// // assemble node status.
	// if c, err = p.assembleNodeStatus(&p.SSH); err != nil {
	// 	return err
	// }
	// c.SSH = p.SSH
	// // initialize K3s cluster.
	if err = p.InitK3sCluster(); err != nil {
		return err
	}

	if p.Clientset != nil {
		p.Log.Infoln("Save kubeconfig to file...")
		opts := k3s.WriteKubeConfigOptions{
			OverwriteExisting: true,
			UpdateCurrentContext: p.Cluster.Spec.KubeconfigOptions.SwitchCurrentContext,
		}
		_, err := k3s.SaveKubeconfig(p.Config, opts)
		if err != nil {
			p.Log.Errorf(err.Error())
		}
	}
	// // deploy custom manifests.
	// if p.Manifests != "" {
	// 	deployCmd, err := p.GetCustomManifests()
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if err = p.DeployExtraManifest(c, deployCmd); err != nil {
	// 		return err
	// 	}
	// 	p.Logger.Infof("[%s] successfully deployed custom manifests", p.Provider)
	// }
	// p.Logger.Infof("[%s] successfully executed create logic", p.GetProviderName())
	return nil
}

// DeleteNode delete node from cluster
func (p *ProviderBase) AddNode(nodeName string) (ok bool) {
	var err error
	var token string

	for _, master := range p.GetMasterNodes() {
		p.Log.Warnf("master node: %s", master.Name)
		token, err = p.GetAgentToken(master)
		if err != nil {
			p.Log.Errorf(err.Error())
		} else {
			// p.Log.Debugf("[AddNode] K3S_TOKEN=%s", token)
			break
		}
	}

	for _, node := range p.GetWorkerNodes() {
		if node.Name == nodeName {
			if node.Role == k3sv1alpha1.Role(types.ServerRole) {
				p.Log.Infof("[AddNode] TODO: Add Master node: %s", node.Name)
			} else{
				p.Log.Infof("[AddNode] Add Worker node: %s", node.Name)
				p.joinWorker(token, node)
				ok = true
			}
		}
	}

	return ok
}

// DeleteNode delete node from cluster
func (p *ProviderBase) DeleteNode(nodeName string, allNodes bool) (cnt int) {
	// p.Log.Infof("------> Delete Worker node: %s", nodeName)
	cnt = 0
	workers := p.GetWorkerNodes()
	// cntWorker := len(workers)
	for _, node := range workers {
		if node.Name == nodeName || allNodes {
			p.Log.Infof("Delete Worker node: %s", node.Name)
			p.setDrain(node)
			p.setDelete(node)
			p.shutdown(node)

			// if bastion, err := p.Cluster.GetBastion(node.Bastion, node); err != nil {
			// 	p.Log.Fatalln(err.Error())
			// } else {
			// 	command := "uname -a"
					
			// 	stdOuts, err := p.ExecuteMaster(command, false)
			// 	if err != nil {
			// 		p.Log.Fatalf("[ExecuteMaster] %v", err.Error())
			// 	}
			// 	for _, stdOut := range stdOuts {
			// 		p.Log.Debugf("[ExecuteMaster] stdOut: %v", stdOut)
			// 	}

			// 	p.Log.Debugf("bastion: %s", bastion.Address)
				
			// 	err = p.SetClientset(p.Cluster.ObjectMeta.Name)
			// 	if err !=nil {
			// 		p.Log.Errorf(err.Error())
			// 	}
			// 	clusterStatus := p.GetClusterStatus()
			// 	p.Log.Infof("clusterStatus: %v", clusterStatus)
			// 	// p.Log.Debugf("------------------------------\n%v\n------------------------------", clusterStatus)
			// 	// n, err := client.ListNodes()
			// 	// if err != nil {
			// 	// 	p.Log.Errorf(err.Error())
			// 	// }
			// 	// p.Log.Warnf("list nodes: %v", n[0].GetName())
				
			// 	p.Log.Infof("Successfully Worker node: %s", node.Name)
				cnt += 1
			// }
		}
	}

	masters := p.GetMasterNodes()
	cntMaster := len(masters)
	for _, node := range masters {
		if node.Name == nodeName || allNodes {
			p.Log.Infof("Delete Master node: %s", node.Name)
			if cntMaster > 1 {
				p.setDrain(node)
			}
			p.shutdown(node)
			
			p.Log.Warnln("TODO: delete kubeconfig context")
			// err := p.SetClientset(p.Cluster.ObjectMeta.Name)
			// if err !=nil {
			// 	p.Log.Errorf(err.Error())
			// }
			// clusterStatus := p.GetClusterStatus()
			// p.Log.Infof("clusterStatus: %v", clusterStatus)
			// p.Log.Infof("Successfully delete Master node: %s", node.Name)
			
			cnt += 1
			cntMaster -= 1
		}
	}

	return cnt
}

// setGroupNodes nodes grouping to role
func (p *ProviderBase) setGroupNodes() {
	// serverNodes := ServerNodes{}
	// agentNodes := AgentNodes{}
	for _, node := range p.Cluster.Spec.Nodes {
		// log.Tracef("Node (%+v): Checking node role %s", node, node.Role)
		role := util.GetNodeRole(string(node.Role))
		node.Role = k3sv1alpha1.Role(role)
		if role == "server" {
			p.Cluster.Status.MasterNodes = append(p.Cluster.Status.MasterNodes, node)
		} else {
			p.Cluster.Status.WorkerNodes = append(p.Cluster.Status.WorkerNodes, node)
		}
	}
	p.Cluster.Spec.Servers = len(p.Cluster.Status.MasterNodes)
	p.Cluster.Spec.Agents = len(p.Cluster.Status.WorkerNodes)
	// if len(serverNodes) == 0 {
	// 	log.Fatalln("Is not set server node :(")
	// }
	// return serverNodes, agentNodes, nil
}

func (p *ProviderBase) GetKubeconfig(master *k3sv1alpha1.Node) (*clientcmdapi.Config, error) {
	// command := fmt.Sprintf("cat %s", types.FileClusterToken)
	var kubeconfig string
	var err error
	err = retry.Do(
		func() error {
			kubeconfig, err = p.Execute(types.CatCfgCommand, master, false)
			if err == nil {
				if len(kubeconfig) > 0 {
					// p.Log.Debugf("--- |%s| ---", strings.Trim(kubeconfig, "\n"))	
					return nil
				}
			}
			return err
		},
		retry.Attempts(10), // количество попыток
		retry.Delay(10 * time.Second), // задержка в секундах
	)
	if err != nil {
		// p.Log.Errorf(err.Error())
		return nil, err
	}
	isExternal := false
	if p.Cluster.Spec.KubeconfigOptions.ConnectType == "ExternalIP" {
		isExternal = true
	}
	apiServerUrl, err := p.GetAPIServerUrl(master, 1, isExternal)
	if err != nil {
		p.Log.Fatalf("[GetKubeconfig] %v", err.Error())
	}
	// timeoutSec := 1
	// counterPing := 2
	// if err := util.PingTCP(apiServerUrl, timeoutSec, counterPing); err != nil {
	// 	p.Log.Errorf(err.Error())
	// 	p.Log.Fatalf("[GetKubeconfig] IS NOT Connect to master node :(")
	// }

	opts := k3s.WriteKubeConfigOptions{
		OverwriteExisting: true,
		UpdateCurrentContext: p.Cluster.Spec.KubeconfigOptions.SwitchCurrentContext,
	}
	// pathKubeConfig, err := k3s.SaveCfg(kubeconfig, apiServerUrl, p.Cluster.GetObjectMeta().GetName(), opts)
	// newKubeconfig, err := k3s.LoadKubeconfig(kubeconfig, apiServerUrl, p.Cluster.GetObjectMeta().GetName(), opts)
	// if err !=nil {
	// 	p.Log.Errorln(err.Error())
	// }
	// k, e := clientcmd.NewClientConfigFromBytes([]byte(newKubeConfig))
	// c, _ := yaml.Marshal(newKubeConfig.Clusters)
	// p.Log.Infof("new kubeconfig: %v", c)
	// p.Log.Infof("new kubeconfig: %v", newKubeConfig.Clusters)

	// p.Log.Warnf("[GetKubeconfig] apiServerUrl: %s", apiServerUrl)

	// for _, item := range masters {
	// 	kubeconfig, err = ExecuteMaster(types.CatCfgCommand, &item, dryRun)
	// 	// if err != nil {
	// 	// 	log.Errorln(err)
	// 	// }
	// 	return kubeconfig, err
	// }
	return k3s.LoadKubeconfig(kubeconfig, apiServerUrl, p.Cluster.GetObjectMeta().GetName(), opts)
}

// initLogging set loging
func (p *ProviderBase) initLogging(cmdFlags *types.CmdFlags) {
	p.Log = logrus.New()
	if cmdFlags.TraceLogging {
		p.Log.SetLevel(logrus.TraceLevel)
	} else if cmdFlags.DebugLogging {
		p.Log.SetLevel(logrus.DebugLevel)
	} else {
		switch cmdFlags.LogLevel {
		case "TRACE":
			p.Log.SetLevel(logrus.TraceLevel)
		case "DEBUG":
			p.Log.SetLevel(logrus.DebugLevel)
		case "WARN":
			p.Log.SetLevel(logrus.WarnLevel)
		case "ERROR":
			p.Log.SetLevel(logrus.ErrorLevel)
		default:
			p.Log.SetLevel(logrus.InfoLevel)
		}
	}

	// log.SetOutput(ioutil.Discard)
	// log.AddHook(&writer.Hook{
	// 	Writer: os.Stderr,
	// 	LogLevels: []log.Level{
	// 		log.PanicLevel,
	// 		log.FatalLevel,
	// 		log.ErrorLevel,
	// 		log.WarnLevel,
	// 	},
	// })
	// log.AddHook(&writer.Hook{
	// 	Writer: os.Stdout,
	// 	LogLevels: []log.Level{
	// 		log.InfoLevel,
	// 		log.DebugLevel,
	// 		log.TraceLevel,
	// 	},
	// })

	// мне это ненужно сейчас
	// formatter := &logrus.TextFormatter{
	// 	ForceColors: true,
	// }
	// // if flags.timestampedLogging || os.Getenv("LOG_TIMESTAMPS") != "" {
	// // 	formatter.FullTimestamp = true
	// // }
	// p.Log.SetFormatter(formatter)

}

func (p *ProviderBase) NewSSH(bastion *k3sv1alpha1.BastionNode) {
	p.SSH = &easyssh.MakeConfig{
		User: bastion.User,
		Port:    fmt.Sprintf("%d", bastion.SshPort),
		Timeout: 60 * time.Second,

		// Parse PrivateKey With Passphrase
		// Passphrase: "XXXX",

		// Optional fingerprint SHA256 verification
		// Get Fingerprint: ssh.FingerprintSHA256(key)
		// Fingerprint: "SHA256:................E"

		// Enable the use of insecure ciphers and key exchange methods.
		// This enables the use of the the following insecure ciphers and key exchange methods:
		// - aes128-cbc
		// - aes192-cbc
		// - aes256-cbc
		// - 3des-cbc
		// - diffie-hellman-group-exchange-sha256
		// - diffie-hellman-group-exchange-sha1
		// Those algorithms are insecure and may allow plaintext data to be recovered by an attacker.
		// UseInsecureCipher: true,
	}
	p.SSH.Server = bastion.Address
	if len(bastion.SSHAuthorizedKey) > 0 {
		p.SSH.KeyPath = util.ExpandPath(bastion.SSHAuthorizedKey)
		p.Log.Debugf("sshKeyPath: %s", p.SSH.KeyPath)
	}
	p.Log.Tracef("ssh -i %s %s@%s:%s", p.SSH.KeyPath, p.SSH.User, p.SSH.Server, p.SSH.Port)
}

// sshExecute выполнить комманду на удаленном компьютере и вернуть результат как строка
func (p *ProviderBase) sshExecute(command string) (stdOut string, stdErr string, err error) {
	if p.CmdFlags.DryRun {
		p.Log.Warnf("Dry RUN: ssh %s@%s -p %s \"%s\"", p.SSH.User, p.SSH.Server, p.SSH.Port, command)
	} else {
		stdOut, stdErr, _, err = p.SSH.Run(command, 60*time.Second)
	}
	return stdOut, stdErr, err
}

// Run command on remote machine
//   Example:
func (p *ProviderBase) Run(command string) (done bool, err error) {
	p.Log.Debugf("RUN command: %s", command)
	stdOut, stdErr, done, err := p.SSH.Run(command, 60*time.Second)
	if len(stdOut) > 0 {
		p.Log.Debugln("===== stdOut ======")
		p.Log.Debugf("%v", stdOut)
		p.Log.Debugln("===================")
	}
	if len(stdErr) > 0 {
		p.Log.Errorln("===== stdErr ======")
		p.Log.Errorf("%v", stdErr)
		p.Log.Errorln("===================")
	}
	return done, err
}

// Stream returns one channel that combines the stdout and stderr of the command
// as it is run on the remote machine, and another that sends true when the
// command is done. The sessions and channels will then be closed.
//  isPrint - выводить результат на экран или в лог
func (p *ProviderBase) sshStream(command string, isPrint bool) {
	if p.CmdFlags.DryRun {
		p.Log.Warnf("Dry Stream: ssh %s@%s -p %s \"%s\"", p.SSH.User, p.SSH.Server, p.SSH.Port, command)
	} else {
		// Call Run method with command you want to run on remote server.
		stdoutChan, stderrChan, doneChan, errChan, err :=  p.SSH.Stream(command, 60*time.Second)
		// Handle errors
		if err != nil {
			p.Log.Fatalln("Can't run remote command: " + err.Error())
		} else {
		// read from the output channel until the done signal is passed
			isTimeout := true
			loop:
				for {
					select {
					case isTimeout = <-doneChan:
						break loop
					case outline := <-stdoutChan:
						if isPrint && len(outline) > 0 {
							// fmt.Println("out:", outline)
							fmt.Println(outline)
						} else if len(outline) > 0 {
							p.Log.Infoln(outline)
						}
					case errline := <-stderrChan:
						if isPrint && len(errline) > 0 {
							// fmt.Println("err:", errline)
							fmt.Println(errline)
						} else if len(errline) > 0 {
							p.Log.Warnln(errline)
						}
					case err = <-errChan:
					}
				}

			// get exit code or command error.
			if err != nil {
				p.Log.Errorln("Error: " + err.Error())
			}

			// command time out
			if !isTimeout {
				p.Log.Errorln("Error: command timeout")
			}
		}
	}
}

// +kubebuilder:rbac:groups=k3s.bbox.kiev.ua,resources=clusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=k3s.bbox.kiev.ua,resources=clusters/status,verbs=get;update;patch

func (r *ClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("cluster", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	// return ctrl.NewControllerManagedBy(mgr).
	// 	For(&k3sv1alpha1.Cluster{}).
	// 	Complete(r)
	// TODO: ошибка Complete(r) 
	return nil
}
