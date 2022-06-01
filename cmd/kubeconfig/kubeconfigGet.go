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
package kubeconfig

import (
	"fmt"
	// "os"

	// "github.com/rancher/k3d/v5/pkg/client"
	// "github.com/rancher/k3d/v5/pkg/runtimes"
	// k3d "github.com/rancher/k3d/v5/pkg/types"
	// conf "github.com/grengojbo/k3ctl/api/v1alpha1"
	"github.com/grengojbo/k3ctl/pkg/config"
	log "github.com/sirupsen/logrus"

	// k3s "github.com/grengojbo/k3ctl/pkg/k3s"

	"github.com/grengojbo/k3ctl/controllers"
	// log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var kubeconfig string
// var kubeconfig []byte

type getKubeconfigFlags struct {
	all bool
}

// NewCmdKubeconfigGet returns a new cobra command
func NewCmdKubeconfigGet() *cobra.Command {

	// writeKubeConfigOptions := client.WriteKubeConfigOptions{
	// 	UpdateExisting:       true,
	// 	UpdateCurrentContext: true,
	// 	OverwriteExisting:    true,
	// }

	getKubeconfigFlags := getKubeconfigFlags{}

	// create new command
	cmd := &cobra.Command{
		Use:     "get [CLUSTER [CLUSTER [...]] | --all]",
		Short:   "Print kubeconfig(s) from cluster(s).",
		Long:    `Print kubeconfig(s) from cluster(s).`,
		Aliases: []string{"print", "show"},
		// ValidArgsFunction: util.ValidArgsAvailableClusters,
		Args: func(cmd *cobra.Command, args []string) error {
			if (len(args) < 1 && !getKubeconfigFlags.all) || (len(args) > 0 && getKubeconfigFlags.all) {
				return fmt.Errorf("Need to specify one or more cluster names *or* set `--all` flag")
			}
			return nil
		},
		PreRunE: func(cmd *cobra.Command, args []string) error {
			cmdFlags.DryRun = viper.GetBool("dry-run")
			cmdFlags.DebugLogging = viper.GetBool("verbose")
			cmdFlags.TraceLogging = viper.GetBool("trace")

			clusterName = args[0]
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

			// обновляем статус нод
			c.LoadNodeStatus()

			// if !getKubeconfigFlags.all {
			// 	for _, clusterName := range args {
			// 		log.Infof("Load kubeconfig from %s cluster", clusterName)
			// 		ConfigFile = config.InitConfig(clusterName, cfgViper, ppViper)
			// 		cfg, err := config.FromViperSimple(cfgViper)
			// 		if err != nil {
			// 			log.Fatalln(err)
			// 		}
			// 		masters := []conf.ContrelPlanNodes{}
			// 		servers, _, err := util.GetGroupNodes(cfg.Spec.Nodes)
			// 		if err != nil {
			// 			log.Fatalln(err)
			// 		}
			// 		if len(cfg.Spec.Nodes) == 0 {
			// 			log.Fatalln("Is Not Nodes to install k3s cluster")
			// 		}
			// 		for _, node := range servers {
			// 			if bastion, err := cfg.GetBastion(node.Bastion, node); err != nil {
			// 				log.Fatalln(err.Error())
			// 			} else {
			// 				masters = append(masters, conf.ContrelPlanNodes{
			// 					Bastion: bastion,
			// 					Node:    node,
			// 				})
			// 			}
			// 		}

			// 		kubeconfig, err := k3s.GetKubeconfig(masters, DryRun)
			// 		if err != nil {
			// 			log.Fatalln(err.Error())
			// 		}

			// 		isExternal := true
			// 		apiServerUrl, err := cfg.GetAPIServerUrl(masters, &cfg.Spec.Networking, isExternal)
			// 		if err != nil {
			// 			if !DryRun {
			// 				log.Fatal(err)
			// 			}
			// 			log.Error(err)
			// 		}
			// 		log.Debugf("apiServerUrl: %s", apiServerUrl)

			// 		opts := k3s.WriteKubeConfigOptions{
			// 			OverwriteExisting:    true,
			// 			UpdateCurrentContext: cfg.Spec.KubeconfigOptions.SwitchCurrentContext,
			// 		}
			// 		if !DryRun {
			// 			log.Debugf("source kubeconfig:\n%v", kubeconfig)
			// 			pathKubeConfig, err := k3s.SaveCfg(kubeconfig, apiServerUrl, clusterName, opts)
			// 			if err != nil {
			// 				log.Errorln(err.Error())
			// 			}
			// 			// c, _ := yaml.Marshal(newKubeConfig.Clusters)
			// 			log.Infof("new kubeconfig: %s", pathKubeConfig)
			// 			// log.Warnf(" cfg.Spec.KubeconfigOptions.SwitchCurrentContext: %v",  cfg.Spec.KubeconfigOptions.SwitchCurrentContext)
			// 		}
			// 		// 		retrievedCluster, err := client.ClusterGet(cmd.Context(), runtimes.SelectedRuntime, &k3d.Cluster{Name: clusterName})
			// 		// 		if err != nil {
			// 		// 			l.Log().Fatalln(err)
			// 		// 		}
			// 		// 		clusters = append(clusters, retrievedCluster)
			// 	}
			// } else {
			// 	log.Fatalln("TODO: load kubeconfig from all cluster")
			// 	// 	clusters, err = client.ClusterList(cmd.Context(), runtimes.SelectedRuntime)
			// 	// 	if err != nil {
			// 	// 		l.Log().Fatalln(err)
			// 	// 	}
			// }
			// // k, err := k3s.KubeconfigGetDefaultFile()
			// // if err != nil {
			// // 	log.Fatalln(err.Error())
			// // }
			// // log.Debugf("KubeconfigGetDefaultFile: %v", k)
			// // c, _ := k3s.KubeconfigGetDefaultFile()
			// // log.Warnf("clusters: %v", c.Clusters)

			// // // get kubeconfigs from all clusters
			// // errorGettingKubeconfig := false
			// // for _, c := range clusters {
			// // 	l.Log().Debugf("Getting kubeconfig for cluster '%s'", c.Name)
			// // 	fmt.Println("---") // YAML document separator
			// // 	if _, err := client.KubeconfigGetWrite(cmd.Context(), runtimes.SelectedRuntime, c, "-", &writeKubeConfigOptions); err != nil {
			// // 		l.Log().Errorln(err)
			// // 		errorGettingKubeconfig = true
			// // 	}
			// // }

			// // // return with non-zero exit code, if there was an error for one of the clusters
			// // if errorGettingKubeconfig {
			// // 	os.Exit(1)
			// // }
			// // log.Errorf("TODO: %s", viper.GetString("kubeconfig"))
		},
	}

	// add flags
	cmd.Flags().BoolVarP(&getKubeconfigFlags.all, "all", "a", false, "Output kubeconfigs from all existing clusters")
	cmd.Flags().Bool("kubeconfig-switch-context", true, "Directly switch the default kubeconfig's current-context to the new cluster's context (requires --kubeconfig-update-default)")
	_ = cfgViper.BindPFlag("spec.kubeconfig.switchcurrentcontext", cmd.Flags().Lookup("kubeconfig-switch-context"))

	// done
	return cmd
}
