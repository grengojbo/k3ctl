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
package node

import (
	// "fmt"
	// "strings"
	// "time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// "github.com/spf13/viper"

	// "github.com/grengojbo/k3ctl/pkg/config"
	"github.com/grengojbo/k3ctl/pkg/types"
	// "github.com/grengojbo/k3ctl/pkg/util"
	// k3s "github.com/grengojbo/k3ctl/pkg/k3s"
	// conf "github.com/grengojbo/k3ctl/api/v1alpha1"
	// // dockerunits "github.com/docker/go-units"
	// // cliutil "github.com/rancher/k3d/v5/cmd/util"
	// // k3dc "github.com/rancher/k3d/v5/pkg/client"
	// // l "github.com/rancher/k3d/v5/pkg/logger"
	log "github.com/sirupsen/logrus"
	// // "github.com/rancher/k3d/v5/pkg/runtimes"
	// // k3d "github.com/rancher/k3d/v5/pkg/types"
	// // "github.com/rancher/k3d/v5/version"
)

type nodeDeleteFlags struct {
	Cluster string
	All               bool
	IncludeRegistries bool
}

// NewCmdNodeCreate returns a new cobra command
func NewCmdNodeDelete() *cobra.Command {

	flags := nodeDeleteFlags{}
	// createNodeOpts := k3d.NodeCreateOpts{}

	// create new command
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete node(s)",
		Long:  `Delete k3s node from cluster.`,
		Args:  cobra.ExactArgs(1), // exactly one name accepted // TODO: if not specified, inherit from cluster that the node shall belong to, if that is specified
		// ValidArgsFunction: cliutil.ValidArgsAvailableNodes,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			DryRun = viper.GetBool("dry-run")
			
			// nodes := parseDeleteNodeCmd(cmd, args, &flags)
			// nodeDeleteOpts := k3d.NodeDeleteOpts{SkipLBUpdate: flags.All} // do not update LB, if we're deleting all nodes anyway

			// NodeName = args[0]
			// // --cluster
			// ClusterName, err := cmd.Flags().GetString("cluster")
			// if err != nil {
			// 	log.Fatalln(err)
			// }
			// ConfigFile = config.InitConfig(ClusterName, CfgViper, PpViper)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			isDeleteNode := false
			// /*************************
			//  * Compute Configuration *
			//  *************************/
			//  cfg, err := config.FromViperSimple(CfgViper)
			//  if err != nil {
			// 	 log.Fatalln(err)
			//  }
			// if len(cfg.Spec.Nodes) == 0 {
			// 	log.Fatalln("Is Not Nodes to install k3s cluster")
			// }

			// servers, agents, err := util.GetGroupNodes(cfg.Spec.Nodes)
			// if err != nil {
			// 	log.Fatalln(err)
			// }
			// k3sOpt := k3s.K3sExecOptions{
			// 	// 	NoExtras:     k3sNoExtras,
			// 	ExtraArgs:           cfg.Spec.K3sOptions.ExtraServerArgs,
			// 	Ingress:             cfg.Spec.Addons.Ingress.Name,
			// 	DisableLoadbalancer: cfg.Spec.Options.DisableLoadbalancer,
			// 	DisableIngress:      cfg.Spec.Options.DisableIngress,
			// 	SecretsEncryption:   cfg.Spec.Options.SecretsEncryption,
			// 	SELinux:             cfg.Spec.Options.SELinux,
			// 	Rootless:            cfg.Spec.Options.Rootless,
			// 	LoadBalancer:        &cfg.Spec.LoadBalancer,
			// 	Networking:          &cfg.Spec.Networking,
			// }
			// masters := []conf.ContrelPlanNodes{}
			// for _, node := range servers {
			// 	if bastion, err := cfg.GetBastion(node.Bastion, node); err != nil {
			// 		log.Fatalln(err.Error())
			// 	} else {
			// 		masters = append(masters, conf.ContrelPlanNodes{
			// 			Bastion: bastion,
			// 			Node: node,
			// 		})
			// 		if node.Name == NodeName {
			// 			log.Infof("TODO: Add master Node: %s", node.Name)		
			// 		}
			// 	}
			// }

			// token, err := k3s.GetAgentToken(masters, DryRun)
			// 	if err != nil {
			// 		log.Fatalln(err.Error())
			// 	}
				
			// 	// log.Debugf("K3S_TOKEN=%s", token)
			// 	for _, node := range agents {
			// 		if node.Name == NodeName {
			// 			if bastion, err := cfg.GetBastion(node.Bastion, node); err != nil {
			// 				log.Fatalln(err.Error())
			// 			} else {
			// 				apiServerAddres, err := cfg.GetAPIServerAddress(node, &cfg.Spec.Networking)
			// 				if err != nil {
			// 					log.Fatal(err)
			// 				}
			// 				cnt := cfg.GetNodeLabels(node)
			// 				log.Warnf("=-> cnt: %d", cnt)
			// 				// log.Warnf("apiServerAddresses: %s", apiServerAddres)
			// 				installk3sAgentExec := k3s.MakeAgentInstallExec(apiServerAddres, token, k3sOpt)
			// 				installk3sAgentExec.K3sChannel = cfg.Spec.K3sChannel
			// 				installk3sAgentExec.K3sVersion = cfg.Spec.KubernetesVersion
			// 				installk3sAgentExec.Node = node

			// 				if err := k3s.RunK3sCommand(bastion, &installk3sAgentExec, DryRun); err != nil {
			// 					log.Fatalln(err.Error())
			// 				}
			// 				log.Infof("Successfully added Agent Node: %s", node.Name)
			// 				isAddNode = true
			// 			}
			// 		}
			// 	}
			// // node, clusterName := parseCreateNodeCmd(cmd, args)
			// // if strings.HasPrefix(clusterName, "https://") {
			// // 	l.Log().Infof("Adding %d node(s) to the remote cluster '%s'...", len(nodes), clusterName)
			// // 	if err := k3dc.NodeAddToClusterMultiRemote(cmd.Context(), runtimes.SelectedRuntime, nodes, clusterName, createNodeOpts); err != nil {
			// // 		l.Log().Fatalf("failed to add %d node(s) to the remote cluster '%s': %v", len(nodes), clusterName, err)
			// // 	}
			// // } else {
			// // 	l.Log().Infof("Adding %d node(s) to the runtime local cluster '%s'...", len(nodes), clusterName)
			// // 	if err := k3dc.NodeAddToClusterMulti(cmd.Context(), runtimes.SelectedRuntime, nodes, &k3d.Cluster{Name: clusterName}, createNodeOpts); err != nil {
			// // 		l.Log().Fatalf("failed to add %d node(s) to the runtime local cluster '%s': %v", len(nodes), clusterName, err)
			// // 	}
			// // }
			if !isDeleteNode {
				log.Errorf("Is NOT set node: %v", NodeName)
			}
		},
	}

	// // add flags
	cmd.Flags().StringVarP(&flags.Cluster, "cluster", "c", types.DefaultClusterName, "Cluster URL or k3d cluster name to connect to.")
	cmd.Flags().BoolVarP(&flags.All, "all", "a", false, "Delete all existing nodes")
	cmd.Flags().BoolVarP(&flags.IncludeRegistries, "registries", "r", false, "Also delete registries")
	
	// done
	return cmd
}

// parseDeleteNodeCmd parses the command input into variables required to delete nodes
// func parseDeleteNodeCmd(cmd *cobra.Command, args []string, flags *nodeDeleteFlags) []*k3d.Node {

// 	var nodes []*k3d.Node
// 	var err error

// 	// --all
// 	if flags.All {
// 		if !flags.IncludeRegistries {
// 			l.Log().Infoln("Didn't set '--registries', so won't delete registries.")
// 		}
// 		nodes, err = client.NodeList(cmd.Context(), runtimes.SelectedRuntime)
// 		if err != nil {
// 			l.Log().Fatalln(err)
// 		}
// 		include := k3d.ClusterInternalNodeRoles
// 		exclude := []k3d.Role{}
// 		if flags.IncludeRegistries {
// 			include = append(include, k3d.RegistryRole)
// 		}
// 		nodes = client.NodeFilterByRoles(nodes, include, exclude)
// 		return nodes
// 	}

// 	if !flags.All && len(args) < 1 {
// 		l.Log().Fatalln("Expecting at least one node name if `--all` is not set")
// 	}

// 	for _, name := range args {
// 		node, err := client.NodeGet(cmd.Context(), runtimes.SelectedRuntime, &k3d.Node{Name: name})
// 		if err != nil {
// 			l.Log().Fatalln(err)
// 		}
// 		nodes = append(nodes, node)
// 	}

// 	return nodes
// }
