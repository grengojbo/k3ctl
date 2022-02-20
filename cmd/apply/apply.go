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
package app

import (
	"github.com/grengojbo/k3ctl/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"

	"github.com/grengojbo/k3ctl/controllers"
	"github.com/grengojbo/k3ctl/pkg/config"
	"github.com/grengojbo/pulumi-modules/automation"
)

var ConfigFile string
var cfgViper = viper.New()
var ppViper = viper.New()

var clusterName string
var cmdFlags types.CmdFlags

const GetScript = "curl -sfL https://get.k3s.io"

// NewCmdCluster returns a new cobra command
func NewCmdApply() *cobra.Command {

	// create new cobra command
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Apply a configuration",
		Long:  `Apply a configuration to a resource in cluster`,
		PreRunE: func(cmd *cobra.Command, args []string) error {

			cmdFlags.DryRun = viper.GetBool("dry-run")
			cmdFlags.DebugLogging = viper.GetBool("verbose")
			cmdFlags.TraceLogging = viper.GetBool("trace")

			// --cluster
			clusterName, err := cmd.Flags().GetString("cluster")
			if err != nil {
				log.Fatalln(err)
			}
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
			 log.Tracef("Simple Config:\n%s", cfg)

			if len(c.Cluster.Spec.Nodes) == 0 {
				log.Fatalln("Is Not Nodes to install k3s cluster")
			}

			// download pulumi plugins
			automation.EnsurePlugins(&c.Plugins)
			// обновляем статус нод
			c.LoadNodeStatus()
			// log.Warnf("kubeconfig:\n%v", c.Kubeconfig)
		},
	}

	// add subcommands
	// cmd.AddCommand(NewCmdClusterCreate())
	// cmd.AddCommand(NewCmdClusterStart())
	// cmd.AddCommand(NewCmdClusterStop())
	// cmd.AddCommand(NewCmdClusterDelete())
	// cmd.AddCommand(NewCmdClusterList())

	// add flags
	cmd.Flags().StringP("cluster", "c", types.DefaultClusterName, "Cluster URL or k3s cluster name to connect to.")

	// done
	return cmd
}
