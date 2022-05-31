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
	"gopkg.in/yaml.v2"

	// conf "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/controllers"
	"github.com/grengojbo/k3ctl/pkg/config"

	// k3s "github.com/grengojbo/k3ctl/pkg/k3s"

	"github.com/grengojbo/k3ctl/pkg/types"

	// dockerunits "github.com/docker/go-units"
	// cliutil "github.com/rancher/k3d/v5/cmd/util"
	// k3dc "github.com/rancher/k3d/v5/pkg/client"
	// l "github.com/rancher/k3d/v5/pkg/logger"
	log "github.com/sirupsen/logrus"
	// "github.com/rancher/k3d/v5/pkg/runtimes"
	// k3d "github.com/rancher/k3d/v5/pkg/types"
	// "github.com/rancher/k3d/v5/version"
)

// var clusterName string
// var cmdFlags types.CmdFlags

// NewCmdNodeCreate returns a new cobra command
func NewCmdNodeDelete() *cobra.Command {

	// createNodeOpts := k3d.NodeCreateOpts{}

	// create new command
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete k3s node",
		Long:  `Delete k3s node from cluster.`,
		// Args:  cobra.ExactArgs(1), // exactly one name accepted // TODO: if not specified, inherit from cluster that the node shall belong to, if that is specified
		PreRunE: func(cmd *cobra.Command, args []string) error {
			// DryRun = viper.GetBool("dry-run")

			cmdFlags.DryRun = viper.GetBool("dry-run")
			cmdFlags.DebugLogging = viper.GetBool("verbose")
			cmdFlags.TraceLogging = viper.GetBool("trace")

			if len(args) > 0 {
				NodeName = args[0]
			}
			// --cluster
			clusterName, err := cmd.Flags().GetString("cluster")
			if err != nil {
				log.Fatalln(err)
			}
			ConfigFile = config.InitConfig(clusterName, CfgViper, PpViper)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			/*************************
			* Compute Configuration *
			*************************/
			c, err := controllers.NewClusterFromConfig(CfgViper, cmdFlags)
			if err != nil {
				log.Fatalln(err)
			}

			cfg, _ := yaml.Marshal(c.Cluster)
			log.Tracef("Simple Config:\n%s", cfg)

			deleteAllNode := false
			cnt := c.DeleteNode(NodeName, deleteAllNode)
			if err != nil {
				log.Errorf("---------- cobra ------------")
				log.Errorf(err.Error())
			}

			log.Infof("Successfully deleted %d node(s)!", cnt)
		},
	}

	// // add flags
	// cmd.Flags().Int("replicas", 1, "Number of replicas of this node specification.")
	// cmd.Flags().String("role", string(k3d.AgentRole), "Specify node role [server, agent]")
	// if err := cmd.RegisterFlagCompletionFunc("role", util.ValidArgsNodeRoles); err != nil {
	// 	l.Log().Fatalln("Failed to register flag completion for '--role'", err)
	// }
	cmd.Flags().StringVarP(&clusterName, "cluster", "c", types.DefaultClusterName, "Cluster URL or k3s cluster name to connect to.")
	// cmd.Flags().StringP("cluster", "c", types.DefaultClusterName, "Cluster URL or k3s cluster name to connect to.")
	// if err := cmd.RegisterFlagCompletionFunc("cluster", util.ValidArgsAvailableClusters); err != nil {
	// 	log.Fatalln("Failed to register flag completion for '--cluster'", err)
	// }

	// done
	return cmd
}
