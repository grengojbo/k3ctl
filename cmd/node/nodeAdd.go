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

	// "github.com/spf13/viper"

	"github.com/grengojbo/k3ctl/controllers"
	"github.com/grengojbo/k3ctl/pkg/config"
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

// NewCmdNodeCreate returns a new cobra command
func NewCmdNodeAdd() *cobra.Command {

	// create new command
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new k3s node",
		Long:  `Add a new containerized k3s node.`,
		Args:  cobra.ExactArgs(1), // exactly one name accepted // TODO: if not specified, inherit from cluster that the node shall belong to, if that is specified
		PreRunE: func(cmd *cobra.Command, args []string) error {

			cmdFlags.DryRun = viper.GetBool("dry-run")
			cmdFlags.DebugLogging = viper.GetBool("verbose")
			cmdFlags.TraceLogging = viper.GetBool("trace")

			NodeName = args[0]
			// --cluster
			clusterName, err := cmd.Flags().GetString("cluster")
			if err != nil {
				log.Fatalln(err)
			}
			ConfigFile = config.InitConfig(clusterName, CfgViper, PpViper)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			// isAddNode := false
			/*************************
			 * Compute Configuration *
			 *************************/
			c, err := controllers.NewClusterFromConfig(CfgViper, cmdFlags)
			if err != nil {
				log.Fatalln(err)
			}

			cfg, _ := yaml.Marshal(c.Cluster)
			log.Tracef("Simple Config:\n%s", cfg)

			if len(c.Cluster.Spec.Nodes) == 0 {
				log.Fatalln("Is Not Nodes to install k3s cluster")
			}

			// if ok := c.AddNode(NodeName); ok {
			// 	isAddNode = true
			// }
			// if !isAddNode {
			// 	log.Errorf("Is NOT set node: %v", NodeName)
			// } else {
			// 	log.Infof("Successfully added %s node(s)!", NodeName)
			// }
		},
	}

	// // add flags
	// cmd.Flags().Int("replicas", 1, "Number of replicas of this node specification.")
	// cmd.Flags().String("role", string(k3d.AgentRole), "Specify node role [server, agent]")
	// if err := cmd.RegisterFlagCompletionFunc("role", util.ValidArgsNodeRoles); err != nil {
	// 	l.Log().Fatalln("Failed to register flag completion for '--role'", err)
	// }
	cmd.Flags().StringP("cluster", "c", types.DefaultClusterName, "Cluster URL or k3s cluster name to connect to.")
	// if err := cmd.RegisterFlagCompletionFunc("cluster", util.ValidArgsAvailableClusters); err != nil {
	// 	log.Fatalln("Failed to register flag completion for '--cluster'", err)
	// }

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

// // parseCreateNodeCmd parses the command input into variables required to create a node
// func parseCreateNodeCmd(cmd *cobra.Command, args []string) ([]*k3d.Node, string) {

// 	// --replicas
// 	replicas, err := cmd.Flags().GetInt("replicas")
// 	if err != nil {
// 		l.Log().Errorln("No replica count specified")
// 		l.Log().Fatalln(err)
// 	}

// 	// --role
// 	roleStr, err := cmd.Flags().GetString("role")
// 	if err != nil {
// 		l.Log().Errorln("No node role specified")
// 		l.Log().Fatalln(err)
// 	}
// 	if _, ok := k3d.NodeRoles[roleStr]; !ok {
// 		l.Log().Fatalf("Unknown node role '%s'\n", roleStr)
// 	}
// 	role := k3d.NodeRoles[roleStr]

// 	// --image
// 	image, err := cmd.Flags().GetString("image")
// 	if err != nil {
// 		l.Log().Errorln("No image specified")
// 		l.Log().Fatalln(err)
// 	}

// 	// --cluster
// 	clusterName, err := cmd.Flags().GetString("cluster")
// 	if err != nil {
// 		l.Log().Fatalln(err)
// 	}

// 	// --memory
// 	memory, err := cmd.Flags().GetString("memory")
// 	if err != nil {
// 		l.Log().Errorln("No memory specified")
// 		l.Log().Fatalln(err)
// 	}
// 	if _, err := dockerunits.RAMInBytes(memory); memory != "" && err != nil {
// 		l.Log().Errorf("Provided memory limit value is invalid")
// 	}

// 	// --runtime-label
// 	runtimeLabelsFlag, err := cmd.Flags().GetStringSlice("runtime-label")
// 	if err != nil {
// 		l.Log().Errorln("No runtime-label specified")
// 		l.Log().Fatalln(err)
// 	}

// 	runtimeLabels := make(map[string]string, len(runtimeLabelsFlag)+1)
// 	for _, label := range runtimeLabelsFlag {
// 		labelSplitted := strings.Split(label, "=")
// 		if len(labelSplitted) != 2 {
// 			l.Log().Fatalf("unknown runtime-label format format: %s, use format \"foo=bar\"", label)
// 		}
// 		cliutil.ValidateRuntimeLabelKey(labelSplitted[0])
// 		runtimeLabels[labelSplitted[0]] = labelSplitted[1]
// 	}

// 	// Internal k3d runtime labels take precedence over user-defined labels
// 	runtimeLabels[k3d.LabelRole] = roleStr

// 	// --k3s-node-label
// 	k3sNodeLabelsFlag, err := cmd.Flags().GetStringSlice("k3s-node-label")
// 	if err != nil {
// 		l.Log().Errorln("No k3s-node-label specified")
// 		l.Log().Fatalln(err)
// 	}

// 	k3sNodeLabels := make(map[string]string, len(k3sNodeLabelsFlag))
// 	for _, label := range k3sNodeLabelsFlag {
// 		labelSplitted := strings.Split(label, "=")
// 		if len(labelSplitted) != 2 {
// 			l.Log().Fatalf("unknown k3s-node-label format format: %s, use format \"foo=bar\"", label)
// 		}
// 		k3sNodeLabels[labelSplitted[0]] = labelSplitted[1]
// 	}

// 	// --network
// 	networks, err := cmd.Flags().GetStringSlice("network")
// 	if err != nil {
// 		l.Log().Fatalf("failed to get --network string slice flag: %v", err)
// 	}

// 	// generate list of nodes
// 	nodes := []*k3d.Node{}
// 	for i := 0; i < replicas; i++ {
// 		node := &k3d.Node{
// 			Name:          fmt.Sprintf("%s-%s-%d", k3d.DefaultObjectNamePrefix, args[0], i),
// 			Role:          role,
// 			Image:         image,
// 			K3sNodeLabels: k3sNodeLabels,
// 			RuntimeLabels: runtimeLabels,
// 			Restart:       true,
// 			Memory:        memory,
// 			Networks:      networks,
// 		}
// 		nodes = append(nodes, node)
// 	}

// 	return nodes, clusterName
// }
