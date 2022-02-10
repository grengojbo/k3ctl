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
	// l "github.com/rancher/k3d/v5/pkg/logger"
	"github.com/grengojbo/k3ctl/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// var ClusterName string
var NodeName string
var ConfigFile string
var CfgViper = viper.New()
var PpViper = viper.New()
// var DryRun bool

var clusterName string
var cmdFlags types.CmdFlags

// NewCmdNode returns a new cobra command
func NewCmdNode() *cobra.Command {

	// create new cobra command
	cmd := &cobra.Command{
		Use:   "node",
		Short: "Manage node(s)",
		Long:  `Manage node(s)`,
		Run: func(cmd *cobra.Command, args []string) {
			if err := cmd.Help(); err != nil {
				// l.Log().Errorln("Couldn't get help text")
				// l.Log().Fatalln(err)
				log.Errorln("Couldn't get help text")
				log.Fatalln(err)
			}
		},
	}

	// add subcommands
	// cmd.AddCommand(NewCmdNodeAdd())
	cmd.AddCommand(NewCmdNodeAdd(),
		// 	NewCmdNodeStart(),
		// 	NewCmdNodeStop(),
		// 	NewCmdNodeList(),
		// 	NewCmdNodeEdit(),
		NewCmdNodeDelete())

	// add flags

	// done
	return cmd
}
