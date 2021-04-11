/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"time"

	"github.com/docker/go-connections/nat"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DefaultConfigTpl for printing
const DefaultConfigTpl = `---
apiVersion: k3d.io/v1alpha2
kind: Simple
name: %s
servers: 1
agents: 0
image: %s
`

// DefaultConfig templated DefaultConfigTpl
// var DefaultConfig = fmt.Sprintf(
// 	DefaultConfigTpl,
// 	k3d.DefaultClusterName,
// 	fmt.Sprintf("%s:%s", k3d.DefaultK3sImageRepo, version.GetK3sVersion(false)),
// )

// Role defines a k3s node role
type Role string

type LabelWithNodeFilters struct {
	Label       string   `mapstructure:"label" yaml:"label" json:"label,omitempty"`
	NodeFilters []string `mapstructure:"nodeFilters" yaml:"nodeFilters" json:"nodeFilters,omitempty"`
}

type EnvVarWithNodeFilters struct {
	EnvVar      string   `mapstructure:"envVar" yaml:"envVar" json:"envVar,omitempty"`
	NodeFilters []string `mapstructure:"nodeFilters" yaml:"nodeFilters" json:"nodeFilters,omitempty"`
}

type Registry struct {
	Use    []string `mapstructure:"use" yaml:"use,omitempty" json:"use,omitempty"`
	Create bool     `mapstructure:"create" yaml:"create,omitempty" json:"create,omitempty"`
	Config string   `mapstructure:"config" yaml:"config,omitempty" json:"config,omitempty"` // registries.yaml (k3s config for containerd registry override)
}

// Node describes a k3d node
type Node struct {
	Name       string            `yaml:"name" json:"name,omitempty"`
	Role       Role              `yaml:"role" json:"role,omitempty"`
	Image      string            `yaml:"image" json:"image,omitempty"`
	Volumes    []string          `yaml:"volumes" json:"volumes,omitempty"`
	Env        []string          `yaml:"env" json:"env,omitempty"`
	Cmd        []string          // filled automatically based on role
	Args       []string          `yaml:"extraArgs" json:"extraArgs,omitempty"`
	Ports      nat.PortMap       `yaml:"portMappings" json:"portMappings,omitempty"`
	Restart    bool              `yaml:"restart" json:"restart,omitempty"`
	Created    string            `yaml:"created" json:"created,omitempty"`
	Labels     map[string]string // filled automatically
	Networks   []string          // filled automatically
	ExtraHosts []string          // filled automatically
	ServerOpts ServerOpts        `yaml:"serverOpts" json:"serverOpts,omitempty"`
	AgentOpts  AgentOpts         `yaml:"agentOpts" json:"agentOpts,omitempty"`
	GPURequest string            // filled automatically
	Memory     string            // filled automatically
	State      NodeState         // filled automatically
}

// ServerOpts describes some additional server role specific opts
type ServerOpts struct {
	IsInit  bool          `yaml:"isInitializingServer" json:"isInitializingServer,omitempty"`
	KubeAPI *ExposureOpts `yaml:"kubeAPI" json:"kubeAPI"`
}

// ExposureOpts describes settings that the user can set for accessing the Kubernetes API
type ExposureOpts struct {
	nat.PortMapping        // filled automatically (reference to normal portmapping)
	Host            string `yaml:"host,omitempty" json:"host,omitempty"`
}

// ExternalDatastore describes an external datastore used for HA/multi-server clusters
type ExternalDatastore struct {
	Endpoint string `yaml:"endpoint" json:"endpoint,omitempty"`
	CAFile   string `yaml:"caFile" json:"caFile,omitempty"`
	CertFile string `yaml:"certFile" json:"certFile,omitempty"`
	KeyFile  string `yaml:"keyFile" json:"keyFile,omitempty"`
	Network  string `yaml:"network" json:"network,omitempty"`
}

// AgentOpts describes some additional agent role specific opts
type AgentOpts struct{}

// NodeState describes the current state of a node
type NodeState struct {
	Running bool
	Status  string
	Started string
}

// K3sOptions k3s options for generate config
type SimpleConfigOptionsK3s struct {
	ExtraServerArgs []string `mapstructure:"extraServerArgs" yaml:"extraServerArgs"`
	ExtraAgentArgs  []string `mapstructure:"extraAgentArgs" yaml:"extraAgentArgs"`
}

// SimpleConfigOptionsKubeconfig describes the set of options referring to the kubeconfig during cluster creation.
type SimpleConfigOptionsKubeconfig struct {
	UpdateDefaultKubeconfig bool `mapstructure:"updateDefaultKubeconfig" yaml:"updateDefaultKubeconfig" json:"updateDefaultKubeconfig,omitempty"` // default: true
	SwitchCurrentContext    bool `mapstructure:"switchCurrentContext" yaml:"switchCurrentContext" json:"switchCurrentContext,omitempty"`          //nolint:lll    // default: true
}

type Options struct {
	Wait                       bool          `mapstructure:"wait" yaml:"wait"`
	Timeout                    time.Duration `mapstructure:"timeout" yaml:"timeout"`
	DisableLoadbalancer        bool          `mapstructure:"disableLoadbalancer" yaml:"disableLoadbalancer"`
	DisableImageVolume         bool          `mapstructure:"disableImageVolume" yaml:"disableImageVolume"`
	NoRollback                 bool          `mapstructure:"disableRollback" yaml:"disableRollback"`
	PrepDisableHostIPInjection bool          `mapstructure:"disableHostIPInjection" yaml:"disableHostIPInjection"`
	// NodeHookActions            []k3d.NodeHookAction `mapstructure:"nodeHookActions" yaml:"nodeHookActions,omitempty"`
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Region            string                        `mapstructure:"region" yaml:"region" json:"region,omitempty"`
	Servers           int                           `mapstructure:"servers" yaml:"servers" json:"servers,omitempty"`         //nolint:lll    // default 1
	Agents            int                           `mapstructure:"agents" yaml:"agents" json:"agents,omitempty"`            //nolint:lll    // default 0
	ClusterToken      string                        `mapstructure:"token" yaml:"clusterToken" json:"clusterToken,omitempty"` // default: auto-generated
	Nodes             []*Node                       `mapstructure:"token" yaml:"nodes" json:"nodes,omitempty"`
	Host              string                        `mapstructure:"host" yaml:"host,omitempty" json:"host,omitempty"`
	HostIP            string                        `mapstructure:"hostIP" yaml:"hostIP,omitempty" json:"hostIP,omitempty"`
	Labels            []LabelWithNodeFilters        `mapstructure:"labels" yaml:"labels" json:"labels,omitempty"`
	Env               []EnvVarWithNodeFilters       `mapstructure:"env" yaml:"env" json:"env,omitempty"`
	Registries        Registry                      `mapstructure:"registries" yaml:"registries,omitempty" json:"registries,omitempty"`
	Options           Options                       `mapstructure:"options" yaml:"options"`
	K3sOptions        SimpleConfigOptionsK3s        `mapstructure:"k3s" yaml:"k3s"`
	KubeconfigOptions SimpleConfigOptionsKubeconfig `mapstructure:"kubeconfig" yaml:"kubeconfig"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

// +kubebuilder:object:root=true

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
