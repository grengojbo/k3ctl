package debug

import (
	"time"

	cliutil "github.com/grengojbo/k3ctl/cmd/util"
	"github.com/grengojbo/k3ctl/controllers"
	"github.com/grengojbo/k3ctl/pkg/cert"
	"github.com/grengojbo/k3ctl/pkg/config"
	"github.com/grengojbo/k3ctl/pkg/types"
	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
)

// var ClusterName string
// var NodeName string
var ConfigFile string
var cfgViper = viper.New()
var ppViper = viper.New()

var mode string
var clusterName string
var cmdFlags types.CmdFlags

// NewCmdDebug returns a new cobra command
func NewCmdDebug() *cobra.Command {
	cmd := &cobra.Command{
		Use:    "debug",
		Hidden: true,
		Short:  "Debug k3ctl cluster(s)",
		Long:   `Debug k3ctl cluster(s)`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				log.Errorln("Couldn't get help text")
				log.Fatalln(err)
			}
		},
	}

	cmd.AddCommand(NewCmdDebugNodeList())
	cmd.AddCommand(NewCmdDebugCert())

	return cmd
}

// NewCmdDebugCert returns a new cobra command
func NewCmdDebugCert() *cobra.Command {

	// create new command
	cmd := &cobra.Command{
		Use:       "cert",
		Short:     "Ceertificate Kubernetes",
		Long:      `Ceertificate Kubernetes API server.`,
		ValidArgs: util.ListClusterName(),
		// Args:  cobra.RangeArgs(0, 1), // exactly one cluster name can be set (default: k3d.DefaultClusterName)
		Args: cobra.ExactValidArgs(1),
		// Args: cobra.MinimumNArgs(1),
		PreRunE: func(cmd *cobra.Command, args []string) error {

			cmdFlags.DryRun = viper.GetBool("dry-run")
			cmdFlags.DebugLogging = viper.GetBool("verbose")
			cmdFlags.TraceLogging = viper.GetBool("trace")

			clusterName = args[0]
			ConfigFile = config.InitConfig(clusterName, cfgViper, ppViper)

			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// /*************************
			//  * Compute Configuration *
			//  *************************/
			c, err := controllers.NewClusterFromConfig(cfgViper, cmdFlags)
			if err != nil {
				log.Fatalln(err)
			}

			cfg, _ := yaml.Marshal(c.Cluster)
			log.Tracef("Simple Config:\n%s", cfg)

			apiServers, err := c.GetAPIServerAddress()
			if err != nil {
				log.Fatalf("[GetAPIServerAddress] %s", err.Error())
			}

			hosts := []cert.HostInfo{}
			timeout := 5 * time.Second
			// expires := 7 * 24 * time.Hour
			for _, apiServer := range apiServers {
				host := cert.HostInfo{
					Host: apiServer,
					Port: int(c.GetAPIServerPort()),
				}
				if err := host.GetCerts(timeout); err == nil {
					log.Infof("Certificates: %+v", host.Certs)
					hosts = append(hosts, host)
					continue
				} else {
					log.Errorf("%s: %s", apiServer, err.Error())
				}
				// enc := json.NewEncoder(os.Stdout)
				// enc.SetIndent("", "  ")
				// enc.SetEscapeHTML(false)
				// err := enc.Encode(&host)
				// if err != nil {
				// 	log.Fatalln(err)
				// }
				// log.Infof("[GetAPIServerAddress] %s", apiServer)
				// cert, err := cliutil.GetCertificate(apiServer, timeout, expires)
				// if err != nil {
				// 	log.Fatalf("[GetCertificate] %s", err.Error())
				// }
				// log.Infof("[GetCertificate] %s", cert)
				log.Warnf("API Server: %s port: %d", apiServers, c.GetAPIServerPort())
			}

			if len(hosts) > 0 {

				if err := cert.ShowCerts(hosts, mode); err != nil {
					log.Fatalf("[CreateCertificate] %s", err.Error())
				}
			}
			// if len(c.Cluster.Spec.Nodes) == 0 {
			// 	log.Fatalln("Is Not Nodes to install k3s cluster")
			// }

			// // // isAddNode := false
			// // if ok := c.AddNodeToCluster(NodeName); ok {
			// // 	isAddNode = true
			// // }
			// // if !isAddNode {
			// // 	log.Errorf("Is NOT set node: %v", NodeName)
			// // } else {
			// // 	log.Infof("Successfully added %s node(s)!", NodeName)
			// // }
			log.Warnf("TODO: certificate expired")
		},
	}
	return cmd
}

// NewCmdNodeCreate returns a new cobra command
func NewCmdDebugNodeList() *cobra.Command {

	// create new command
	cmd := &cobra.Command{
		Use:   "node-list",
		Short: "Lisr Kubernetes node",
		Long:  `Add a new containerized Kubernetes node.`,
		PreRunE: func(cmd *cobra.Command, args []string) error {

			cmdFlags.DryRun = viper.GetBool("dry-run")
			cmdFlags.DebugLogging = viper.GetBool("verbose")
			cmdFlags.TraceLogging = viper.GetBool("trace")

			// NodeName = args[0]
			// // --cluster
			// clusterName, err := cmd.Flags().GetString("cluster")
			// if err != nil {
			// 	log.Fatalln(err)
			// }
			// ConfigFile = config.InitConfig(clusterName, CfgViper, PpViper)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// /*************************
			//  * Compute Configuration *
			//  *************************/
			// c, err := controllers.NewClusterFromConfig(CfgViper, cmdFlags)
			// if err != nil {
			// 	log.Fatalln(err)
			// }

			// cfg, _ := yaml.Marshal(c.Cluster)
			// log.Tracef("Simple Config:\n%s", cfg)

			// if len(c.Cluster.Spec.Nodes) == 0 {
			// 	log.Fatalln("Is Not Nodes to install k3s cluster")
			// }

			nodes, _ := cliutil.ValidArgsAvailableNodes(cmd, args, "toComplete")
			log.Warnf("Nodes: %+v", nodes)

			// // // isAddNode := false
			// // if ok := c.AddNodeToCluster(NodeName); ok {
			// // 	isAddNode = true
			// // }
			// // if !isAddNode {
			// // 	log.Errorf("Is NOT set node: %v", NodeName)
			// // } else {
			// // 	log.Infof("Successfully added %s node(s)!", NodeName)
			// // }
		},
	}

	// // add flags
	// cmd.Flags().Int("replicas", 1, "Number of replicas of this node specification.")
	// cmd.Flags().String("role", string(k3d.AgentRole), "Specify node role [server, agent]")
	// if err := cmd.RegisterFlagCompletionFunc("role", util.ValidArgsNodeRoles); err != nil {
	// 	l.Log().Fatalln("Failed to register flag completion for '--role'", err)
	// }
	cmd.Flags().StringP("cluster", "c", types.DefaultClusterName, "Cluster URL or k3s cluster name to connect to.")
	cmd.MarkFlagRequired("cluster")
	// if err := cmd.RegisterFlagCompletionFunc("cluster", util.ValidArgsAvailableClusters); err != nil {
	if err := cmd.RegisterFlagCompletionFunc("cluster", cliutil.ValidArgsAvailableClusters); err != nil {
		log.Fatalln("Failed to register flag completion for '--cluster'", err)
	}

	cmd.Flags().StringVarP(&mode, "output", "o", "text", "Output result [text, json, none]")

	// cmd.Flags().StringP("image", "i", fmt.Sprintf("%s:%s", k3d.DefaultK3sImageRepo, version.K3sVersion), "Specify k3s image used for the node(s)")
	// cmd.Flags().String("memory", "", "Memory limit imposed on the node [From docker]")

	// cmd.Flags().BoolVar(&createNodeOpts.Wait, "wait", true, "Wait for the node(s) to be ready before returning.")
	// cmd.Flags().DurationVar(&createNodeOpts.Timeout, "timeout", 0*time.Second, "Maximum waiting time for '--wait' before canceling/returning.")

	// cmd.Flags().StringSliceP("runtime-label", "", []string{}, "Specify container runtime labels in format \"foo=bar\"")
	// cmd.Flags().StringSliceP("k3s-node-label", "", []string{}, "Specify k3s node labels in format \"foo=bar\"")

	// cmd.Flags().StringSliceP("network", "n", []string{}, "Add node to (another) runtime network")

	// cmd.Flags().StringVarP(&createNodeOpts.ClusterToken, "token", "t", "", "Override cluster token (required when connecting to an external cluster)")

	// done
	return cmd
}
