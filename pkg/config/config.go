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
package config

import (
	"fmt"
	"strings"

	"github.com/grengojbo/k3ctl/pkg/util"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"

	"github.com/spf13/viper"

	k3sv1alpha1 "github.com/grengojbo/k3ctl/api/v1alpha1"
	// conf "github.com/grengojbo/api/v1alpha1"
)

func FromViperSimple(config *viper.Viper) (k3sv1alpha1.Cluster, error) {

	var cfg k3sv1alpha1.Cluster

	// determine config kind
	if config.GetString("kind") != "" && strings.ToLower(config.GetString("kind")) != "cluster" {
		return cfg, fmt.Errorf("Wrong `kind` '%s' != 'Cluster' in config file", config.GetString("kind"))
	}

	if err := config.Unmarshal(&cfg); err != nil {
		log.Errorln("Failed to unmarshal File config")

		return cfg, err
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
	return cfg, nil
}

// var configFile string
// var cfgViper = viper.New()
// var ppViper = viper.New()
// var dryRun bool

func InitConfig(clusterName string, cfgViper *viper.Viper, ppViper *viper.Viper) (configFile string) {

	// dryRun = viper.GetBool("dry-run")
	// Viper for pre-processed config options
	ppViper.SetEnvPrefix("K3S")

	// viper for the general config (file, env and non pre-processed flags)
	cfgViper.SetEnvPrefix("K3S")
	cfgViper.AutomaticEnv()

	cfgViper.SetConfigType("yaml")

	configFile = util.GerConfigFileName(clusterName)
	cfgViper.SetConfigFile(configFile)
	// log.Tracef("Schema: %+v", conf.JSONSchema)

	// if err := config.ValidateSchemaFile(configFile, []byte(conf.JSONSchema)); err != nil {
	// 	log.Fatalf("Schema Validation failed for config file %s: %+v", configFile, err)
	// }

	// try to read config into memory (viper map structure)
	if err := cfgViper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file %s not found: %+v", configFile, err)
		}
		// config file found but some other error happened
		log.Fatalf("Failed to read config file %s: %+v", configFile, err)
	}

	log.Infof("Using config file %s", cfgViper.ConfigFileUsed())
	// }

	// TODO: Default Configs
	cfgViper.SetDefault("spec.kubeconfig.updateDefaultKubeconfig", true)
	cfgViper.SetDefault("spec.kubeconfig.switchCurrentContext", true)

	if log.GetLevel() >= log.DebugLevel {
		c, _ := yaml.Marshal(cfgViper.AllSettings())
		log.Debugf("Configuration:\n%s", c)

		c, _ = yaml.Marshal(ppViper.AllSettings())
		log.Debugf("Additional CLI Configuration:\n%s", c)
	}
	return configFile
}

// func FromViper(config *viper.Viper) (conf.Config, error) {

// 	var cfg conf.Config

// 	// determine config kind
// 	switch strings.ToLower(config.GetString("kind")) {
// 	case "simple":
// 		cfg = conf.SimpleConfig{}
// 	case "cluster":
// 		cfg = conf.ClusterConfig{}
// 	case "clusterlist":
// 		cfg = conf.ClusterListConfig{}
// 	case "":
// 		return nil, fmt.Errorf("Missing `kind` in config file")
// 	default:
// 		return nil, fmt.Errorf("Unknown `kind` '%s' in config file", config.GetString("kind"))
// 	}

// 	if err := config.Unmarshal(&cfg); err != nil {
// 		log.Errorln("Failed to unmarshal File config")

// 		return nil, err
// 	}

// 	return cfg, nil
// }
