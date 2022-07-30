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
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

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

// ListClusterName
func ListClusterName() (clusterNames []string) {
	clusterNames, dirs, _ := ShowFilesInDirectory("./variables", ".yaml")
	for _, dir := range dirs {
		configFileHomeDir := fmt.Sprintf("./variables/%s/cluster.yaml", dir)
		if _, err := os.Stat(configFileHomeDir); err == nil {
			clusterNames = append(clusterNames, dir)
		}
	}
	defPathName, dirs, _ := ShowFilesInDirectory(fmt.Sprintf("~/%s", types.DefaultConfigDirName), ".yaml")
	for _, dir := range dirs {
		configFileHomeDir := fmt.Sprintf("~/%s/%s/cluster.yaml", types.DefaultConfigDirName, dir)
		if _, err := os.Stat(configFileHomeDir); err == nil {
			defPathName = append(defPathName, dir)
		}
	}
	clusterNames = append(clusterNames, defPathName...)
	return clusterNames
}

// ShowFilesInDirectory show files in directory filter extension
func ShowFilesInDirectory(dir string, extension string) ([]string, []string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}
	var filteredFiles []string
	var filteredDirs []string
	for _, file := range files {
		if file.IsDir() {
			filteredDirs = append(filteredDirs, file.Name())
			continue
		}
		if path.Ext(file.Name()) == extension {
			filteredFiles = append(filteredFiles, strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())))
		}
	}
	return filteredFiles, filteredDirs, nil
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

// CheckExitValueFile
func CheckExitValueFile(clusterName string, addonsName string) (file string, err error) {
	file = fmt.Sprintf("./variables/%s/%s-value.yml", clusterName, addonsName)
	if _, err = os.Stat(file); err != nil {
		file = fmt.Sprintf("./variables/%s/%s-value.yaml", clusterName, addonsName)
		if _, err = os.Stat(file); err != nil {
			return file, fmt.Errorf("IS NOT file: %s", file)
		}
	}
	return file, nil
}

// CheckExitFile
func CheckExitFile(file string) (err error) {
	_, err = os.Stat(file)
	return err
}

// CheckСredentials - проверяем естьли файл с credentials
func CheckСredentials(clusterName string, provider string) (ok bool, secretFile string) {
	secretFile = fmt.Sprintf("./variables/%s/secret-%s.ini", clusterName, provider)
	if err := CheckExitFile(secretFile); err != nil {
		return false, secretFile
	}
	return true, secretFile
}

// LoadTemplate - load template from url
func LoadTemplate(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
	// byte to string
	// return string(byte), nil
}

// embed template
// func embedTemplate(template string, data interface{}) (string, error) {
// 	t, err := template.New("").Parse(template)
// 	if err != nil {
// 		return "", err
// 	}
// 	var buf bytes.Buffer
// 	if err := t.Execute(&buf, data); err != nil {
// 		return "", err
// 	}
// 	return buf.String(), nil
// }
