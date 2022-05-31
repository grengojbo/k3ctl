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
package util

import (
	"fmt"
	"os"
	"path"

	"github.com/grengojbo/k3ctl/pkg/types"
	homedir "github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	// "github.com/spf13/afero"
)

// GetConfigDirOrCreate will return the base path of the k3d config directory or create it if it doesn't exist yet
// k3d's config directory will be $HOME/.k3d (Unix)
func GetConfigDirOrCreate() (string, error) {

	// build the path
	homeDir, err := homedir.Dir()
	if err != nil {
		log.Errorln("Failed to get user's home directory")
		return "", err
	}
	configDir := path.Join(homeDir, ".k3d")

	// create directories if necessary
	if err := createDirIfNotExists(configDir); err != nil {
		log.Errorf("Failed to create config path '%s'", configDir)
		return "", err
	}

	return configDir, nil

}

// createDirIfNotExists checks for the existence of a directory and creates it along with all required parents if not.
// It returns an error if the directory (or parents) couldn't be created and nil if it worked fine or if the path already exists.
func createDirIfNotExists(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return os.MkdirAll(path, os.ModePerm)
	}
	return nil
}

// GetEnvDir
func GetEnvDir(clusterName string) (envPath string) {
	envFile := ".env"
	// file, err := afero.ReadFile(v.fs, filename)
	// if err != nil {
	// 	return err
	// }
	// if _, err := afero.Exists(configFile); err != nil {
	// 	log.Errorf("")
	// }
	envFileVariablesDir := fmt.Sprintf("./variables/%s/%s", clusterName, envFile)
	if _, err := os.Stat(envFileVariablesDir); err != nil {
		envFileVariables := fmt.Sprintf("./variables/%s", envFile)
		if _, err := os.Stat(envFileVariables); err != nil {
			if _, err := os.Stat(envFile); err != nil {
				envFileHomeDir := fmt.Sprintf("~/%s/%s", clusterName, envFile)
				if _, err := os.Stat(envFileHomeDir); err != nil {
					envFileDefaultDir := fmt.Sprintf("~/%s/%s/%s", types.DefaultConfigDirName, clusterName, envFile)
					if _, err := os.Stat(envFileDefaultDir); err != nil {
						envFileDefaultFile := fmt.Sprintf("~/%s/%s", types.DefaultConfigDirName, envFile)
						if _, err := os.Stat(envFileDefaultFile); err == nil {
							return envFileDefaultFile
						}
					} else {
						return envFileDefaultDir
					}
				} else {
					return envFileHomeDir
				}
				// // log.Fatalf("%+v", err)
				// log.Fatalln(messageError)
			} else {
				return envFile
			}
		} else {
			return envFileVariables
		}
	} else {
		return envFileVariablesDir
	}
	return ""
}

// GetConfigFileName load config from file
// seach path <clusterName>.yaml, ./variables/<clusterName>.yaml, ~/<clusterName>/cluster.yaml, ~/.k3s/<clusterName>.yaml, ~/.k3s/<clusterName>/cluster.yaml
func GetConfigFileName(configFile string) (configFilePath string) {
	messageError := "Is NOT cluster config file:"
	if configFile == "sample" {
		return "config/samples/k3s_v1alpha1_cluster.yaml"
	}
	// file, err := afero.ReadFile(v.fs, filename)
	// if err != nil {
	// 	return err
	// }
	// if _, err := afero.Exists(configFile); err != nil {
	// 	log.Errorf("")
	// }
	if _, err := os.Stat(configFile); err != nil {
		messageError = fmt.Sprintf("%s %s", messageError, configFile)
		configFileCurrentDir := fmt.Sprintf("./variables/%s.yaml", configFile)
		if _, err := os.Stat(configFileCurrentDir); err != nil {
			messageError = fmt.Sprintf("%s, %s", messageError, configFileCurrentDir)
			configFileHomeDir := fmt.Sprintf("~/%s/cluster.yaml", configFile)
			if _, err := os.Stat(configFileHomeDir); err != nil {
				messageError = fmt.Sprintf("%s, %s", messageError, configFileHomeDir)
				configFileDefaultFile := fmt.Sprintf("~/%s/%s.yaml", types.DefaultConfigDirName, configFile)
				if _, err := os.Stat(configFileDefaultFile); err != nil {
					messageError = fmt.Sprintf("%s, %s", messageError, configFileDefaultFile)
					configFileDefaultDir := fmt.Sprintf("~/%s/%s/cluster.yaml", types.DefaultConfigDirName, configFile)
					if _, err := os.Stat(configFileDefaultDir); err != nil {
						messageError = fmt.Sprintf("%s, %s", messageError, configFileDefaultDir)
					} else {
						return configFileDefaultDir
					}
				} else {
					return configFileDefaultFile
				}
			} else {
				return configFileHomeDir
			}
		} else {
			return configFileCurrentDir
		}
		// log.Fatalf("%+v", err)
		log.Fatalln(messageError)
	}
	return configFile
}
