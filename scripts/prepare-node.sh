#!/usr/bin/env bash
# Підготовка ноди Ubuntu під k3s (iwis-ai).
# Викликається на цільовій ноді з sudo.

set -euo pipefail

# 1) Hostname задається ззовні (через --hostname або автоматично).
if [[ -n "${NODE_HOSTNAME:-}" ]]; then
  echo "==> hostnamectl set-hostname ${NODE_HOSTNAME}"
  sudo hostnamectl set-hostname "${NODE_HOSTNAME}"
fi

# 2) Пакети
echo "==> apt-get update && install"
sudo DEBIAN_FRONTEND=noninteractive apt-get update -y
sudo DEBIAN_FRONTEND=noninteractive apt-get install -y \
  curl iptables open-iscsi nfs-common jq apparmor-utils cryptsetup dmsetup

# iscsid треба лише для Longhorn (worker-ноди).
case "${NODE_ROLE:-}" in
  worker)
    echo "==> enable iscsid (NODE_ROLE=worker, для Longhorn)"
    sudo systemctl enable --now iscsid
    ;;
  *)
    echo "==> skip iscsid (NODE_ROLE='${NODE_ROLE:-}')"
    ;;
esac

# 3) Модулі ядра
echo "==> /etc/modules-load.d/k3s.conf"
sudo tee /etc/modules-load.d/k3s.conf >/dev/null <<'EOF'
br_netfilter
overlay
ip_vs
ip_vs_rr
ip_vs_wrr
ip_vs_sh
nf_conntrack
EOF
sudo modprobe br_netfilter overlay ip_vs ip_vs_rr ip_vs_wrr ip_vs_sh nf_conntrack

# 4) sysctl
echo "==> /etc/sysctl.d/99-k3s.conf"
sudo tee /etc/sysctl.d/99-k3s.conf >/dev/null <<'EOF'
net.ipv4.ip_forward = 1
net.bridge.bridge-nf-call-iptables = 1
net.bridge.bridge-nf-call-ip6tables = 1
fs.inotify.max_user_instances = 8192
fs.inotify.max_user_watches = 524288
EOF
sudo sysctl --system >/dev/null

# 5) swap off
echo "==> swap off"
sudo swapoff -a
if grep -E '^[^#].* swap ' /etc/fstab >/dev/null; then
  sudo sed -i.bak '/ swap / s/^/#/' /etc/fstab
fi

echo "==> Готово на $(hostname)"
