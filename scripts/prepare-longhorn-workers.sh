#!/usr/bin/env bash
# Підготовка Longhorn на всіх worker-нодах: диск + iscsi_tcp модуль + вимкнення multipathd

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
DISK_SCRIPT="${SCRIPT_DIR}/prepare-longhorn-disk.sh"
SSH_KEY="${SSH_KEY:-$HOME/.ssh/id_ed25519}"
SSH_USER="${SSH_USER:-ubuntu}"
SSH_OPTS=(-i "${SSH_KEY}" -o StrictHostKeyChecking=accept-new -o ConnectTimeout=10)

WORKERS=(
  "worker-1:10.0.40.101"
  "worker-2:10.0.40.102"
  "worker-3:10.0.40.103"
)

run_on_worker() {
  local name="$1" ip="$2"
  echo
  echo "============================================================"
  echo " >>> ${name} (${ip})"
  echo "============================================================"
  
  # 1. Підготовка диска
  scp "${SSH_OPTS[@]}" "${DISK_SCRIPT}" "${SSH_USER}@${ip}:/tmp/prepare-longhorn-disk.sh"
  ssh "${SSH_OPTS[@]}" "${SSH_USER}@${ip}" \
    "chmod +x /tmp/prepare-longhorn-disk.sh && sudo bash /tmp/prepare-longhorn-disk.sh && rm -f /tmp/prepare-longhorn-disk.sh"
  
  # 2. Увімкнення модуля iscsi_tcp (потрібен для Longhorn)
  echo "==> Увімкнення iscsi_tcp модуля"
  ssh "${SSH_OPTS[@]}" "${SSH_USER}@${ip}" '
    echo "iscsi_tcp" | sudo tee /etc/modules-load.d/iscsi_tcp.conf
    sudo modprobe iscsi_tcp
    lsmod | grep -q iscsi_tcp && echo "iscsi_tcp loaded OK" || echo "ERROR: iscsi_tcp not loaded"
  '
  
  # 3. Вимкнення multipathd (конфліктує з Longhorn)
  echo "==> Вимкнення multipathd"
  ssh "${SSH_OPTS[@]}" "${SSH_USER}@${ip}" '
    sudo systemctl stop multipathd 2>/dev/null || true
    sudo systemctl disable multipathd 2>/dev/null || true
    sudo systemctl mask multipathd 2>/dev/null || true
    echo "multipathd status: $(systemctl is-active multipathd 2>/dev/null || echo inactive)"
  '
  
  # 4. Перезапуск iscsid
  echo "==> Перезапуск iscsid"
  ssh "${SSH_OPTS[@]}" "${SSH_USER}@${ip}" \
    'sudo systemctl restart iscsid && sudo systemctl is-active iscsid'
}

FILTER=("$@")

for entry in "${WORKERS[@]}"; do
  name="${entry%%:*}"
  ip="${entry##*:}"
  if [[ ${#FILTER[@]} -gt 0 ]]; then
    skip=1
    for f in "${FILTER[@]}"; do [[ "$f" == "$name" ]] && skip=0; done
    [[ $skip -eq 1 ]] && continue
  fi
  run_on_worker "$name" "$ip"
done

echo
echo "==> Усі worker-ноди підготовлені для Longhorn."
echo "==> Рекомендується перевірити: curl -sSfL https://raw.githubusercontent.com/longhorn/longhorn/v1.7.2/scripts/environment_check.sh | bash"
