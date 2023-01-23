#!/bin/bash

# set -ex

[ ${UID} -ne 0 ] && {
  echo "RUN only ROOT user"
  exit 1;
}

{{- if not .noShow }}
# echo "prepare the k3s config directory"
# mkdir -p /etc/rancher/k3s
# # move the config file into place
# mv /tmp/config.yaml /etc/rancher/k3s/config.yaml
# # if the server has already been initialized just stop here
# [ -e /etc/rancher/k3s/k3s.yaml ] && exit 0;

# # apply_k3s_selinux
# /sbin/semodule -v -i /usr/share/selinux/packages/k3s.pp

# curl -sfL https://get.k3s.io | INSTALL_K3S_SKIP_START=true INSTALL_K3S_SKIP_SELINUX_RPM=true INSTALL_K3S_CHANNEL=stable INSTALL_K3S_EXEC=server sh -
# curl -sfL https://get.k3s.io | INSTALL_K3S_SKIP_START=true INSTALL_K3S_SKIP_SELINUX_RPM=true INSTALL_K3S_CHANNEL=stable INSTALL_K3S_EXEC=agent sh -

# # Default k3s node labels
# # default_agent_labels = ["k3s_upgrade=true"]
# # default_control_plane_labels = ["k3s_upgrade=true"]
# #  default_control_plane_taints = concat([], local.allow_scheduling_on_control_plane ? [] : ["node-role.kubernetes.io/master:NoSchedule"])
{{- end }}

InstallK3s() {
  echo "Install k3s {{ .role }}"
  {{ .command }}
}

SaveK3sConfigFile() {
	local CONFIG_FILE_K3S=/etc/rancher/k3s/config.yaml
	# Save k3s config file
	echo "Save k3s config file"
	mkdir -p /etc/rancher/k3s
	# cp /tmp/config.yaml /etc/rancher/k3s/config.yaml
	cat > ${CONFIG_FILE_K3S} << EOF
{{ .configK3s }}
EOF

  # chmod 755 ${CONFIG_FILE_K3S}
  # ${CONFIG_FILE_K3S}
}

WaitStartServer() {
systemctl start k3s 2> /dev/null
# systemctl start k3s
echo "prepare the post_install directory"
mkdir -p /var/post_install
timeout 120 bash <<EOF
  until systemctl status k3s > /dev/null; do
    systemctl start k3s 2> /dev/null
  	echo "Waiting for the k3s server to start..."
  	sleep 2
  done
  until [ -e /etc/rancher/k3s/k3s.yaml ]; do
    echo "Waiting for kubectl config..."
    sleep 2
  done
  until [[ "\$(kubectl get --raw='/readyz' 2> /dev/null)" == "ok" ]]; do
    echo "Waiting for the cluster to become ready..."
    sleep 2
  done
EOF
}

WaitStartAgent() {
systemctl start k3s-agent 2> /dev/null
timeout 120 bash <<EOF
  until systemctl status k3s-agent > /dev/null; do
  	systemctl start k3s-agent 2> /dev/null
  	echo "Waiting for the k3s agent to start..."
  	sleep 2
  done
EOF
}

SaveK3sConfigFile
InstallK3s

{{- if eq .role "server" }}
WaitStartServer
{{- else }}
WaitStartAgent
{{- end }}

{{- if not .noShow }}
# Upon reboot start k3s and wait for it to be ready to receive commands


# provisioner "file" {
#     content = yamlencode(merge({
#       node-name                   = module.control_planes[each.key].name
#       server                      = length(module.control_planes) == 1 ? null : "https://${module.control_planes[each.key].private_ipv4_address == module.control_planes[keys(module.control_planes)[0]].private_ipv4_address ? module.control_planes[keys(module.control_planes)[1]].private_ipv4_address : module.control_planes[keys(module.control_planes)[0]].private_ipv4_address}:6443"
#       token                       = random_password.k3s_token.result
#       disable-cloud-controller    = true
#       disable                     = local.disable_extras
#       flannel-iface               = "eth1"
#       kubelet-arg                 = ["cloud-provider=external", "volume-plugin-dir=/var/lib/kubelet/volumeplugins"]
#       kube-controller-manager-arg = "flex-volume-plugin-dir=/var/lib/kubelet/volumeplugins"
#       node-ip                     = module.control_planes[each.key].private_ipv4_address
#       advertise-address           = module.control_planes[each.key].private_ipv4_address
#       node-label                  = each.value.labels
#       node-taint                  = each.value.taints
#       disable-network-policy      = var.cni_plugin == "calico" ? true : var.disable_network_policy
#       write-kubeconfig-mode       = "0644" # needed for import into rancher
#       },
#       var.cni_plugin == "calico" ? {
#         flannel-backend = "none"
#     } : {}))

#     destination = "/tmp/config.yaml"
#   }

# # Generating k3s master config file
#   provisioner "file" {
#     content = yamlencode(merge({
#       node-name                   = module.control_planes[keys(module.control_planes)[0]].name
#       token                       = random_password.k3s_token.result
#       cluster-init                = true
#       disable-cloud-controller    = true
#       disable                     = local.disable_extras
#       flannel-iface               = "eth1"
#       kubelet-arg                 = ["cloud-provider=external", "volume-plugin-dir=/var/lib/kubelet/volumeplugins"]
#       kube-controller-manager-arg = "flex-volume-plugin-dir=/var/lib/kubelet/volumeplugins"
#       node-ip                     = module.control_planes[keys(module.control_planes)[0]].private_ipv4_address
#       advertise-address           = module.control_planes[keys(module.control_planes)[0]].private_ipv4_address
#       node-taint                  = local.control_plane_nodes[keys(module.control_planes)[0]].taints
#       node-label                  = local.control_plane_nodes[keys(module.control_planes)[0]].labels
#       disable-network-policy      = var.cni_plugin == "calico" ? true : var.disable_network_policy
#       },
#       var.cni_plugin == "calico" ? {
#         flannel-backend = "none"
#     } : {}))

#     destination = "/tmp/config.yaml"
#   }
{{- end }}