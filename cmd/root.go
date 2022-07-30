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
package cmd

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	// "context"
	// "io/ioutil"
	// "strings"
	// rt "runtime"
	"github.com/grengojbo/k3ctl/cmd/cluster"
	// cfg "github.com/rancher/k3d/v4/cmd/config"
	// "github.com/rancher/k3d/v4/cmd/image"

	app "github.com/grengojbo/k3ctl/cmd/apply"
	"github.com/grengojbo/k3ctl/cmd/kubeconfig"
	"github.com/grengojbo/k3ctl/cmd/node"

	// "github.com/rancher/k3d/v4/cmd/registry"
	cliutil "github.com/grengojbo/k3ctl/cmd/util"
	// "github.com/rancher/k3d/v4/pkg/runtimes"

	log "github.com/sirupsen/logrus"
	"github.com/sirupsen/logrus/hooks/writer"
)

// RootFlags describes a struct that holds flags that can be set on root level of the command
type RootFlags struct {
	// debugLogging       bool
	// traceLogging       bool
	timestampedLogging bool
	version            bool
}

var flags = RootFlags{}
var version string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "k3ctl",
	Short: "Run k3s cluster",
	Long: `
k3ctl is a wrapper CLI that helps you to easily create k3s clusters.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Println("Start...")
		if flags.version {
			printVersion()
		} else {
			if err := cmd.Usage(); err != nil {
				log.Fatalln(err)
			}
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(ver string) {
	version = ver
	if len(os.Args) > 1 {
		parts := os.Args[1:]
		// Check if it's a built-in command, else try to execute it as a plugin
		if _, _, err := rootCmd.Find(parts); err != nil {
			pluginFound, err := cliutil.HandlePlugin(context.Background(), parts)
			if err != nil {
				log.Errorf("Failed to execute plugin '%+v'", parts)
				log.Fatalln(err)
			} else if pluginFound {
				os.Exit(0)
			}
		}
	}
	if err := rootCmd.Execute(); err != nil {
		log.Fatalln(err)
	}
}

func init() {

	rootCmd.PersistentFlags().String("kubeconfig", "", "Local path for your kubeconfig file")
	_ = viper.BindPFlag("kubeconfig", rootCmd.PersistentFlags().Lookup("kubeconfig"))
	// viper.AutomaticEnv()
	// _ = viper.BindEnv("kubeconfig")
	rootCmd.PersistentFlags().Bool("verbose", false, "Enable verbose output (debug logging)")
	rootCmd.PersistentFlags().Bool("trace", false, "Enable super verbose output (trace logging)")
	// rootCmd.PersistentFlags().BoolVar(&flags.debugLogging, "verbose", false, "Enable verbose output (debug logging)")
	// rootCmd.PersistentFlags().BoolVar(&flags.traceLogging, "trace", false, "Enable super verbose output (trace logging)")
	rootCmd.PersistentFlags().BoolVar(&flags.timestampedLogging, "timestamps", false, "Enable Log timestamps")
	rootCmd.PersistentFlags().Bool("dry-run", false, "Show run command and skip the k3s installer")
	_ = viper.BindPFlag("dry-run", rootCmd.PersistentFlags().Lookup("dry-run"))
	_ = viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
	_ = viper.BindPFlag("trace", rootCmd.PersistentFlags().Lookup("trace"))

	// add local flags
	rootCmd.Flags().BoolVar(&flags.version, "version", false, "Show k3ctl and default k3s version")

	// add subcommands
	rootCmd.AddCommand(NewCmdCompletion())
	rootCmd.AddCommand(cluster.NewCmdCluster())
	rootCmd.AddCommand(kubeconfig.NewCmdKubeconfig())
	rootCmd.AddCommand(node.NewCmdNode())
	rootCmd.AddCommand(app.NewCmdApply())
	// rootCmd.AddCommand(image.NewCmdImage())
	// rootCmd.AddCommand(cfg.NewCmdConfig())
	// rootCmd.AddCommand(registry.NewCmdRegistry())

	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Show k3ctl and default k3s version",
		Long:  "Show k3ctl and default k3s version",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	})

	// Init
	cobra.OnInitialize(initLogging)
	// cobra.OnInitialize(initLogging, initRuntime)
}

// initLogging initializes the logger
func initLogging() {
	if viper.GetBool("trace") {
		log.SetLevel(log.TraceLevel)
	} else if viper.GetBool("verbose") {
		log.SetLevel(log.DebugLevel)
	} else {
		switch logLevel := strings.ToUpper(os.Getenv("LOG_LEVEL")); logLevel {
		case "TRACE":
			log.SetLevel(log.TraceLevel)
		case "DEBUG":
			log.SetLevel(log.DebugLevel)
		case "WARN":
			log.SetLevel(log.WarnLevel)
		case "ERROR":
			log.SetLevel(log.ErrorLevel)
		default:
			log.SetLevel(log.InfoLevel)
		}
	}
	log.SetOutput(ioutil.Discard)
	log.AddHook(&writer.Hook{
		Writer: os.Stderr,
		LogLevels: []log.Level{
			log.PanicLevel,
			log.FatalLevel,
			log.ErrorLevel,
			log.WarnLevel,
		},
	})
	log.AddHook(&writer.Hook{
		Writer: os.Stdout,
		LogLevels: []log.Level{
			log.InfoLevel,
			log.DebugLevel,
			log.TraceLevel,
		},
	})

	formatter := &log.TextFormatter{
		ForceColors: true,
	}

	if flags.timestampedLogging || os.Getenv("LOG_TIMESTAMPS") != "" {
		formatter.FullTimestamp = true
	}

	log.SetFormatter(formatter)

}

// func initRuntime() {
// 	runtime, err := runtimes.GetRuntime("docker")
// 	if err != nil {
// 		log.Fatalln(err)
// 	}
// 	runtimes.SelectedRuntime = runtime
// 	log.Debugf("Selected runtime is '%T' on GOOS '%s/%s'", runtimes.SelectedRuntime, rt.GOOS, rt.GOARCH)
// }

func printVersion() {
	fmt.Printf("k3ctl version %s\n", version)
	// fmt.Printf("k3s version %s (default)\n", version.K3sVersion)
}

func generateFishCompletion(writer io.Writer) error {
	return rootCmd.GenFishCompletion(writer, true)
}

// Completion
var completionFunctions = map[string]func(io.Writer) error{
	"bash": rootCmd.GenBashCompletion,
	"zsh": func(writer io.Writer) error {
		if err := rootCmd.GenZshCompletion(writer); err != nil {
			return err
		}

		fmt.Fprintf(writer, "\n# source completion file\ncompdef _k3d k3d\n")

		return nil
	},
	"psh":        rootCmd.GenPowerShellCompletion,
	"powershell": rootCmd.GenPowerShellCompletion,
	"fish":       generateFishCompletion,
}

// NewCmdCompletion creates a new completion command
func NewCmdCompletion() *cobra.Command {
	// create new cobra command
	cmd := &cobra.Command{
		Use:   "completion SHELL",
		Short: "Generate completion scripts for [bash, zsh, fish, powershell | psh]",
		Long: fmt.Sprintf(`To load completions:

Bash:

  $ source <(%[1]s completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ %[1]s completion bash > /etc/bash_completion.d/%[1]s
  # macOS:
  $ %[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ %[1]s completion zsh > "${fpath[1]}/_%[1]s"

  # You will need to start a new shell for this setup to take effect.

fish:

  $ %[1]s completion fish | source

  # To load completions for each session, execute once:
  $ %[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

PowerShell:

  PS> %[1]s completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> %[1]s completion powershell > %[1]s.ps1
  # and source this file from your PowerShell profile.
`, rootCmd.Name()),
		ValidArgs: []string{"bash", "zsh", "fish", "psh", "powershell"},
		Args:      cobra.ExactValidArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			if completionFunc, ok := completionFunctions[args[0]]; ok {
				if err := completionFunc(os.Stdout); err != nil {
					log.Fatalf("Failed to generate completion script for shell '%s'", args[0])
				}
				return
			}
			log.Fatalf("Shell '%s' not supported for completion", args[0])
		},
	}
	return cmd
}
