//go:build !ignore_autogenerated
// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1alpha1

import (
	"github.com/docker/go-connections/nat"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/cluster-api/api/v1alpha3"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Addons) DeepCopyInto(out *Addons) {
	*out = *in
	in.Ingress.DeepCopyInto(&out.Ingress)
	in.CertManager.DeepCopyInto(&out.CertManager)
	in.Registries.DeepCopyInto(&out.Registries)
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Addons.
func (in *Addons) DeepCopy() *Addons {
	if in == nil {
		return nil
	}
	out := new(Addons)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *AgentOpts) DeepCopyInto(out *AgentOpts) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new AgentOpts.
func (in *AgentOpts) DeepCopy() *AgentOpts {
	if in == nil {
		return nil
	}
	out := new(AgentOpts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *BastionNode) DeepCopyInto(out *BastionNode) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new BastionNode.
func (in *BastionNode) DeepCopy() *BastionNode {
	if in == nil {
		return nil
	}
	out := new(BastionNode)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *CertManager) DeepCopyInto(out *CertManager) {
	*out = *in
	// if in.Values != nil {
	// 	in, out := &in.Values, &out.Values
	// 	*out = make([]string, len(*in))
	// 	copy(*out, *in)
	// }
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new CertManager.
func (in *CertManager) DeepCopy() *CertManager {
	if in == nil {
		return nil
	}
	out := new(CertManager)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Cluster) DeepCopyInto(out *Cluster) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Cluster.
func (in *Cluster) DeepCopy() *Cluster {
	if in == nil {
		return nil
	}
	out := new(Cluster)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *Cluster) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterList) DeepCopyInto(out *ClusterList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Cluster, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterList.
func (in *ClusterList) DeepCopy() *ClusterList {
	if in == nil {
		return nil
	}
	out := new(ClusterList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ClusterList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterSpec) DeepCopyInto(out *ClusterSpec) {
	*out = *in
	if in.Bastions != nil {
		in, out := &in.Bastions, &out.Bastions
		*out = make([]*BastionNode, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(BastionNode)
				**out = **in
			}
		}
	}
	if in.Nodes != nil {
		in, out := &in.Nodes, &out.Nodes
		*out = make([]*Node, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(Node)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make([]LabelWithNodeFilters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]EnvVarWithNodeFilters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Options = in.Options
	in.K3sOptions.DeepCopyInto(&out.K3sOptions)
	out.LoadBalancer = in.LoadBalancer
	in.Addons.DeepCopyInto(&out.Addons)
	out.KubeconfigOptions = in.KubeconfigOptions
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]VolumeWithNodeFilters, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	out.Datastore = in.Datastore
	in.Networking.DeepCopyInto(&out.Networking)
	if in.Files != nil {
		in, out := &in.Files, &out.Files
		*out = make([]File, len(*in))
		copy(*out, *in)
	}
	if in.PreCommands != nil {
		in, out := &in.PreCommands, &out.PreCommands
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.PostCommands != nil {
		in, out := &in.PostCommands, &out.PostCommands
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Users != nil {
		in, out := &in.Users, &out.Users
		*out = make([]User, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.NTP != nil {
		in, out := &in.NTP, &out.NTP
		*out = new(NTP)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterSpec.
func (in *ClusterSpec) DeepCopy() *ClusterSpec {
	if in == nil {
		return nil
	}
	out := new(ClusterSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ClusterStatus) DeepCopyInto(out *ClusterStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ClusterStatus.
func (in *ClusterStatus) DeepCopy() *ClusterStatus {
	if in == nil {
		return nil
	}
	out := new(ClusterStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ContrelPlanNodes) DeepCopyInto(out *ContrelPlanNodes) {
	*out = *in
	if in.Bastion != nil {
		in, out := &in.Bastion, &out.Bastion
		*out = new(BastionNode)
		**out = **in
	}
	if in.Node != nil {
		in, out := &in.Node, &out.Node
		*out = new(Node)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ContrelPlanNodes.
func (in *ContrelPlanNodes) DeepCopy() *ContrelPlanNodes {
	if in == nil {
		return nil
	}
	out := new(ContrelPlanNodes)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Datastore) DeepCopyInto(out *Datastore) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Datastore.
func (in *Datastore) DeepCopy() *Datastore {
	if in == nil {
		return nil
	}
	out := new(Datastore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *EnvVarWithNodeFilters) DeepCopyInto(out *EnvVarWithNodeFilters) {
	*out = *in
	if in.NodeFilters != nil {
		in, out := &in.NodeFilters, &out.NodeFilters
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new EnvVarWithNodeFilters.
func (in *EnvVarWithNodeFilters) DeepCopy() *EnvVarWithNodeFilters {
	if in == nil {
		return nil
	}
	out := new(EnvVarWithNodeFilters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExposureOpts) DeepCopyInto(out *ExposureOpts) {
	*out = *in
	out.PortMapping = in.PortMapping
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExposureOpts.
func (in *ExposureOpts) DeepCopy() *ExposureOpts {
	if in == nil {
		return nil
	}
	out := new(ExposureOpts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ExternalDatastore) DeepCopyInto(out *ExternalDatastore) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ExternalDatastore.
func (in *ExternalDatastore) DeepCopy() *ExternalDatastore {
	if in == nil {
		return nil
	}
	out := new(ExternalDatastore)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *File) DeepCopyInto(out *File) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new File.
func (in *File) DeepCopy() *File {
	if in == nil {
		return nil
	}
	out := new(File)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Ingress) DeepCopyInto(out *Ingress) {
	*out = *in
	// if in.Values != nil {
	// 	in, out := &in.Values, &out.Values
	// 	// *out = make([]string, len(*in))
	// 	copy(*out, *in)
	// }
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Ingress.
func (in *Ingress) DeepCopy() *Ingress {
	if in == nil {
		return nil
	}
	out := new(Ingress)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LabelWithNodeFilters) DeepCopyInto(out *LabelWithNodeFilters) {
	*out = *in
	if in.NodeFilters != nil {
		in, out := &in.NodeFilters, &out.NodeFilters
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LabelWithNodeFilters.
func (in *LabelWithNodeFilters) DeepCopy() *LabelWithNodeFilters {
	if in == nil {
		return nil
	}
	out := new(LabelWithNodeFilters)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *LoadBalancer) DeepCopyInto(out *LoadBalancer) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new LoadBalancer.
func (in *LoadBalancer) DeepCopy() *LoadBalancer {
	if in == nil {
		return nil
	}
	out := new(LoadBalancer)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NTP) DeepCopyInto(out *NTP) {
	*out = *in
	if in.Servers != nil {
		in, out := &in.Servers, &out.Servers
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Enabled != nil {
		in, out := &in.Enabled, &out.Enabled
		*out = new(bool)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NTP.
func (in *NTP) DeepCopy() *NTP {
	if in == nil {
		return nil
	}
	out := new(NTP)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Networking) DeepCopyInto(out *Networking) {
	*out = *in
	if in.APIServerAddresses != nil {
		in, out := &in.APIServerAddresses, &out.APIServerAddresses
		*out = make(v1alpha3.MachineAddresses, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Networking.
func (in *Networking) DeepCopy() *Networking {
	if in == nil {
		return nil
	}
	out := new(Networking)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Node) DeepCopyInto(out *Node) {
	*out = *in
	if in.Volumes != nil {
		in, out := &in.Volumes, &out.Volumes
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Env != nil {
		in, out := &in.Env, &out.Env
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Cmd != nil {
		in, out := &in.Cmd, &out.Cmd
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Args != nil {
		in, out := &in.Args, &out.Args
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.Ports != nil {
		in, out := &in.Ports, &out.Ports
		*out = make(nat.PortMap, len(*in))
		for key, val := range *in {
			var outVal []nat.PortBinding
			if val == nil {
				(*out)[key] = nil
			} else {
				in, out := &val, &outVal
				*out = make([]nat.PortBinding, len(*in))
				copy(*out, *in)
			}
			(*out)[key] = outVal
		}
	}
	if in.Labels != nil {
		in, out := &in.Labels, &out.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if in.Networks != nil {
		in, out := &in.Networks, &out.Networks
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ExtraHosts != nil {
		in, out := &in.ExtraHosts, &out.ExtraHosts
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	in.ServerOpts.DeepCopyInto(&out.ServerOpts)
	out.AgentOpts = in.AgentOpts
	out.State = in.State
	if in.Addresses != nil {
		in, out := &in.Addresses, &out.Addresses
		*out = make(v1alpha3.MachineAddresses, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Node.
func (in *Node) DeepCopy() *Node {
	if in == nil {
		return nil
	}
	out := new(Node)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *NodeState) DeepCopyInto(out *NodeState) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new NodeState.
func (in *NodeState) DeepCopy() *NodeState {
	if in == nil {
		return nil
	}
	out := new(NodeState)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Options) DeepCopyInto(out *Options) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Options.
func (in *Options) DeepCopy() *Options {
	if in == nil {
		return nil
	}
	out := new(Options)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *Registry) DeepCopyInto(out *Registry) {
	*out = *in
	if in.Use != nil {
		in, out := &in.Use, &out.Use
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new Registry.
func (in *Registry) DeepCopy() *Registry {
	if in == nil {
		return nil
	}
	out := new(Registry)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ServerOpts) DeepCopyInto(out *ServerOpts) {
	*out = *in
	if in.KubeAPI != nil {
		in, out := &in.KubeAPI, &out.KubeAPI
		*out = new(ExposureOpts)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ServerOpts.
func (in *ServerOpts) DeepCopy() *ServerOpts {
	if in == nil {
		return nil
	}
	out := new(ServerOpts)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimpleConfigOptionsK3s) DeepCopyInto(out *SimpleConfigOptionsK3s) {
	*out = *in
	if in.ExtraServerArgs != nil {
		in, out := &in.ExtraServerArgs, &out.ExtraServerArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if in.ExtraAgentArgs != nil {
		in, out := &in.ExtraAgentArgs, &out.ExtraAgentArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimpleConfigOptionsK3s.
func (in *SimpleConfigOptionsK3s) DeepCopy() *SimpleConfigOptionsK3s {
	if in == nil {
		return nil
	}
	out := new(SimpleConfigOptionsK3s)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *SimpleConfigOptionsKubeconfig) DeepCopyInto(out *SimpleConfigOptionsKubeconfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new SimpleConfigOptionsKubeconfig.
func (in *SimpleConfigOptionsKubeconfig) DeepCopy() *SimpleConfigOptionsKubeconfig {
	if in == nil {
		return nil
	}
	out := new(SimpleConfigOptionsKubeconfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *User) DeepCopyInto(out *User) {
	*out = *in
	if in.Gecos != nil {
		in, out := &in.Gecos, &out.Gecos
		*out = new(string)
		**out = **in
	}
	if in.Groups != nil {
		in, out := &in.Groups, &out.Groups
		*out = new(string)
		**out = **in
	}
	if in.HomeDir != nil {
		in, out := &in.HomeDir, &out.HomeDir
		*out = new(string)
		**out = **in
	}
	if in.Inactive != nil {
		in, out := &in.Inactive, &out.Inactive
		*out = new(bool)
		**out = **in
	}
	if in.Shell != nil {
		in, out := &in.Shell, &out.Shell
		*out = new(string)
		**out = **in
	}
	if in.Passwd != nil {
		in, out := &in.Passwd, &out.Passwd
		*out = new(string)
		**out = **in
	}
	if in.PrimaryGroup != nil {
		in, out := &in.PrimaryGroup, &out.PrimaryGroup
		*out = new(string)
		**out = **in
	}
	if in.LockPassword != nil {
		in, out := &in.LockPassword, &out.LockPassword
		*out = new(bool)
		**out = **in
	}
	if in.Sudo != nil {
		in, out := &in.Sudo, &out.Sudo
		*out = new(string)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new User.
func (in *User) DeepCopy() *User {
	if in == nil {
		return nil
	}
	out := new(User)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *VolumeWithNodeFilters) DeepCopyInto(out *VolumeWithNodeFilters) {
	*out = *in
	if in.NodeFilters != nil {
		in, out := &in.NodeFilters, &out.NodeFilters
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new VolumeWithNodeFilters.
func (in *VolumeWithNodeFilters) DeepCopy() *VolumeWithNodeFilters {
	if in == nil {
		return nil
	}
	out := new(VolumeWithNodeFilters)
	in.DeepCopyInto(out)
	return out
}
