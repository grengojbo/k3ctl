/*

Copyright © 2020 The k3d Author(s)

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
package cluster

import (
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	// "github.com/docker/go-connections/nat"
	// "github.com/docker/go-connections/nat"

	cliutil "github.com/grengojbo/k3ctl/cmd/util"
	"github.com/grengojbo/k3ctl/controllers"

	// k3dCluster "github.com/rancher/k3d/v4/pkg/client"
	conf "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/config"
	"github.com/grengojbo/k3ctl/pkg/types"

	// "github.com/rancher/k3d/v4/pkg/runtimes"

	"github.com/grengojbo/k3ctl/version"

	log "github.com/sirupsen/logrus"
)

var configFile string

const clusterCreateDescription = `
Create a new k3s cluster with containerized nodes (k3s in docker).
Every cluster will consist of one or more containers:
	- 1 (or more) server node container (k3s)
	- (optionally) 1 loadbalancer container as the entrypoint to the cluster (nginx)
	- (optionally) 1 (or more) agent node containers (k3s)
`

// func initConfig(args []string) {

// 	dryRun = viper.GetBool("dry-run")
// 	// Viper for pre-processed config options
// 	ppViper.SetEnvPrefix("K3S")

// 	// viper for the general config (file, env and non pre-processed flags)
// 	cfgViper.SetEnvPrefix("K3S")
// 	cfgViper.AutomaticEnv()

// 	cfgViper.SetConfigType("yaml")

// 	configFile = util.GerConfigFileName(args[0])
// 	cfgViper.SetConfigFile(configFile)
// 	// log.Tracef("Schema: %+v", conf.JSONSchema)

// 	// if err := config.ValidateSchemaFile(configFile, []byte(conf.JSONSchema)); err != nil {
// 	// 	log.Fatalf("Schema Validation failed for config file %s: %+v", configFile, err)
// 	// }

// 	// try to read config into memory (viper map structure)
// 	if err := cfgViper.ReadInConfig(); err != nil {
// 		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
// 			log.Fatalf("Config file %s not found: %+v", configFile, err)
// 		}
// 		// config file found but some other error happened
// 		log.Fatalf("Failed to read config file %s: %+v", configFile, err)
// 	}

// 	log.Infof("Using config file %s", cfgViper.ConfigFileUsed())
// 	// }

// 	if log.GetLevel() >= log.DebugLevel {
// 		c, _ := yaml.Marshal(cfgViper.AllSettings())
// 		log.Debugf("Configuration:\n%s", c)

// 		c, _ = yaml.Marshal(ppViper.AllSettings())
// 		log.Debugf("Additional CLI Configuration:\n%s", c)
// 	}
// }

// NewCmdClusterCreate returns a new cobra command
func NewCmdClusterCreate() *cobra.Command {

	// create new command
	cmd := &cobra.Command{
		Use:   "create NAME",
		Short: "Create a new cluster",
		Long:  clusterCreateDescription,
		// Args:  cobra.RangeArgs(0, 1), // exactly one cluster name can be set (default: k3d.DefaultClusterName)
		Args: cobra.ExactArgs(1), // exactly one name accepted // TODO: if not specified, inherit from cluster that the node shall belong to, if that is specified
		PreRunE: func(cmd *cobra.Command, args []string) error {
			
			cmdFlags.DryRun = viper.GetBool("dry-run")
			cmdFlags.DebugLogging = viper.GetBool("verbose")
			cmdFlags.TraceLogging = viper.GetBool("trace")
			
			clusterName = args[0]
			// NodeName = args[0]
			// --cluster
			// clusterName, err := cmd.Flags().GetString("cluster")
			// if err != nil {
			// 	log.Fatalln(err)
			// }
			ConfigFile = config.InitConfig(clusterName, cfgViper, ppViper)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {

			/*************************
			 * Compute Configuration *
			 *************************/
			 c, err := controllers.NewClusterFromConfig(cfgViper, cmdFlags)
			 if err != nil {
				 log.Fatalln(err)
			 }
 
			cfg, _ := yaml.Marshal(c.Cluster)

			log.Debugf("========== Simple Config ==========\n%s\n==========================\n", cfg)

			log.Infof("cni: %s (backend: %s)", c.Cluster.Spec.Networking.CNI, c.Cluster.Spec.Networking.Backend)
			log.Infof("secretsEncryption: %v", c.Cluster.Spec.Options.SecretsEncryption)
			log.Infof("datastore: %s", c.Cluster.Spec.Datastore.Provider)

			// c, err = applyCLIOverrides(c)
			// if err != nil {
			// 	log.Fatalf("Failed to apply CLI overrides: %+v", err)
			// }

			// if log.GetLevel() >= log.DebugLevel {
			// 	log.Debugf("========== Merged Simple Config ==========\n%+v\n==========================\n", cfg)
			// 	c, _ := yaml.Marshal(cfg)
			// 	log.Debugf("Merge Configuration:\n%s", c)
			// }

			/**************************************
			 * Transform, Process & Validate Configuration *
			 **************************************/

			// Set the name
			// if len(args) != 0 {
			// 	cfg.Name = args[0]
			// }

			// clusterConfig, err := config.TransformSimpleToClusterConfig(cmd.Context(), runtimes.SelectedRuntime, cfg)
			// if err != nil {
			// 	log.Fatalln(err)
			// }
			// log.Debugf("===== Merged Cluster Config =====\n%+v\n===== ===== =====\n", clusterConfig)

			// clusterConfig, err = config.ProcessClusterConfig(*clusterConfig)
			// if err != nil {
			// 	log.Fatalln(err)
			// }
			// log.Debugf("===== Processed Cluster Config =====\n%+v\n===== ===== =====\n", clusterConfig)

			// if err := config.ValidateClusterConfig(cmd.Context(), runtimes.SelectedRuntime, *clusterConfig); err != nil {
			// 	log.Fatalln("Failed Cluster Configuration Validation: ", err)
			// }

			/**************************************
			 * Create cluster if it doesn't exist *
			**************************************/

			if len(c.Cluster.Spec.Nodes) == 0 {
				log.Fatalln("Is Not Nodes to install k3s cluster")
			}

			// обновляем статус нод
			c.LoadNodeStatus()

			// servers, agents, err := util.GetGroupNodes(cfg.Spec.Nodes)
			if err = c.CreateK3sCluster(); err != nil {
				log.Fatalln(err)
			}
			
			// nodes, err := c.ListNodes()
			// if err != nil {
			// 	log.Errorf(err.Error())
			// }
			// for _, node := range nodes {
			// 	status := client.GetStatus(&node)
			// 	log.Infof("node: %s (%s)", node.GetObjectMeta().GetName(), status)
			// }
			
			// nodes, err := c.DescribeClusterNodes()
			// if err != nil {
			// 	log.Errorf(err.Error())
			// }
			// for _, node := range nodes {
			// 	y, _ := yaml.Marshal(node)
			// 	log.Debugf("========== Node Info ==========\n%s\n==========================\n", y)
			// }
			
			log.Infoln("Creating initializing server node")
			// masters := []conf.ContrelPlanNodes{}
			// for _, node := range servers {
			// 		if err := k3s.RunK3sCommand(bastion, &installk3sExec, dryRun); err != nil {
			// 			log.Fatalln(err.Error())
			// 		}
			// 		masters = append(masters, conf.ContrelPlanNodes{
			// 			Bastion: bastion,
			// 			Node:    node,
			// 		})
			// 		log.Infof("Name: %s (Role: %v) User: %v\n", node.Name, node.Role, cfg.GetUser(node.User).Name)
			// 		log.Infoln("-------------------")
			// 	}
			// }

			// if len(agents) > 0 {
			// 	// if len(masters) == 0 {
			// 	// 	log.Fatalln("Is NOT set control plane nodes")
			// 	// }
			// 	token, err := k3s.GetAgentToken(masters, dryRun)
			// 	if err != nil {
			// 		log.Fatalln(err.Error())
			// 	}
			// 	log.Infoln("=====================")
			// 	log.Infoln("Install agents")

			// 	// log.Debugf("K3S_TOKEN=%s", token)
			// 	for _, node := range agents {
			// 		if bastion, err := cfg.GetBastion(node.Bastion, node); err != nil {
			// 			log.Fatalln(err.Error())
			// 		} else {
			// 			apiServerAddres, err := cfg.GetAPIServerAddress(node, &cfg.Spec.Networking)
			// 			if err != nil {
			// 				log.Fatal(err)
			// 			}
			// 			// log.Warnf("apiServerAddresses: %s", apiServerAddres)

			// 			installk3sAgentExec := k3s.MakeAgentInstallExec(apiServerAddres, token, k3sOpt)
			// 			installk3sAgentExec.K3sChannel = cfg.Spec.K3sChannel
			// 			installk3sAgentExec.K3sVersion = cfg.Spec.KubernetesVersion
			// 			installk3sAgentExec.Node = node

			// 			if err := k3s.RunK3sCommand(bastion, &installk3sAgentExec, dryRun); err != nil {
			// 				log.Fatalln(err.Error())
			// 			}

			// 			log.Infof("Name: %s (Role: %v) User: %v\n", node.Name, node.Role, cfg.GetUser(node.User).Name)
			// 			// log.Debugln(bastion)
			// 			log.Infoln("-------------------")
			// 		}
			// 	}
			// }
			// log.Infoln("DRY RUN: ", dryRun)
			// // // check if a cluster with that name exists already
			// // if _, err := k3dCluster.ClusterGet(cmd.Context(), runtimes.SelectedRuntime, &clusterConfig.Cluster); err == nil {
			// // 	log.Fatalf("Failed to create cluster '%s' because a cluster with that name already exists", clusterConfig.Cluster.Name)
			// // }

			// // // create cluster
			// // if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
			// // 	log.Debugln("'--kubeconfig-update-default set: enabling wait-for-server")
			// // 	clusterConfig.ClusterCreateOpts.WaitForServer = true
			// // }
			// // //if err := k3dCluster.ClusterCreate(cmd.Context(), runtimes.SelectedRuntime, &clusterConfig.Cluster, &clusterConfig.ClusterCreateOpts); err != nil {
			// // if err := k3dCluster.ClusterRun(cmd.Context(), runtimes.SelectedRuntime, clusterConfig); err != nil {
			// // 	// rollback if creation failed
			// // 	log.Errorln(err)
			// // 	if cfg.Options.K3dOptions.NoRollback { // TODO: move rollback mechanics to pkg/
			// // 		log.Fatalln("Cluster creation FAILED, rollback deactivated.")
			// // 	}
			// // 	// rollback if creation failed
			// // 	log.Errorln("Failed to create cluster >>> Rolling Back")
			// // 	if err := k3dCluster.ClusterDelete(cmd.Context(), runtimes.SelectedRuntime, &clusterConfig.Cluster, k3d.ClusterDeleteOpts{SkipRegistryCheck: true}); err != nil {
			// // 		log.Errorln(err)
			// // 		log.Fatalln("Cluster creation FAILED, also FAILED to rollback changes!")
			// // 	}
			// // 	log.Fatalln("Cluster creation FAILED, all changes have been rolled back!")
			// // }
			// // log.Infof("Cluster '%s' created successfully!", clusterConfig.Cluster.Name)

			// /**************
			//  * Kubeconfig *
			//  **************/

			// // if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig && clusterConfig.KubeconfigOpts.SwitchCurrentContext {
			// // 	log.Infoln("--kubeconfig-update-default=false --> sets --kubeconfig-switch-context=false")
			// // 	clusterConfig.KubeconfigOpts.SwitchCurrentContext = false
			// // }

			// // if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig {
			// // 	log.Debugf("Updating default kubeconfig with a new context for cluster %s", clusterConfig.Cluster.Name)
			// // 	if _, err := k3dCluster.KubeconfigGetWrite(cmd.Context(), runtimes.SelectedRuntime, &clusterConfig.Cluster, "", &k3dCluster.WriteKubeConfigOptions{UpdateExisting: true, OverwriteExisting: false, UpdateCurrentContext: cfg.Options.KubeconfigOptions.SwitchCurrentContext}); err != nil {
			// // 		log.Warningln(err)
			// // 	}
			// // }

			// /*****************
			//  * User Feedback *
			//  *****************/

			// // // print information on how to use the cluster with kubectl
			// // log.Infoln("You can now use it like this:")
			// // if clusterConfig.KubeconfigOpts.UpdateDefaultKubeconfig && !clusterConfig.KubeconfigOpts.SwitchCurrentContext {
			// // 	fmt.Printf("kubectl config use-context %s\n", fmt.Sprintf("%s-%s", k3d.DefaultObjectNamePrefix, clusterConfig.Cluster.Name))
			// // } else if !clusterConfig.KubeconfigOpts.SwitchCurrentContext {
			// // 	if runtime.GOOS == "windows" {
			// // 		fmt.Printf("$env:KUBECONFIG=(%s kubeconfig write %s)\n", os.Args[0], clusterConfig.Cluster.Name)
			// // 	} else {
			// // 		fmt.Printf("export KUBECONFIG=$(%s kubeconfig write %s)\n", os.Args[0], clusterConfig.Cluster.Name)
			// // 	}
			// // }

			// // k3sup install --ip 192.168.192.103 --print-command \
			// // --k3s-extra-args="--tls-san developer.iwis.io --disable servicelb --disable traefik
			// // --cluster-cidr 10.42.0.0/19 --service-cidr 10.42.32.0/19 --cluster-dns 10.42.32.10
			// // --flannel-backend=none --secrets-encryption --node-taint CriticalAddonsOnly=true:NoExecute" \
			// // --user ubuntu  --local-path ~/.kube/developer.yaml --context developer
			fmt.Println("kubectl cluster-info :)")
		},
	}

	/***************
	 * Config File *
	 ***************/

	cmd.Flags().StringVarP(&configFile, "config", "c", "", "Path of a config file to use")
	if err := cobra.MarkFlagFilename(cmd.Flags(), "config", "yaml", "yml"); err != nil {
		log.Fatalln("Failed to mark flag 'config' as filename flag")
	}

	/***********************
	 * Pre-Processed Flags *
	 ***********************
	 *
	 * Flags that have a different style in the CLI than their internal representation.
	 * Also, we cannot set (viper) default values just here for those.
	 * Example:
	 *   CLI: `--api-port 0.0.0.0:6443`
	 *   Config File:
	 *	   exposeAPI:
	 *			 hostIP: 0.0.0.0
	 *       port: 6443
	 *
	 * Note: here we also use Slice-type flags instead of Array because of https://github.com/spf13/viper/issues/380
	 */

	cmd.Flags().String("api-port", "", "Specify the Kubernetes API server port exposed on the LoadBalancer (Format: `[HOST:]HOSTPORT`)\n - Example: `k3d cluster create --servers 3 --api-port 0.0.0.0:6550`")
	_ = ppViper.BindPFlag("cli.api-port", cmd.Flags().Lookup("api-port"))

	cmd.Flags().StringArrayP("env", "e", nil, "Add environment variables to nodes (Format: `KEY[=VALUE][@NODEFILTER[;NODEFILTER...]]`\n - Example: `k3d cluster create --agents 2 -e \"HTTP_PROXY=my.proxy.com\" -e \"SOME_KEY=SOME_VAL@server[0]\"`")
	_ = ppViper.BindPFlag("cli.env", cmd.Flags().Lookup("env"))

	cmd.Flags().StringArrayP("volume", "v", nil, "Mount volumes into the nodes (Format: `[SOURCE:]DEST[@NODEFILTER[;NODEFILTER...]]`\n - Example: `k3d cluster create --agents 2 -v /my/path@agent[0,1] -v /tmp/test:/tmp/other@server[0]`")
	_ = ppViper.BindPFlag("cli.volumes", cmd.Flags().Lookup("volume"))

	cmd.Flags().StringArrayP("port", "p", nil, "Map ports from the node containers to the host (Format: `[HOST:][HOSTPORT:]CONTAINERPORT[/PROTOCOL][@NODEFILTER]`)\n - Example: `k3d cluster create --agents 2 -p 8080:80@agent[0] -p 8081@agent[1]`")
	_ = ppViper.BindPFlag("cli.ports", cmd.Flags().Lookup("port"))

	cmd.Flags().StringArrayP("label", "l", nil, "Add label to node container (Format: `KEY[=VALUE][@NODEFILTER[;NODEFILTER...]]`\n - Example: `k3d cluster create --agents 2 -l \"my.label@agent[0,1]\" -v \"other.label=somevalue@server[0]\"`")
	_ = ppViper.BindPFlag("cli.labels", cmd.Flags().Lookup("label"))

	/******************
	 * "Normal" Flags *
	 ******************
	 *
	 * No pre-processing needed on CLI level.
	 * Bound to Viper config value.
	 * Default Values set via Viper.
	 */

	cmd.Flags().IntP("servers", "s", 0, "Specify how many servers you want to create")
	_ = cfgViper.BindPFlag("spec.servers", cmd.Flags().Lookup("servers"))
	cfgViper.SetDefault("spec.servers", 1)

	cmd.Flags().IntP("agents", "a", 0, "Specify how many agents you want to create")
	_ = cfgViper.BindPFlag("spec.agents", cmd.Flags().Lookup("agents"))
	cfgViper.SetDefault("spec.agents", 0)

	// // cmd.Flags().StringP("image", "i", "", "Specify k3s image that you want to use for the nodes")
	// // _ = cfgViper.BindPFlag("image", cmd.Flags().Lookup("image"))
	// // cfgViper.SetDefault("image", fmt.Sprintf("%s:%s", k3d.DefaultK3sImageRepo, version.GetK3sVersion(false)))

	// cmd.Flags().String("network", "", "Join an existing network")
	// _ = cfgViper.BindPFlag("network", cmd.Flags().Lookup("network"))

	cmd.Flags().String("token", "", "Specify a cluster token. By default, we generate one.")
	_ = cfgViper.BindPFlag("spec.token", cmd.Flags().Lookup("token"))

	cmd.Flags().Bool("secrets-encryption", false, "Enable Secret encryption at rest")
	_ = cfgViper.BindPFlag("spec.options.secretsEncryption", cmd.Flags().Lookup("secrets-encryption"))

	cmd.Flags().Bool("selinux", false, "To leverage SELinux, specify the flag when starting K3s servers and agents.")
	_ = cfgViper.BindPFlag("spec.options.selinux", cmd.Flags().Lookup("selinux"))

	cmd.Flags().Bool("rootless", false, "Running Servers and Agents with Rootless")
	_ = cfgViper.BindPFlag("spec.options.rootless", cmd.Flags().Lookup("rootless"))

	cmd.Flags().Bool("wait", true, "Wait for the server(s) to be ready before returning. Use '--timeout DURATION' to not wait forever.")
	_ = cfgViper.BindPFlag("spec.options.wait", cmd.Flags().Lookup("wait"))

	cmd.Flags().Duration("timeout", 0*time.Second, "Rollback changes if cluster couldn't be created in specified duration.")
	_ = cfgViper.BindPFlag("spec.options.timeout", cmd.Flags().Lookup("timeout"))

	cmd.Flags().Bool("kubeconfig-update-default", true, "Directly update the default kubeconfig with the new cluster's context")
	_ = cfgViper.BindPFlag("spec.kubeconfig.updatedefaultkubeconfig", cmd.Flags().Lookup("kubeconfig-update-default"))

	cmd.Flags().Bool("kubeconfig-switch-context", true, "Directly switch the default kubeconfig's current-context to the new cluster's context (requires --kubeconfig-update-default)")
	_ = cfgViper.BindPFlag("spec.kubeconfig.switchcurrentcontext", cmd.Flags().Lookup("kubeconfig-switch-context"))

	cmd.Flags().String("k3s-channel", version.PinnedK3sChannel, fmt.Sprintf("Release channel: stable, latest, or i.e. %v", version.PinnedK3sChannel))
	_ = cfgViper.BindPFlag("spec.channel", cmd.Flags().Lookup("k3s-channel"))

	cmd.Flags().String("k3s-version", "", "Set a version to install, overrides k3s-version")
	// log.Infoln("k3s-version: ", cmd.Flags().Lookup("k3s-version").Value.String())
	_ = cfgViper.BindPFlag("spec.KubernetesVersion", cmd.Flags().Lookup("k3s-version"))

	cmd.Flags().Bool("no-lb", false, "Disable the creation of a LoadBalancer in front of the server nodes")
	_ = cfgViper.BindPFlag("spec.options.disableloadbalancer", cmd.Flags().Lookup("no-lb"))

	cmd.Flags().Bool("no-ingress", false, "Disable the creation of a Ingress Controller in front of the server nodes")
	_ = cfgViper.BindPFlag("spec.options.disableIngress", cmd.Flags().Lookup("no-ingress"))

	// cmd.Flags().Bool("no-rollback", false, "Disable the automatic rollback actions, if anything goes wrong")
	// _ = cfgViper.BindPFlag("spec.options.disablerollback", cmd.Flags().Lookup("no-rollback"))

	// cmd.Flags().Bool("no-hostip", false, "Disable the automatic injection of the Host IP as 'host.options.internal' into the containers and CoreDNS")
	// _ = cfgViper.BindPFlag("spec.options.disablehostipinjection", cmd.Flags().Lookup("no-hostip"))

	// cmd.Flags().String("gpus", "", "GPU devices to add to the cluster node containers ('all' to pass all GPUs) [From docker]")
	// _ = cfgViper.BindPFlag("spec.runtime.gpurequest", cmd.Flags().Lookup("gpus"))

	// cmd.Flags().String("servers-memory", "", "Memory limit imposed on the server nodes [From docker]")
	// _ = cfgViper.BindPFlag("spec.runtime.serversmemory", cmd.Flags().Lookup("servers-memory"))

	// cmd.Flags().String("agents-memory", "", "Memory limit imposed on the agents nodes [From docker]")
	// _ = cfgViper.BindPFlag("spec.runtime.agentsmemory", cmd.Flags().Lookup("agents-memory"))

	// /* Image Importing */
	// cmd.Flags().Bool("no-image-volume", false, "Disable the creation of a volume for importing images")
	// _ = cfgViper.BindPFlag("spec.options.disableimagevolume", cmd.Flags().Lookup("no-image-volume"))

	/* Registry */
	cmd.Flags().StringArray("registry-use", nil, "Connect to one or more k3d-managed registries running locally")
	_ = cfgViper.BindPFlag("spec.registries.use", cmd.Flags().Lookup("registry-use"))

	cmd.Flags().Bool("registry-create", false, "Create a k3d-managed registry and connect it to the cluster")
	_ = cfgViper.BindPFlag("spec.registries.create", cmd.Flags().Lookup("registry-create"))

	cmd.Flags().String("registry-config", "", "Specify path to an extra registries.yaml file")
	_ = cfgViper.BindPFlag("spec.registries.config", cmd.Flags().Lookup("registry-config"))

	/* k3s */
	cmd.Flags().StringArray("k3s-server-arg", nil, "Additional args passed to the `k3s server` command on server nodes (new flag per arg)")
	_ = cfgViper.BindPFlag("spec.k3s.extraserverargs", cmd.Flags().Lookup("k3s-server-arg"))

	cmd.Flags().StringArray("k3s-agent-arg", nil, "Additional args passed to the `k3s agent` command on agent nodes (new flag per arg)")
	_ = cfgViper.BindPFlag("spec.k3s.extraagentargs", cmd.Flags().Lookup("k3s-agent-arg"))

	/* Subcommands */

	// done
	return cmd
}

func applyCLIOverrides(cfg conf.Cluster) (conf.Cluster, error) {

	/****************************
	 * Parse and validate flags *
	 ****************************/

	// -> API-PORT
	// parse the port mapping
	// var (
	// 	err       error
	// 	exposeAPI *types.ExposureOpts
	// )

	// // Apply config file values as defaults
	// exposeAPI = &types.ExposureOpts{
	// 	PortMapping: nat.PortMapping{
	// 		Binding: nat.PortBinding{
	// 			HostIP: cfg.Spec.HostIP,
	// 			// HostPort: cfg.Spec.HostPort,
	// 		},
	// 	},
	// 	Host: cfg.Spec.Host,
	// }

	// Overwrite if cli arg is set
	// if ppViper.IsSet("cli.api-port") {
	// 	if cfg.ExposeAPI.HostPort != "" {
	// 		log.Debugf("Overriding pre-defined kubeAPI Exposure Spec %+v with CLI argument %s", cfg.ExposeAPI, ppViper.GetString("cli.api-port"))
	// 	}
	// 	exposeAPI, err = cliutil.ParsePortExposureSpec(ppViper.GetString("cli.api-port"), k3d.DefaultAPIPort)
	// 	if err != nil {
	// 		return cfg, err
	// 	}
	// }

	// Set to random port if port is empty string
	// if len(exposeAPI.Binding.HostPort) == 0 {
	// 	exposeAPI, err = cliutil.ParsePortExposureSpec("random", k3d.DefaultAPIPort)
	// 	if err != nil {
	// 		return cfg, err
	// 	}
	// }

	// cfg.ExposeAPI = conf.SimpleExposureOpts{
	// 	Host:     exposeAPI.Host,
	// 	HostIP:   exposeAPI.Binding.HostIP,
	// 	HostPort: exposeAPI.Binding.HostPort,
	// }

	// -> VOLUMES
	// volumeFilterMap will map volume mounts to applied node filters
	volumeFilterMap := make(map[string][]string, 1)
	for _, volumeFlag := range ppViper.GetStringSlice("cli.volumes") {

		// split node filter from the specified volume
		volume, filters, err := cliutil.SplitFiltersFromFlag(volumeFlag)
		if err != nil {
			log.Fatalln(err)
		}

		if strings.Contains(volume, types.DefaultRegistriesFilePath) && (cfg.Spec.Addons.Registries.Create || cfg.Spec.Addons.Registries.Config != "" || len(cfg.Spec.Addons.Registries.Use) != 0) {
			log.Warnf("Seems like you're mounting a file at '%s' while also using a referenced registries config or k3d-managed registries: Your mounted file will probably be overwritten!", types.DefaultRegistriesFilePath)
		}

		// create new entry or append filter to existing entry
		if _, exists := volumeFilterMap[volume]; exists {
			volumeFilterMap[volume] = append(volumeFilterMap[volume], filters...)
		} else {
			volumeFilterMap[volume] = filters
		}
	}

	for volume, nodeFilters := range volumeFilterMap {
		cfg.Spec.Volumes = append(cfg.Spec.Volumes, conf.VolumeWithNodeFilters{
			Volume:      volume,
			NodeFilters: nodeFilters,
		})
	}

	log.Tracef("VolumeFilterMap: %+v", volumeFilterMap)

	// -> PORTS
	portFilterMap := make(map[string][]string, 1)
	for _, portFlag := range ppViper.GetStringSlice("cli.ports") {
		// split node filter from the specified volume
		portmap, filters, err := cliutil.SplitFiltersFromFlag(portFlag)
		if err != nil {
			log.Fatalln(err)
		}

		if len(filters) > 1 {
			log.Fatalln("Can only apply a Portmap to one node")
		}

		// create new entry or append filter to existing entry
		if _, exists := portFilterMap[portmap]; exists {
			log.Fatalln("Same Portmapping can not be used for multiple nodes")
		} else {
			portFilterMap[portmap] = filters
		}
	}

	// for port, nodeFilters := range portFilterMap {
	// 	cfg.Ports = append(cfg.Ports, conf.PortWithNodeFilters{
	// 		Port:        port,
	// 		NodeFilters: nodeFilters,
	// 	})
	// }

	log.Tracef("PortFilterMap: %+v", portFilterMap)

	// --label
	// labelFilterMap will add container label to applied node filters
	labelFilterMap := make(map[string][]string, 1)
	for _, labelFlag := range ppViper.GetStringSlice("cli.labels") {

		// split node filter from the specified label
		label, nodeFilters, err := cliutil.SplitFiltersFromFlag(labelFlag)
		if err != nil {
			log.Fatalln(err)
		}

		// create new entry or append filter to existing entry
		if _, exists := labelFilterMap[label]; exists {
			labelFilterMap[label] = append(labelFilterMap[label], nodeFilters...)
		} else {
			labelFilterMap[label] = nodeFilters
		}
	}

	for label, nodeFilters := range labelFilterMap {
		cfg.Spec.Labels = append(cfg.Spec.Labels, conf.LabelWithNodeFilters{
			Label:       label,
			NodeFilters: nodeFilters,
		})
	}

	log.Tracef("LabelFilterMap: %+v", labelFilterMap)

	// --env
	// envFilterMap will add container env vars to applied node filters
	envFilterMap := make(map[string][]string, 1)
	for _, envFlag := range ppViper.GetStringSlice("cli.env") {

		// split node filter from the specified env var
		env, filters, err := cliutil.SplitFiltersFromFlag(envFlag)
		if err != nil {
			log.Fatalln(err)
		}

		// create new entry or append filter to existing entry
		if _, exists := envFilterMap[env]; exists {
			envFilterMap[env] = append(envFilterMap[env], filters...)
		} else {
			envFilterMap[env] = filters
		}
	}

	for envVar, nodeFilters := range envFilterMap {
		cfg.Spec.Env = append(cfg.Spec.Env, conf.EnvVarWithNodeFilters{
			EnvVar:      envVar,
			NodeFilters: nodeFilters,
		})
	}

	log.Tracef("EnvFilterMap: %+v", envFilterMap)

	return cfg, nil
}
