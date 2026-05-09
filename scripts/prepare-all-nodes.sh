#!/usr/bin/env bash
# Запускає prepare-node.sh на всіх нодах кластера iwis-ai через ssh.
# Використовує SSH-ключ ~/.ssh/id_ed25519 і користувача ubuntu (з variables/iwis-ai.yaml).

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NODE_SCRIPT="${SCRIPT_DIR}/prepare-node.sh"
SSH_KEY="${SSH_KEY:-$HOME/.ssh/id_ed25519}"
SSH_USER="${SSH_USER:-ubuntu}"
SSH_OPTS=(-i "${SSH_KEY}" -o StrictHostKeyChecking=accept-new -o ConnectTimeout=10)

# name:ip:role  (role: master | worker | lb)
NODES=(
  "master-1:10.0.40.10:master"
  "worker-1:10.0.40.101:worker"
  "worker-2:10.0.40.102:worker"
  "worker-3:10.0.40.103:worker"
  "lb-1:10.0.40.98:lb"
  "lb-2:10.0.40.99:lb"
)

run_on_node() {
  local name="$1" ip="$2" role="$3"
  echo
  echo "============================================================"
  echo " >>> ${name} (${ip}) role=${role}"
  echo "============================================================"
  scp "${SSH_OPTS[@]}" "${NODE_SCRIPT}" "${SSH_USER}@${ip}:/tmp/prepare-node.sh"
  ssh "${SSH_OPTS[@]}" "${SSH_USER}@${ip}" \
    "chmod +x /tmp/prepare-node.sh && NODE_HOSTNAME='${name}' NODE_ROLE='${role}' sudo -E bash /tmp/prepare-node.sh && rm -f /tmp/prepare-node.sh"
}

# Якщо передані аргументи — обробити лише ці ноди (за іменем).
FILTER=("$@")

for entry in "${NODES[@]}"; do
  IFS=':' read -r name ip role <<<"$entry"
  if [[ ${#FILTER[@]} -gt 0 ]]; then
    skip=1
    for f in "${FILTER[@]}"; do [[ "$f" == "$name" ]] && skip=0; done
    [[ $skip -eq 1 ]] && continue
  fi
  run_on_node "$name" "$ip" "$role"
done

echo
echo "==> Усі ноди підготовлені."
