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
	"errors"
	"fmt"
	"time"

	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"

	// log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	// go get sigs.k8s.io/cluster-api

	clusterv1 "sigs.k8s.io/cluster-api/api/v1alpha3"
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
	// Local      string   = "local"
	// Public              string = "ExternalIP"
	// Private             string = "InternalIP"
	ExternalIP          string = "ExternalIP"
	ExternalDNS         string = "ExternalDNS"
	InternalIP          string = "InternalIP"
	InternalDNS         string = "InternalDNS"
	SshKeyDefault       string = "~/.ssh/id_rsa"
	SshPortDefault      int32  = 22
	DatastoreMySql      string = "mysql"
	DatastorePostgreSql string = "postgres"
	DatastoreEtcd       string = "etcd"
)

var PrivateHost = []string{InternalIP, InternalDNS}
var PublicHost = []string{ExternalIP, InternalIP}
var ConnectionHosts = []string{ExternalIP, ExternalDNS, InternalDNS, InternalIP}

// var LocalHost = []string{Local, "localhost", "127.0.0.1"}

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

	// SSHAuthorizedKey specifies a list of ssh authorized keys for the user
	// +optional
	SSHAuthorizedKey string `yaml:"sshAuthorizedKey" json:"sshAuthorizedKey,omitempty"`
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
	// APIServerAddresses is a list of addresses assigned to the API Server.
	// +optional
	APIServerAddresses clusterv1.MachineAddresses `json:"apiServerAddresses,omitempty"`
	// APIServerPort specifies the port the API Server should bind to.
	// Defaults to 6443.
	// +optional
	APIServerPort int32 `json:"apiServerPort,omitempty"`
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
	// ClusterDns	Cluster IP for coredns service. Should be in your service-cidr range --cluster-dns value	“10.43.0.10”
	// +optional
	ClusterDns string `json:"clusterDns,omitempty"`
	// CNI plugins ("flannel", "calico", "cilium", "aws")
	// +optional
	CNI string `json:"cni,omitempty"`
	// Backend The default backend for flannel is VXLAN
	// --flannel-backend value	“vxlan”	One of ‘none’, ‘vxlan’, ‘ipsec’, ‘host-gw’, or ‘wireguard’
	// Calico "ipip", “vxlan”
	// Cilium "ipip", “vxlan”
	// +optional
	Backend string `json:"backend,omitempty"`
}

// ######### END

type Registry struct {
	Name      string   `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Namespace string   `mapstructure:"namespace" yaml:"namespace" json:"namespace,omitempty"`
	Type      string   `mapstructure:"type" yaml:"type" json:"type,omitempty"` // helm, manifest, orchestration
	Use       []string `mapstructure:"use" yaml:"use" json:"use,omitempty"`
	Create    bool     `mapstructure:"create" yaml:"create" json:"create,omitempty"`
	Config    string   `mapstructure:"config" yaml:"config" json:"config,omitempty"` // registries.yaml (k3s config for containerd registry override)
}

type PulumiModule struct {
	Name       string            `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Namespace  string            `mapstructure:"namespace" yaml:"namespace" json:"namespace,omitempty"`
	Disabled   bool              `mapstructure:"disabled" yaml:"disabled" json:"disabled,omitempty"`
	Version    string            `mapstructure:"version" yaml:"version" json:"version,omitempty"`
	URL        string            `mapstructure:"url" yaml:"url" json:"url,omitempty"`
	Values     map[string]string `mapstructure:"values" yaml:"values" json:"values,omitempty"`
	ValuesFile string            `mapstructure:"valuesFile" yaml:"valuesFile" json:"valuesFile,omitempty"`
}

type ExternalDns struct {
	Name       string            `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Namespace  string            `mapstructure:"namespace" yaml:"namespace" json:"namespace,omitempty"`
	Disabled   bool              `mapstructure:"disabled" yaml:"disabled" json:"disabled,omitempty"`
	Version    string            `mapstructure:"version" yaml:"version" json:"version,omitempty"`
	URL        string            `mapstructure:"url" yaml:"url" json:"url,omitempty"`
	Values     map[string]string `mapstructure:"values" yaml:"values" json:"values,omitempty"`
	ValuesFile string            `mapstructure:"valuesFile" yaml:"valuesFile" json:"valuesFile,omitempty"`
	Repo       HelmRepo          `mapstructure:"repo" yaml:"repo" json:"repo,omitempty"`
}

type MetalLB struct {
	Name       string            `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Namespace  string            `mapstructure:"namespace" yaml:"namespace" json:"namespace,omitempty"`
	Disabled   bool              `mapstructure:"disabled" yaml:"disabled" json:"disabled,omitempty"`
	Version    string            `mapstructure:"version" yaml:"version" json:"version,omitempty"`
	URL        string            `mapstructure:"url" yaml:"url" json:"url,omitempty"`
	Values     map[string]string `mapstructure:"values" yaml:"values" json:"values,omitempty"`
	ValuesFile string            `mapstructure:"valuesFile" yaml:"valuesFile" json:"valuesFile,omitempty"`
	Repo       HelmRepo          `mapstructure:"repo" yaml:"repo" json:"repo,omitempty"`
}

type CertManager struct {
	Name       string            `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Namespace  string            `mapstructure:"namespace" yaml:"namespace" json:"namespace,omitempty"`
	Disabled   bool              `mapstructure:"disabled" yaml:"disabled" json:"disabled,omitempty"`
	Version    string            `mapstructure:"version" yaml:"version" json:"version,omitempty"`
	URL        string            `mapstructure:"url" yaml:"url" json:"url,omitempty"`
	Values     map[string]string `mapstructure:"values" yaml:"values" json:"values,omitempty"`
	ValuesFile string            `mapstructure:"valuesFile" yaml:"valuesFile" json:"valuesFile,omitempty"`
	Repo       HelmRepo          `mapstructure:"repo" yaml:"repo" json:"repo,omitempty"`
}

// observer
type Monitoring struct {
	Name       string            `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Namespace  string            `mapstructure:"namespace" yaml:"namespace" json:"namespace,omitempty"`
	Disabled   bool              `mapstructure:"disabled" yaml:"disabled" json:"disabled,omitempty"`
	Version    string            `mapstructure:"version" yaml:"version" json:"version,omitempty"`
	URL        string            `mapstructure:"url" yaml:"url" json:"url,omitempty"`
	Values     map[string]string `mapstructure:"values" yaml:"values" json:"values,omitempty"`
	ValuesFile string            `mapstructure:"valuesFile" yaml:"valuesFile" json:"valuesFile,omitempty"`
	Repo       HelmRepo          `mapstructure:"repo" yaml:"repo" json:"repo,omitempty"`
}

type VeleroStorage struct {
	Name    string `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Image   string `mapstructure:"image" yaml:"image" json:"image,omitempty"`
	Default bool   `mapstructure:"default" yaml:"default" json:"default,omitempty"`
	Region  string `mapstructure:"region" yaml:"region" json:"region,omitempty"`
	Bucket  string `mapstructure:"bucket" yaml:"bucket" json:"bucket,omitempty"`
}

type ResticBackup struct {
	DeployRestic bool `mapstructure:"deployRestic" yaml:"deployRestic" json:"deployRestic,omitempty"`
}

type VeleroBackup struct {
	Providers  []string        `mapstructure:"providers" yaml:"providers" json:"providers,omitempty"` // aws azure s3
	Storages   []VeleroStorage `mapstructure:"storages" yaml:"storages" json:"storages,omitempty"`
	Repo       HelmRepo        `mapstructure:"repo" yaml:"repo" json:"repo,omitempty"`
	Restic     ResticBackup    `mapstructure:"restic" yaml:"restic" json:"restic,omitempty"`
	SecretFile string          `mapstructure:"secretFile" yaml:"secretFile" json:"secretFile,omitempty"`
}

type Backup struct {
	Name       string            `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Namespace  string            `mapstructure:"namespace" yaml:"namespace" json:"namespace,omitempty"`
	Disabled   bool              `mapstructure:"disabled" yaml:"disabled" json:"disabled,omitempty"`
	Version    string            `mapstructure:"version" yaml:"version" json:"version,omitempty"`
	URL        string            `mapstructure:"url" yaml:"url" json:"url,omitempty"`
	Values     map[string]string `mapstructure:"values" yaml:"values" json:"values,omitempty"`
	ValuesFile string            `mapstructure:"valuesFile" yaml:"valuesFile" json:"valuesFile,omitempty"`
	Velero     VeleroBackup      `mapstructure:"velero" yaml:"velero" json:"velero,omitempty"`
	Provider   string            `mapstructure:"provider" yaml:"provider" json:"provider,omitempty"` // aws azure s3
	Bucket     string            `mapstructure:"bucket" yaml:"bucket" json:"bucket,omitempty"`
	Region     string            `mapstructure:"region" yaml:"region" json:"region,omitempty"`
	Repo       HelmRepo          `mapstructure:"repo" yaml:"repo" json:"repo,omitempty"`
	Schedules  []SchedulesBackup `mapstructure:"schedules" yaml:"schedules" json:"schedules,omitempty"`
}
type SchedulesBackup struct {
	Name                       string            `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Disabled                   bool              `mapstructure:"disabled" yaml:"disabled" json:"disabled,omitempty"`
	Schedule                   string            `mapstructure:"schedule" yaml:"schedule" json:"schedule,omitempty"`
	Labels                     map[string]string `mapstructure:"labels" yaml:"labels" json:"labels,omitempty"`
	Annotations                map[string]string `mapstructure:"annotations" yaml:"annotations" json:"annotations,omitempty"`
	UseOwnerReferencesInBackup bool              `mapstructure:"useOwnerReferencesInBackup" yaml:"useOwnerReferencesInBackup" json:"useOwnerReferencesInBackup,omitempty"`
	Ttl                        string            `mapstructure:"ttl" yaml:"ttl" json:"ttl,omitempty"`
	IncludedNamespaces         []string          `mapstructure:"includedNamespaces" yaml:"includedNamespaces" json:"includedNamespaces,omitempty"`
	// #     template:
	// #       includedNamespaces:
	// #       - foo
}

type Ingress struct {
	Name           string            `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	Namespace      string            `mapstructure:"namespace" yaml:"namespace" json:"namespace,omitempty"`
	Disabled       bool              `mapstructure:"disabled" yaml:"disabled" json:"disabled,omitempty"`
	HostMode       bool              `mapstructure:"hostMode" yaml:"hostMode" json:"hostMode,omitempty"`
	Version        string            `mapstructure:"version" yaml:"version" json:"version,omitempty"`
	URL            string            `mapstructure:"url" yaml:"url" json:"url,omitempty"`
	Values         map[string]string `mapstructure:"values" yaml:"values" json:"values,omitempty"`
	ValuesFile     string            `mapstructure:"valuesFile" yaml:"valuesFile" json:"valuesFile,omitempty"`
	DefaultBackend DefaultBackend    `mapstructure:"defaultBackend" yaml:"defaultBackend" json:"defaultBackend,omitempty"`
	Repo           HelmRepo          `mapstructure:"repo" yaml:"repo" json:"repo,omitempty"`
}

type DefaultBackend struct {
	Registry string `mapstructure:"registry" yaml:"registry" json:"registry,omitempty"`
	Image    string `mapstructure:"image" yaml:"image" json:"image,omitempty"`
	Tag      string `mapstructure:"tag" yaml:"tag" json:"tag,omitempty"`
}

type AddonOptions struct {
	UpdateStrategy string `mapstructure:"updateStrategy,omitempty" yaml:"updateStrategy,omitempty" json:"updateStrategy,omitempty"` // none, latest
	NoWait         bool   `mapstructure:"noWait" yaml:"noWait" json:"noWait,omitempty"`
	// helmAddons []string `mapstructure:"helmAddons" yaml:"helmAddons" json:"helmAddons,omitempty"`
}

type Addons struct {
	Ingress      Ingress      `mapstructure:"ingress" yaml:"ingress" json:"ingress,omitempty"`
	CertManager  CertManager  `mapstructure:"certManager" yaml:"certManager,omitempty" json:"certManager,omitempty"`
	MetalLB      MetalLB      `mapstructure:"metallb" yaml:"metallb,omitempty" json:"metallb,omitempty"`
	Monitoring   Monitoring   `mapstructure:"monitoring" yaml:"monitoring,omitempty" json:"monitoring,omitempty"`
	Backup       Backup       `mapstructure:"backup" yaml:"backup,omitempty" json:"backup,omitempty"`
	ExternalDns  ExternalDns  `mapstructure:"externalDns" yaml:"externalDns,omitempty" json:"externalDns,omitempty"`
	PulumiModule PulumiModule `mapstructure:"pulumi" yaml:"pulumi,omitempty" json:"pulumi,omitempty"`
	Registries   Registry     `mapstructure:"registries" yaml:"registries,omitempty" json:"registries,omitempty"`
	Options      AddonOptions `mapstructure:"options" yaml:"options" json:"options,omitempty"`
}

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

type VolumeWithNodeFilters struct {
	Volume      string   `mapstructure:"volume" yaml:"volume" json:"volume,omitempty"`
	NodeFilters []string `mapstructure:"nodeFilters" yaml:"nodeFilters" json:"nodeFilters,omitempty"`
}

// ClusterNode struct for cluster node.
type ClusterNode struct {
	InstanceID              string   `json:"instance-id,omitempty" yaml:"instance-id,omitempty"`
	InstanceStatus          string   `json:"instance-status,omitempty" yaml:"instance-statu,omitempty"`
	ExternalIP              []string `json:"external-ip,omitempty" yaml:"external-ip,omitempty"`
	InternalIP              []string `json:"internal-ip,omitempty" yaml:"internal-ip,omitempty"`
	Roles                   string   `json:"roles,omitempty" yaml:"status,omitempty"`
	Status                  string   `json:"status,omitempty" yaml:"status,omitempty"`
	HostName                string   `json:"hostname,omitempty" yaml:"hostname,omitempty"`
	ContainerRuntimeVersion string   `json:"containerRuntimeVersion,omitempty" yaml:"containerRuntimeVersion,omitempty"`
	Version                 string   `json:"version,omitempty" yaml:"version,omitempty"`
	Master                  bool     `json:"-" yaml:"master,omitempty"`
}

type ContrelPlanNodes struct {
	ClusterName     string       `mapstructure:"clusterName" yaml:"clusterName" json:"clusterName,omitempty"`
	ApiServerAddres string       `mapstructure:"apiServerAddres" yaml:"apiServerAddres" json:"apiServerAddres,omitempty"`
	Bastion         *BastionNode `mapstructure:"bastion" yaml:"bastion" json:"bastion"`
	Node            *Node        `mapstructure:"node" yaml:"node" json:"node"`
}
type BastionNode struct {
	// Name The bastion Name
	Name string `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	// User Ssh user is empty use bastion name
	// +optional
	User string `mapstructure:"user" yaml:"user" json:"user,omitempty"`
	// Address bastion host
	Address string `mapstructure:"address" yaml:"address" json:"address"`
	// SshPort specifies the port the SSH bastion host.
	// Defaults to 22.
	// +optional
	SshPort int32 `mapstructure:"sshPort" yaml:"sshPort" json:"sshPort,omitempty"`
	// SSHAuthorizedKey specifies a list of ssh authorized keys for the user
	// +optional
	SSHAuthorizedKey string `mapstructure:"sshAuthorizedKey" yaml:"sshAuthorizedKey" json:"sshAuthorizedKey,omitempty"`
	// RemoteConnectionString TODO: tranclate строка подключения к удаленному хосту если через bastion
	// +optional
	RemoteConnectionString string `mapstructure:"remoteConnectionString,omitempty" yaml:"remoteConnectionString,omitempty" json:"remoteConnectionString,omitempty"`
	// RemoteSudo TODO: tranclate если через bastion и пользовател на приватном хосте не root устанавливается true
	// +optional
	RemoteSudo string `mapstructure:"remoteSudo,omitempty" yaml:"remoteSudo,omitempty" json:"remoteSudo,omitempty"`
	// RemoteAddress адрес хоста за бастионом
	// TODO: translate
	// +optional
	RemoteAddress string `mapstructure:"remoteAddress,omitempty" yaml:"remoteAddress,omitempty" json:"remoteAddress,omitempty"`
	// +optional
	RemoteUser string `mapstructure:"remoteUser,omitempty" yaml:"remoteUser,omitempty" json:"remoteUser,omitempty"`
	// +optional
	RemotePort int32 `mapstructure:"remotePort,omitempty" yaml:"remotePort,omitempty" json:"remotePort,omitempty"`
}

// Node describes a k3d node
type Node struct {
	Name       string            `yaml:"name" json:"name,omitempty"`
	User       string            `yaml:"user" json:"user,omitempty"`
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
	// Bastion имя ssh bastion сервера если local то запускается на локальном хосте
	// +optional
	Bastion string `yaml:"bastion" json:"bastion,omitempty"`
	// Addresses is a list of addresses assigned to the machine.
	// This field is copied from the infrastructure provider reference.
	// https://github.com/kubernetes-sigs/cluster-api/blob/2cbeb175b243da6953c4edf9e7ec99eac4e2a4a2/api/v1alpha3/common_types.go
	// +optional
	Addresses clusterv1.MachineAddresses `json:"addresses,omitempty"`
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
	// kubeconfig connection from k3ctl to server ExternalIP | InternalIP
	ConnectType             string `mapstructure:"connectType" yaml:"connectType" json:"connectType,omitempty"`
	Patch                   string `mapstructure:"patch" yaml:"patch" json:"patch,omitempty"`
	UpdateDefaultKubeconfig bool   `mapstructure:"updateDefaultKubeconfig" yaml:"updateDefaultKubeconfig" json:"updateDefaultKubeconfig,omitempty"` // default: true
	SwitchCurrentContext    bool   `mapstructure:"switchCurrentContext" yaml:"switchCurrentContext" json:"switchCurrentContext,omitempty"`          //nolint:lll    // default: true
}

type Options struct {
	Protected           bool          `mapstructure:"protected" yaml:"protected" json:"protected,omitempty"`
	Wait                bool          `mapstructure:"wait" yaml:"wait" json:"wait,omitempty"`
	Timeout             time.Duration `mapstructure:"timeout" yaml:"timeout" json:"timeout,omitempty"`
	DisableLoadbalancer bool          `mapstructure:"disableLoadbalancer" yaml:"disableLoadbalancer" json:"disableLoadbalancer,omitempty"`
	// EnableIngress             bool          `mapstructure:"enableIngress" yaml:"enableIngress" json:"enableIngress,omitempty"`
	DisableImageVolume         bool `mapstructure:"disableImageVolume" yaml:"disableImageVolume" json:"disableImageVolume,omitempty"`
	NoRollback                 bool `mapstructure:"disableRollback" yaml:"disableRollback" json:"disableRollback,omitempty"`
	PrepDisableHostIPInjection bool `mapstructure:"disableHostIPInjection" yaml:"disableHostIPInjection" json:"disableHostIPInjection,omitempty"`
	// SELinux To leverage SELinux, specify the --selinux flag when starting K3s servers and agents.
	// https://rancher.com/docs/k3s/latest/en/advanced/
	// +optional
	SELinux bool `mapstructure:"selinux" yaml:"selinux" json:"selinux,omitempty"`
	// Rootless --rootless Running Servers and Agents with Rootless
	// k3s Experimental Options
	// +optional
	Rootless bool `mapstructure:"rootless" yaml:"rootless" json:"rootless,omitempty"`
	// SecretsEncryption --secrets-encryption	Enable Secret encryption at rest
	// +optional
	SecretsEncryption bool `mapstructure:"secretsEncryption" yaml:"secretsEncryption" json:"secretsEncryption,omitempty"`
	// --agent-token value	K3S_AGENT_TOKEN	Shared secret used to join agents to the cluster, but not servers
	// --agent-token-file value	K3S_AGENT_TOKEN_FILE	File containing the agent secret
	// --server value, -s value	K3S_URL	Server to connect to, used to join a cluster
	// --cluster-init	K3S_CLUSTER_INIT	Initialize new cluster master
	// --cluster-reset	K3S_CLUSTER_RESET	Forget all peers and become a single cluster new cluster master
	// NodeHookActions            []k3d.NodeHookAction `mapstructure:"nodeHookActions" yaml:"nodeHookActions,omitempty"`
}

type Datastore struct {
	// Name Database name ("mysql", "postgresql", "etcd")
	Name string `mapstructure:"name" yaml:"name" json:"name,omitempty"`
	// Provider Database name ("mysql", "postgres", "etcd")
	Provider string `mapstructure:"provider" yaml:"provider" json:"provider,omitempty"`
	Username string `mapstructure:"username" yaml:"username" json:"username,omitempty"`
	Password string `mapstructure:"DB_PASSWORD"`
	Host     string `mapstructure:"host" yaml:"host,omitempty" json:"host,omitempty"`
	// Port DataBase port
	// +optional
	Port int32 `mapstructure:"port" yaml:"port,omitempty" json:"port,omitempty"`
	// CertFile K3S_DATASTORE_CERTFILE='/path/to/client.crt'
	// +optional
	CertFile string `mapstructure:"certFile" yaml:"certFile,omitempty" json:"certFile,omitempty"`
	// KeyFile K3S_DATASTORE_KEYFILE='/path/to/client.key'
	// +optional
	KeyFile string `mapstructure:"keyFile" yaml:"keyFile,omitempty" json:"keyFile,omitempty"`
}

type LoadBalancer struct {
	Name string `mapstructure:"name" yaml:"name" json:"name"`
	// IP addresse at which the controller services are available
	ExternalIP string `mapstructure:"externalIPs" yaml:"externalIPs" json:"externalIPs,omitempty"`
	MetalLb    string `mapstructure:"metalLb" yaml:"metalLb" json:"metalLb,omitempty"`
	KubeVip    string `mapstructure:"kubeVip" yaml:"kubeVip" json:"kubeVip,omitempty"`
}

// ClusterSpec defines the desired state of Cluster
// https://github.com/kubernetes-sigs/cluster-api-bootstrap-provider-kubeadm/blob/master/kubeadm/v1beta1/types.go
type ClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Region            string                        `mapstructure:"region" yaml:"region" json:"region,omitempty"`
	Provider          string                        `mapstructure:"provider" json:"provider" yaml:"provider"`
	Operator          bool                          `mapstructure:"operator" yaml:"operator" json:"operator,omitempty"`
	Servers           int                           `mapstructure:"servers" yaml:"servers" json:"servers,omitempty"` //nolint:lll    // default 1
	Agents            int                           `mapstructure:"agents" yaml:"agents" json:"agents,omitempty"`    //nolint:lll    // default 0
	ClusterToken      string                        `mapstructure:"clusterToken" yaml:"clusterToken" json:"clusterToken,omitempty"`
	AgentToken        string                        `mapstructure:"agentToken" yaml:"agentToken" json:"agentToken,omitempty"`
	Bastions          []*BastionNode                `mapstructure:"bastions" yaml:"bastions" json:"bastions,omitempty"`
	Nodes             []*Node                       `mapstructure:"nodes" yaml:"nodes" json:"nodes,omitempty"`
	Labels            []LabelWithNodeFilters        `mapstructure:"labels" yaml:"labels" json:"labels,omitempty"`
	Env               []EnvVarWithNodeFilters       `mapstructure:"env" yaml:"env" json:"env,omitempty"`
	Options           Options                       `mapstructure:"options" yaml:"options" json:"options,omitempty"`
	K3sOptions        SimpleConfigOptionsK3s        `mapstructure:"k3s" yaml:"k3s" json:"k3s,omitempty"`
	LoadBalancer      LoadBalancer                  `mapstructure:"loadBalancer" yaml:"loadBalancer" json:"loadBalancer,omitempty"`
	Addons            Addons                        `mapstructure:"addons" yaml:"addons" json:"addons,omitempty"`
	KubeconfigOptions SimpleConfigOptionsKubeconfig `mapstructure:"kubeconfig" yaml:"kubeconfig" json:"kubeconfig,omitempty"`
	Volumes           []VolumeWithNodeFilters       `mapstructure:"volumes" yaml:"volumes" json:"volumes,omitempty"`
	// Host              string                        `mapstructure:"host" yaml:"host,omitempty" json:"host,omitempty"`
	// HostIP            string                        `mapstructure:"hostIP" yaml:"hostIP,omitempty" json:"hostIP,omitempty"`
	// Datastore k3s datastore to enable HA https://rancher.com/docs/k3s/latest/en/installation/datastore/
	// +optional
	Datastore Datastore `mapstructure:"datastore" yaml:"datastore" json:"datastore,omitempty"`
	// The cluster name
	// +optional
	ClusterName string `json:"clusterName,omitempty"`
	// KubernetesVersion is the target version of the control plane.
	// NB: This value defaults to the Machine object spec.kuberentesVersion
	// +optional
	KubernetesVersion string `mapstructure:"kubernetesVersion" yaml:"kubernetesVersion" json:"kubernetesVersion,omitempty"`
	// K3sChannel Release channel: stable, latest, or i.e. v1.19
	// +optional
	K3sChannel string `mapstructure:"channel,omitempty" yaml:"channel,omitempty" json:"channel,omitempty"`
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
	Status      string  `json:"status,omitempty"`
	MasterNodes []*Node `json:"masters,omitempty"`
	WorkerNodes []*Node `json:"workers,omitempty"`
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

// type K3sExecOptions struct {
// 	Datastore           string
// 	ExtraArgs           []string
// 	FlannelIPSec        bool
// 	NoExtras            bool
// 	LoadBalancer        *LoadBalancer
// 	Networking          *Networking
// 	Options 						*Options
// 	K3sChannel          string
// 	KubernetesVersion 	string
// 	ClusterToken 				string
// 	AgentToken 					string
// 	Ingress             string
// 	// DisableLoadbalancer bool
// 	// DisableIngress      bool
// 	// SELinux             bool
// 	// Rootless            bool

// 	// SecretsEncryption   bool
// }

type K3sIstallOptions struct {
	ExecString   string
	LoadBalancer string
	Ingress      string
	CNI          string
	Backend      string
	K3sVersion   string
	K3sChannel   string
	IsCluster    bool
	Node         *Node
}
type K3sWorkerOptions struct {
	JoinAgentCommand string `json:"joinAgentCommand,omitempty"`
	ExecString       string `json:"execString,omitempty"`
	ApiServerAddres  string `json:"apiServerAddres,omitempty"`
	ApiServerPort    int32  `json:"apiServerPort,omitempty"`
	Token            string `json:"token,omitempty"`
	K3sVersion       string `json:"k3sVersion,omitempty"`
	K3sChannel       string `json:"k3sChannel,omitempty"`
}

type EnvServer struct {
	K3sDatastoreEndpoint string `mapstructure:"K3S_DATASTORE_ENDPOINT"`
	K3sDatastoreXafile   string `mapstructure:"K3S_DATASTORE_CAFILE"`
	K3sDatastoreCertfile string `mapstructure:"K3S_DATASTORE_CERTFILE"`
	K3sDatastoreKeyfile  string `mapstructure:"K3S_DATASTORE_KEYFILE"`
	AwsAccessKeyId       string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey   string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	K3sToken             string `mapstructure:"K3S_TOKEN"`
	K3sTokenFile         string `mapstructure:"K3S_TOKEN_FILE"`
	K3sKubeconfigOutput  string `mapstructure:"K3S_KUBECONFIG_OUTPUT"`
	K3sKubeconfigMode    string `mapstructure:"K3S_KUBECONFIG_MODE"`
	K3sNodeName          string `mapstructure:"K3S_NODE_NAME"`
	K3sSelinux           string `mapstructure:"K3S_SELINUX"`
	K3sLbServerPort      string `mapstructure:"K3S_LB_SERVER_PORT"`
	K3sResolvConf        string `mapstructure:"K3S_RESOLV_CONF"`
	K3sAgentToken        string `mapstructure:"K3S_AGENT_TOKEN"`
	K3sAgentTokenFile    string `mapstructure:"K3S_AGENT_TOKEN_FILE"`
	K3sUrl               string `mapstructure:"K3S_URL"`
	K3sClusterInit       string `mapstructure:"K3S_CLUSTER_INIT"`
	K3sClusterReset      string `mapstructure:"K3S_CLUSTER_RESET"`
	K3sClusterSecret     string `mapstructure:"K3S_CLUSTER_SECRET"`
}

type EnvConfig struct {
	DBPassword          string `mapstructure:"DB_PASSWORD"`
	HcloudToken         string `mapstructure:"HCLOUD_TOKEN"`
	AwsAccessKeyId      string `mapstructure:"AWS_ACCESS_KEY_ID"`
	AwsSecretAccessKey  string `mapstructure:"AWS_SECRET_ACCESS_KEY"`
	AzureClientId       string `mapstructure:"ARM_CLIENT_ID"`
	AzureClientSecret   string `mapstructure:"ARM_CLIENT_SECRET"`
	AzureTenantId       string `mapstructure:"ARM_TENANT_ID"`
	AzureSubscriptionId string `mapstructure:"ARM_SUBSCRIPTION_ID"`
}
type HelmInterfaces struct {
	Name             string            `mapstructure:"name" yaml:"name" json:"name"`
	Namespace        string            `mapstructure:"namespace" yaml:"namespace" json:"namespace"`
	Revision         int               `mapstructure:"revision" yaml:"revision" json:"revision"`
	Updated          string            `mapstructure:"updated" yaml:"updated" json:"updated"`
	Deleted          bool              `mapstructure:"deleted" yaml:"deleted" json:"deleted,omitempty"`
	DependencyUpdate bool              `mapstructure:"dependencyUpdate" yaml:"dependencyUpdate" json:"dependencyUpdate,omitempty"`
	Status           string            `mapstructure:"status" yaml:"status" json:"status"`
	Chart            string            `mapstructure:"chart" yaml:"chart" json:"chart"`
	AppVersion       string            `mapstructure:"appVersion" yaml:"app_version" json:"app_version"`
	Values           map[string]string `mapstructure:"values" yaml:"values" json:"values,omitempty"`
	ValuesFile       string            `mapstructure:"valuesFile" yaml:"valuesFile" json:"valuesFile,omitempty"`
	// TODO: delete
	RepoName string `mapstructure:"repoName" yaml:"repoName" json:"repoName"`
	Repo     string `mapstructure:"repo" yaml:"repo" json:"repo"`
	Url      string `mapstructure:"url" yaml:"url" json:"url"`
	Version  string `mapstructure:"version" yaml:"version" json:"version"`
}

type HelmRelease struct {
	ClusterName    string           `mapstructure:"clusterName" yaml:"clusterName" json:"clusterName,omitempty"`
	ServiceMonitor bool             `mapstructure:"serviceMonitor" yaml:"serviceMonitor" json:"serviceMonitor,omitempty"`
	Verbose        bool             `mapstructure:"verbose" yaml:"verbose" json:"verbose,omitempty"`
	Wait           bool             `mapstructure:"wait" yaml:"wait" json:"wait,omitempty"`
	UpdateStrategy string           `mapstructure:"updateStrategy" yaml:"updateStrategy" json:"updateStrategy,omitempty"` // none, latest
	Releases       []HelmInterfaces `mapstructure:"releases" yaml:"releases" json:"releases,omitempty"`
	Repo           []HelmRepo       `mapstructure:"repo" yaml:"repo" json:"repo,omitempty"`
	// LoadBalancer   LoadBalancer       `mapstructure:"loadBalancer" yaml:"loadBalancer" json:"loadBalancer,omitempty"`
}

type HelmOptions struct {
	ClusterName     string
	CreateNamespace bool
	KubeconfigPath  string
	Overrides       map[string]string
	NodeArch        string
	Helm            *HelmInterfaces
	Wait            bool
	Verbose         bool
	DryRun          bool
	Secrets         []K8sSecret
}

type HelmRepo struct {
	Name    string `mapstructure:"name" yaml:"name" json:"name"`
	Repo    string `mapstructure:"repo" yaml:"repo" json:"repo"`
	Url     string `mapstructure:"url" yaml:"url" json:"url"`
	Version string `mapstructure:"version" yaml:"version" json:"version"`
}

type K8sSecret struct {
	Type        string        `mapstructure:"type" yaml:"type" json:"type"`
	Name        string        `mapstructure:"name" yaml:"name" json:"name"`
	Namespace   string        `mapstructure:"namespace" yaml:"namespace" json:"namespace"`
	SecretsData []SecretsData `mapstructure:"secrets" yaml:"secrets" json:"secrets"`
	Status      string        `mapstructure:"status" yaml:"status" json:"status"`
	Updated     bool          `mapstructure:"updated" yaml:"updated" json:"updated"`
}

type SecretsData struct {
	Type  string // file or literal
	Key   string
	Value string
}

// GetHelmRelease return installed release
func GetHelmRelease(name string, releases []HelmInterfaces) (ok bool, result HelmInterfaces) {
	for _, item := range releases {
		// log.Debugf("Release: %s VERSION: %s STATUS: %s", item.Name, item.AppVersion, item.Status)
		// log.Warnf("%s == %s", item.Name, name)
		if item.Name == name {
			return true, item
		}
	}
	return false, result
}

func (r *Cluster) GetProvider() string {
	return r.Spec.Provider
}

func (r *Cluster) GetUser(name string) User {
	if name == "" {
		name = "root"
	}

	user := User{
		Name: name,
	}
	return user
}

// GetDatastore connection string
func (r *Cluster) GetDatastore(dbPassword string) (string, error) {
	conUrl := ""
	if len(r.Spec.Datastore.Provider) == 0 {
		return "", errors.New("Is not set datastore.provider")
	}
	if len(r.Spec.Datastore.Name) == 0 {
		r.Spec.Datastore.Name = "k3s"
	}
	if len(r.Spec.Datastore.Host) == 0 {
		return "", errors.New("Is not set datastore.host")
	}
	if len(dbPassword) == 0 {
		return "", errors.New("Is not set DB_PASSWORD")
	}
	if len(r.Spec.Datastore.Username) == 0 {
		return "", errors.New("Is not set datastore.username")
	}
	if r.Spec.Datastore.Provider == DatastoreMySql {
		if r.Spec.Datastore.Port == 0 {
			r.Spec.Datastore.Port = 3306
		}
		// K3S_DATASTORE_ENDPOINT='mysql://username:password@tcp(hostname:3306)/k3s' \
		// K3S_DATASTORE_CERTFILE='/path/to/client.crt' \
		// K3S_DATASTORE_KEYFILE='/path/to/client.key' \
		// k3s server
		conUrl = fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/%s", r.Spec.Datastore.Username, dbPassword, r.Spec.Datastore.Host, r.Spec.Datastore.Port, r.Spec.Datastore.Name)
	} else if r.Spec.Datastore.Provider == DatastorePostgreSql {
		if r.Spec.Datastore.Port == 0 {
			r.Spec.Datastore.Port = 5432
		}
		// K3S_DATASTORE_ENDPOINT='postgres://username:password@hostname:5432/k3s' k3s server
		conUrl = fmt.Sprintf("postgres://%s:%s@%s:%d/%s", r.Spec.Datastore.Username, dbPassword, r.Spec.Datastore.Host, r.Spec.Datastore.Port, r.Spec.Datastore.Name)
	} else {
		return "", errors.New(fmt.Sprintf("Is not suport Datastore provider %s.", r.Spec.Datastore.Provider))
	}
	return conUrl, nil
}

// GetTlsSan Add additional hostname or IP as a Subject Alternative Name in the TLS cert
func (r *Cluster) GetTlsSan(node *Node, vpc *Networking) (tlsSAN []string) {
	for _, addr := range vpc.APIServerAddresses {
		if _, isset := Find(tlsSAN, addr.Address); !isset {
			tlsSAN = append(tlsSAN, addr.Address)
		}
	}
	for _, addr := range node.Addresses {
		if _, isset := Find(tlsSAN, addr.Address); !isset {
			tlsSAN = append(tlsSAN, addr.Address)
		}
	}
	return tlsSAN
}

// GetAPIServerUrl url для подключения к API серверу
// TODO: delete
func (r *Cluster) GetAPIServerUrl(masters []ContrelPlanNodes, vpc *Networking, isExternal bool) (apiServerUrl string, err error) {
	if isExternal {
		for _, item := range vpc.APIServerAddresses {
			if string(item.Type) == ExternalDNS {
				return fmt.Sprintf("https://%s:%d", item.Address, vpc.APIServerPort), nil
			} else if string(item.Type) == ExternalIP {
				return fmt.Sprintf("https://%s:%d", item.Address, vpc.APIServerPort), nil
			}
		}
	}
	for _, item := range vpc.APIServerAddresses {
		if string(item.Type) == InternalDNS {
			return fmt.Sprintf("https://%s:%d", item.Address, vpc.APIServerPort), nil
		} else if string(item.Type) == InternalIP {
			return fmt.Sprintf("https://%s:%d", item.Address, vpc.APIServerPort), nil
		} else if string(item.Type) == ExternalDNS {
			return fmt.Sprintf("https://%s:%d", item.Address, vpc.APIServerPort), nil
		} else if string(item.Type) == ExternalIP {
			return fmt.Sprintf("https://%s:%d", item.Address, vpc.APIServerPort), nil
		}
	}

	for _, item := range masters {
		if isExternal {
			nodeIP, ok := r.GetNodeAddress(item.Node, "external")
			if ok {
				return fmt.Sprintf("https://%s:%d", nodeIP, vpc.APIServerPort), nil
			}
		}
		nodeIP, ok := r.GetNodeAddress(item.Node, "internal")
		if ok {
			return fmt.Sprintf("https://%s:%d", nodeIP, vpc.APIServerPort), nil
		}
	}
	return "", fmt.Errorf("Is NOT set Api server URL")
}

// GetAPIServerAddress возвращает hostname  or ip API Server
func (r *Cluster) GetAPIServerAddress(node *Node, vpc *Networking) (apiServerAddres string, err error) {
	for _, item := range vpc.APIServerAddresses {
		apiServerAddres = item.Address
		nodeIP, ok := r.GetNodeAddress(node, "internal")
		// log.Warnf("==> item: %s = val: %s", item, valType)
		if ok {
			if string(item.Type) == InternalDNS {
				return item.Address, nil
			} else if string(item.Type) == InternalIP {
				return item.Address, nil
			} else if string(item.Type) == ExternalDNS {
				return item.Address, nil
			} else if string(item.Type) == ExternalIP {
				return item.Address, nil
			}
			err = errors.New(fmt.Sprintf("is set node internal IP: %s not set APIServerAddresses type (InternalDNS, InternalIP, ExternalDNS, ExternalIP)", nodeIP))
		}

		nodeIP, ok = r.GetNodeAddress(node, "external")
		if ok {
			if string(item.Type) == ExternalDNS {
				return item.Address, nil
			} else if string(item.Type) == ExternalIP {
				return item.Address, nil
			}
			err = errors.New(fmt.Sprintf("is set node external IP: %s not set APIServerAddresses type (InternalDNS, InternalIP, ExternalDNS, ExternalIP)", nodeIP))
		}
		err = errors.New(fmt.Sprintf("is not set node internal or external type (InternalDNS, InternalIP, ExternalDNS, ExternalIP) IP: %s", nodeIP))
	}
	return apiServerAddres, err
}

// GetBastion search and return bastion host
// для работы через baston смотреть README
func (r *Cluster) GetBastion(name string, node *Node) (bastion *BastionNode, err error) {
	log.Debugf("[GetBastion] 1) name: %s node: %s", name, node.Name)
	bastion = &BastionNode{
		SshPort:          SshPortDefault,
		SSHAuthorizedKey: SshKeyDefault,
	}
	if name == "localhost" || name == "127.0.0.1" || name == "local" {
		bastion.Name = "local"
		bastion.Address = "127.0.0.1"
		log.Debugf("[GetBastion] 2) connect address: %s", bastion.Address)
		return bastion, nil
	}
	if len(node.Addresses) == 0 {
		return bastion, fmt.Errorf("Is not set addresses in node %s", node.Name)
	}

	if len(name) == 0 || name == InternalIP || name == InternalDNS || name == ExternalDNS || name == ExternalIP {
		for _, addr := range node.Addresses {
			if name == string(addr.Type) {
				bastion.Address = addr.Address
				bastion.Name = string(addr.Type)
				// bastion.RemoteAddress = addr.Address
				log.Warnf("[GetBastion] 3) connect address: %s, remote address: %s", bastion.Address, bastion.RemoteAddress)
				return bastion, nil
			}
		}
		bastion.Address = node.Addresses[0].Address
		bastion.Name = string(node.Addresses[0].Type)
		bastion.User = node.User
		log.Debugf("[GetBastion] 4) connect address: %s, remote address: %s", bastion.Address, bastion.RemoteAddress)
		return bastion, nil
	}

	for _, item := range r.Spec.Bastions {
		if item.SshPort == 0 {
			item.SshPort = SshPortDefault
		}
		if len(item.SSHAuthorizedKey) == 0 {
			item.SSHAuthorizedKey = SshKeyDefault
		}
		if name != node.Name {
			remoteAddress, ok := r.GetNodeAddress(node, "internal")
			if ok {
				item.RemoteAddress = remoteAddress
				item.RemoteUser = node.User
				item.RemotePort = SshPortDefault
			}
			log.Warnf("[GetBastion] 5) connect address: %s, remote address: %s", item.Address, item.RemoteAddress)
			return item, nil
		}
	}
	return bastion, fmt.Errorf("Is not bastion %s host.", name)
}

func Find(slice []string, val string) (string, bool) {
	// log.Errorln(slice)
	for _, item := range slice {
		// log.Warnf("==> item: %s = val: %s", item, val)
		if item == val {
			return item, true
		}
	}
	return "", false
}

// FindRelease - Find installed Release
func FindRelease(slice []HelmInterfaces, val string) (HelmInterfaces, bool) {
	// log.Errorln(slice)
	for _, item := range slice {
		// log.Warnf("==> item: %s = val: %s", item.Name, val)
		if item.Name == val {
			return item, true
		}
	}
	return HelmInterfaces{}, false
}

func (r *Cluster) GetNodeLabels(node *Node) (cnt int) {
	labels := []string{}
	if len(node.Labels) > 0 {
		if node.Role != "master" {
			labels = append(labels, "")
		}
		return len(node.Labels)
	}
	return 0
}

// GetNodeAddress возращает hostnamr или ip в зависимости от valType (internal|external)
func (r *Cluster) GetNodeAddress(node *Node, valType string) (string, bool) {
	// log.Errorln(slice)
	res := ""
	for _, item := range node.Addresses {
		res = item.Address
		// log.Warnf("==> item: %s = val: %s", item, valType)
		if valType == "internal" {
			if string(item.Type) == InternalDNS {
				return item.Address, true
			} else if string(item.Type) == InternalIP {
				return item.Address, true
			}
		} else if valType == "external" {
			if string(item.Type) == ExternalDNS {
				return item.Address, true
			} else if string(item.Type) == ExternalIP {
				return item.Address, true
			}
		}
	}
	return res, false
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
