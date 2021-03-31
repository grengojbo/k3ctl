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

// K3sOptions k3s options for generate config
type K3sOptions struct {
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	Region       string                  `mapstructure:"region" yaml:"region" json:"region,omitempty"`
	Servers      int                     `mapstructure:"servers" yaml:"servers" json:"servers,omitempty"`         //nolint:lll    // default 1
	Agents       int                     `mapstructure:"agents" yaml:"agents" json:"agents,omitempty"`            //nolint:lll    // default 0
	ClusterToken string                  `mapstructure:"token" yaml:"clusterToken" json:"clusterToken,omitempty"` // default: auto-generated
	Host         string                  `mapstructure:"host" yaml:"host,omitempty" json:"host,omitempty"`
	HostIP       string                  `mapstructure:"hostIP" yaml:"hostIP,omitempty" json:"hostIP,omitempty"`
	Labels       []LabelWithNodeFilters  `mapstructure:"labels" yaml:"labels" json:"labels,omitempty"`
	Env          []EnvVarWithNodeFilters `mapstructure:"env" yaml:"env" json:"env,omitempty"`
	Registries   Registry                `mapstructure:"registries" yaml:"registries,omitempty" json:"registries,omitempty"`
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
