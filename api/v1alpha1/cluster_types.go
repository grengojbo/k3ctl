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
	// go get sigs.k8s.io/cluster-api
	// clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
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

// ###### From cluster api kubeadmin
// https://github.com/kubernetes-sigs/cluster-api-bootstrap-provider-kubeadm

// Encoding specifies the cloud-init file encoding.
// +kubebuilder:validation:Enum=base64;gzip;gzip+base64
type Encoding string

const (
	// Base64 implies the contents of the file are encoded as base64.
	Base64 Encoding = "base64"
	// Gzip implies the contents of the file are encoded with gzip.
	Gzip Encoding = "gzip"
	// GzipBase64 implies the contents of the file are first base64 encoded and then gzip encoded.
	GzipBase64 Encoding = "gzip+base64"
)

// File defines the input for generating write_files in cloud-init.
type File struct {
	// Path specifies the full path on disk where to store the file.
	Path string `json:"path"`

	// Owner specifies the ownership of the file, e.g. "root:root".
	// +optional
	Owner string `json:"owner,omitempty"`

	// Permissions specifies the permissions to assign to the file, e.g. "0640".
	// +optional
	Permissions string `json:"permissions,omitempty"`

	// Encoding specifies the encoding of the file contents.
	// +optional
	Encoding Encoding `json:"encoding,omitempty"`

	// Content is the actual content of the file.
	Content string `json:"content"`
}

// User defines the input for a generated user in cloud-init.
type User struct {
	// Name specifies the user name
	Name string `json:"name"`

	// Gecos specifies the gecos to use for the user
	// +optional
	Gecos *string `json:"gecos,omitempty"`

	// Groups specifies the additional groups for the user
	// +optional
	Groups *string `json:"groups,omitempty"`

	// HomeDir specifies the home directory to use for the user
	// +optional
	HomeDir *string `json:"homeDir,omitempty"`

	// Inactive specifies whether to mark the user as inactive
	// +optional
	Inactive *bool `json:"inactive,omitempty"`

	// Shell specifies the user's shell
	// +optional
	Shell *string `json:"shell,omitempty"`

	// Passwd specifies a hashed password for the user
	// +optional
	Passwd *string `json:"passwd,omitempty"`

	// PrimaryGroup specifies the primary group for the user
	// +optional
	PrimaryGroup *string `json:"primaryGroup,omitempty"`

	// LockPassword specifies if password login should be disabled
	// +optional
	LockPassword *bool `json:"lockPassword,omitempty"`

	// Sudo specifies a sudo role for the user
	// +optional
	Sudo *string `json:"sudo,omitempty"`

	// SSHAuthorizedKeys specifies a list of ssh authorized keys for the user
	// +optional
	SSHAuthorizedKeys []string `json:"sshAuthorizedKeys,omitempty"`
}

// NTP defines input for generated ntp in cloud-init
type NTP struct {
	// Servers specifies which NTP servers to use
	// +optional
	Servers []string `json:"servers,omitempty"`

	// Enabled specifies whether NTP should be enabled
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

// Networking contains elements describing cluster's networking configuration
type Networking struct {
	// APIServerPort specifies the port the API Server should bind to.
	// Defaults to 6443.
	// +optional
	APIServerPort *int32 `json:"apiServerPort,omitempty"`
	// ServiceSubnet is the subnet used by k8s services.
	// Defaults to the first element of the Cluster object's spec.clusterNetwork.pods.cidrBlocks field, or
	// to "10.96.0.0/12" if that's unset.
	// +optional
	ServiceSubnet string `json:"serviceSubnet,omitempty"`
	// PodSubnet is the subnet used by pods.
	// If unset, the API server will not allocate CIDR ranges for every node.
	// Defaults to the first element of the Cluster object's spec.clusterNetwork.services.cidrBlocks if that is set
	// +optional
	PodSubnet string `json:"podSubnet,omitempty"`
	// DNSDomain is the dns domain used by k8s services. Defaults to "cluster.local".
	// +optional
	DNSDomain string `json:"dnsDomain,omitempty"`
}

// ######### END

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

type VolumeWithNodeFilters struct {
	Volume      string   `mapstructure:"volume" yaml:"volume" json:"volume,omitempty"`
	NodeFilters []string `mapstructure:"nodeFilters" yaml:"nodeFilters" json:"nodeFilters,omitempty"`
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
// https://github.com/kubernetes-sigs/cluster-api-bootstrap-provider-kubeadm/blob/master/kubeadm/v1beta1/types.go
type ClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Region            string                        `mapstructure:"region" yaml:"region" json:"region,omitempty"`
	Servers           int                           `mapstructure:"servers" yaml:"servers" json:"servers,omitempty"`         //nolint:lll    // default 1
	Agents            int                           `mapstructure:"agents" yaml:"agents" json:"agents,omitempty"`            //nolint:lll    // default 0
	ClusterToken      string                        `mapstructure:"token" yaml:"clusterToken" json:"clusterToken,omitempty"` // default: auto-generated
	Nodes             []*Node                       `mapstructure:"nodes" yaml:"nodes" json:"nodes,omitempty"`
	Host              string                        `mapstructure:"host" yaml:"host,omitempty" json:"host,omitempty"`
	HostIP            string                        `mapstructure:"hostIP" yaml:"hostIP,omitempty" json:"hostIP,omitempty"`
	Labels            []LabelWithNodeFilters        `mapstructure:"labels" yaml:"labels" json:"labels,omitempty"`
	Env               []EnvVarWithNodeFilters       `mapstructure:"env" yaml:"env" json:"env,omitempty"`
	Registries        Registry                      `mapstructure:"registries" yaml:"registries,omitempty" json:"registries,omitempty"`
	Options           Options                       `mapstructure:"options" yaml:"options" json:"options,omitempty"`
	K3sOptions        SimpleConfigOptionsK3s        `mapstructure:"k3s" yaml:"k3s" json:"k3s,omitempty"`
	KubeconfigOptions SimpleConfigOptionsKubeconfig `mapstructure:"kubeconfig" yaml:"kubeconfig" json:"kubeconfig,omitempty"`
	Volumes           []VolumeWithNodeFilters       `mapstructure:"volumes" yaml:"volumes" json:"volumes,omitempty"`
	// The cluster name
	// +optional
	ClusterName string `json:"clusterName,omitempty"`
	// KubernetesVersion is the target version of the control plane.
	// NB: This value defaults to the Machine object spec.kuberentesVersion
	// +optional
	KubernetesVersion string `mapstructure:"kubernetesVersion" yaml:"kubernetesVersion" json:"kubernetesVersion,omitempty"`
	// K3sChannel Release channel: stable, latest, or i.e. v1.19
	// +optional
	K3sChannel string `mapstructure:"channel" yaml:"channel" json:"channel,omitempty"`
	// Networking holds configuration for the networking topology of the cluster.
	// NB: This value defaults to the Cluster object spec.clusterNetwork.
	// +optional
	Networking Networking `mapstructure:"networking" yaml:"networking" json:"networking,omitempty"`
	// Cluster network configuration.
	// +optional
	// ClusterNetwork *clusterv1.ClusterNetwork `json:"clusterNetwork,omitempty"`
	// CertificatesDir specifies where to store or look for all required certificates.
	// NB: if not provided, this will default to `/etc/kubernetes/pki`
	// +optional
	CertificatesDir string `json:"certificatesDir,omitempty"`
	// Files specifies extra files to be passed to user_data upon creation.
	// +optional
	Files []File `json:"files,omitempty"`
	// PreCommands specifies extra commands to run before kubeadm runs
	// +optional
	PreCommands []string `json:"preCommands,omitempty"`
	// PostCommands specifies extra commands to run after kubeadm runs
	// +optional
	PostCommands []string `json:"postCommands,omitempty"`
	// Users specifies extra users to add
	// +optional
	Users []User `json:"users,omitempty"`
	// NTP specifies NTP configuration
	// +optional
	NTP *NTP `json:"ntp,omitempty"`
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
