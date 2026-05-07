# Runbook: розгортання k3s-кластера `iwis-ai` (домен `iwis.dev`)

HA-готовий k3s в Hetzner: embedded etcd, Cilium (kube-proxy replacement, MTU 1400), Kube-VIP для control-plane та service VIP, HAProxy Ingress на dedicated lb-нодах, Longhorn для PV, ExternalDNS + cert-manager на Cloudflare, публікація через DNAT `116.202.72.52` + Hetzner VSwitch IP `168.119.131.148`, Cloudflare Tunnel для адмін-сервісів.

> Конфіг: `variables/iwis-ai.yaml`, values: `variables/iwis-ai/*.yaml`.

---

## 1. Топологія

| Роль | Hostname | Internal IP | Зовнішній IP / Примітка |
|---|---|---|---|
| master / etcd | `master-1` | `10.0.40.10/24` | (надалі +`master-2` 10.0.40.11, +`master-3` 10.0.40.12) |
| worker | `worker-1` | `10.0.40.101/24` | Longhorn |
| worker | `worker-2` | `10.0.40.102/24` | Longhorn |
| worker | `worker-3` | `10.0.40.103/24` | Longhorn |
| lb / ingress | `lb-1` | `10.0.40.98/24` | `168.119.131.148/29` (gw `168.119.131.145`, MTU=1400) на VSwitch |
| lb / ingress | `lb-2` | `10.0.40.99/24` | (опційно ще IP з `168.119.131.144/29`) |

**Віртуальні IP (Kube-VIP, ARP, у мережі 10.0.40.0/24):**
- `10.0.40.20` — VIP API-сервера (`6443`), плаває між master-ами.
- `10.0.40.100` — VIP HAProxy Ingress (`80/443`), плаває між `lb-1` / `lb-2`.

**Зовнішня публікація:**
- `116.202.72.52` → DNAT `TCP 80,443` → `10.0.40.100`. DNS `*.iwis.dev` (Cloudflare).
- `168.119.131.148` — secondary IP на VSwitch-інтерфейсі `lb-1` (MTU=1400) для прямого виходу/резерву.
- Cloudflare Tunnel — приватні `*.admin.iwis.dev` без публічного IP.

---

## 2. Технологічний стек

| Шар | Рішення | Чому |
|---|---|---|
| Datastore | embedded etcd (`datastore.provider: etcd` → `--cluster-init` авто) | HA з 3 master, без зовнішньої БД |
| CNI | **Cilium** (kube-proxy replacement, eBPF) | observability, NetworkPolicy L7, BGP/L2, MTU гнучкий |
| L4/VIP | Kube-VIP (control-plane + cloud-provider) | без MetalLB, ARP у L2 10.0.40.0/24 |
| Ingress | HAProxy Ingress (DaemonSet, hostNetwork) | високий RPS, гнучкий L4/L7 |
| Storage | Longhorn (тільки workers) | PV, snapshots, backup до S3 |
| TLS | cert-manager + Cloudflare DNS-01 | wildcard для `*.iwis.dev` |
| DNS | ExternalDNS (Cloudflare) | авто-A-записи з Ingress |
| Tunnel | cloudflared Deployment | приватні адмін-UI |

---

## 3. Підготовка VM (Ubuntu 22.04+) — на КОЖНІЙ ноді

Повна підготовка ноди винесена у [`scripts/prepare-node.sh`](../scripts/prepare-node.sh).
Скрипт: встановлює пакети, завантажує модулі ядра (`br_netfilter`, `overlay`, `ip_vs`, `nf_conntrack`), налаштовує sysctl, вимикає swap. Для worker-нод (`NODE_ROLE=worker`) додатково вмикає `iscsid` (Longhorn).

```bash
# Передумова: користувач ubuntu з sudo NOPASSWD + ssh-ключ у ~/.ssh/authorized_keys
# Hostname задати відповідно до таблиці §1

# Запустити на ноді (або через ssh-pipe):
curl -fsSL https://raw.githubusercontent.com/grengojbo/k3ctl/main/scripts/prepare-node.sh \
  | NODE_HOSTNAME=master-1 NODE_ROLE=master sudo -E bash

# Або скопіювати і запустити вручну:
scp scripts/prepare-node.sh ubuntu@10.0.40.10:~
ssh ubuntu@10.0.40.10 "NODE_HOSTNAME=master-1 NODE_ROLE=master bash prepare-node.sh"

# Для worker-1:
ssh ubuntu@10.0.40.101 "NODE_HOSTNAME=worker-1 NODE_ROLE=worker bash prepare-node.sh"
# Для lb-1:
ssh ubuntu@10.0.40.98  "NODE_HOSTNAME=lb-1 NODE_ROLE=lb bash prepare-node.sh"
```

### 3.1. Тільки на `lb-1` — VSwitch інтерфейс з MTU 1400

`/etc/netplan/60-vswitch.yaml` (vlan id уточнити в Hetzner Robot → vSwitch):

```yaml
network:
  version: 2
  vlans:
    enp7s0.4000:
      id: 4000
      link: enp7s0
      mtu: 1400
      addresses: [168.119.131.148/29]
      routes:
        - to: 0.0.0.0/0
          via: 168.119.131.145
          table: 100
          metric: 200
      routing-policy:
        - from: 168.119.131.148
          table: 100
```

```bash
sudo chmod 600 /etc/netplan/60-vswitch.yaml
sudo netplan apply
ip a show enp7s0.4000
ping -M do -s 1372 -c 2 168.119.131.145    # перевірка MTU 1400
```

### 3.2. DNAT 116.202.72.52 → 10.0.40.100

Налаштувати на стороні шлюза/Hetzner Cloud Firewall:
- `TCP 80  116.202.72.52 -> 10.0.40.100:80`
- `TCP 443 116.202.72.52 -> 10.0.40.100:443`

DNS у Cloudflare: `*.iwis.dev` A `116.202.72.52` (proxied=false на час видачі LE; потім за бажанням proxied=true).

---

## 4. Розгортання через `k3ctl`

> З робочої машини. `KUBECONFIG=~/.kube/iwis-ai.yaml`.

### 4.0. Ролі команд `k3ctl`

| Команда | Що робить |
|---|---|
| `k3ctl cluster create iwis-ai -c ./variables/iwis-ai.yaml` | Створює k3s-кластер (SSH на ноди, встановлює `k3s server`/`agent`, формує `--tls-san`, `--cluster-init` тощо). Alias: `k3ctl cluster apply`. |
| `k3ctl apply -c iwis-ai [addon]` | Ставить Helm addons з `spec.addons`: `cert-manager`, `ingress`, `external-dns`, `monitoring`. Без аргументу — всі разом. |
| ручно | Cilium, Kube-VIP, Longhorn, Cloudflare Tunnel, ClusterIssuer — `k3ctl` їх не вміє. |

### 4.1. Dry-run
```bash
k3ctl cluster create iwis-ai -c ./variables/iwis-ai.yaml --dry-run
```

### 4.2. Створити кластер (всі ноди)

`k3ctl cluster create` бере всі ноди зі `spec.nodes`: ставить перший master з `--cluster-init`, решту master/agent приєднує через адресу з `apiServerAddresses`. Ноди залишаться в `NotReady` — це очікувано (CNI вимкнено через `--flannel-backend=none`, стануть Ready після §4.3 Cilium).

```bash
k3ctl cluster create iwis-ai -c ./variables/iwis-ai.yaml

k3ctl kubeconfig get iwis-ai
kubectl get nodes -o wide   # всі NotReady — окей
kubectl -n kube-system get pods
```

> **Рекомендація після створення**: замінити аргументи CLI у `/etc/systemd/system/k3s.service` на `/etc/rancher/k3s/config.yaml` (div. §9.3).

### 4.3. Cilium (kube-proxy replacement)

> Поки Kube-VIP ще не піднято, тимчасово виставити `k8sServiceHost=10.0.40.10` у `cilium-values.yaml`, після Kube-VIP — повернути `10.0.40.20`.

```bash
helm repo add cilium https://helm.cilium.io && helm repo update
helm upgrade --install cilium cilium/cilium -n kube-system \
  -f ./variables/iwis-ai/cilium-values.yaml
cilium status --wait
cilium connectivity test    # опційно
```

### 4.4. Kube-VIP — control-plane VIP `10.0.40.20` (DaemonSet)

Рекомендований спосіб — DaemonSet з офіційною RBAC (підтримує leader election, auto-failover):

```bash
# 1. RBAC (k3s авто-застосує з server/manifests)
ssh ubuntu@10.0.40.10 "sudo curl -fsSL https://kube-vip.io/manifests/rbac.yaml \
  -o /var/lib/rancher/k3s/server/manifests/kube-vip-rbac.yaml"

# 2. DaemonSet (вже налаштовано у kube-vip-cp.yaml)
kubectl apply -f ./variables/iwis-ai/kube-vip-cp.yaml

# 3. Перевірити (~10-15с)
kubectl -n kube-system get pods -l app.kubernetes.io/name=kube-vip-ds
kubectl -n kube-system logs -l app.kubernetes.io/name=kube-vip-ds --tail=10
ping -c3 10.0.40.20
```

Очікуваний результат логів:
```
successfully acquired lease kube-system/plndr-cp-lock
Node [master-1] is assuming leadership of the cluster
Gratuitous Arp broadcast will repeat every 3 seconds for [10.0.40.20/eth0]
```

Після підняття VIP — оновити Cilium на `k8sServiceHost: 10.0.40.20`:
```bash
helm upgrade cilium cilium/cilium -n kube-system -f ./variables/iwis-ai/cilium-values.yaml
```

> **Важливо**: переконатися, що `10.0.40.20` присутній у TLS SANs API-сервера (див. §9.2).

### 4.5. Лейбли та taint-и нодам

Лейбл `role: lb` вже задано в `spec.nodes` для lb-1/lb-2 і `k3ctl` його застосує. Додаємо taint і лейбли для Longhorn:

```bash
# Taint для lb-нод (k3ctl єму не виставляє)
kubectl taint node lb-1 lb-2 dedicated=lb:NoSchedule --overwrite

# Longhorn диски — тільки на worker-ах
kubectl label node worker-1 worker-2 worker-3 \
  node.longhorn.io/create-default-disk=true --overwrite
kubectl label node lb-1 lb-2 master-1 \
  node.longhorn.io/create-default-disk=false --overwrite
```

### 4.6. Kube-VIP Cloud Provider (service VIP `10.0.40.100`)
```bash
kubectl apply -f https://kube-vip.io/manifests/rbac.yaml
kubectl apply -f https://raw.githubusercontent.com/kube-vip/kube-vip-cloud-provider/main/manifest/kube-vip-cloud-controller.yaml
kubectl apply -f ./variables/iwis-ai/kube-vip-cloud.yaml

# Kube-VIP DaemonSet на lb-нодах для анонсу LoadBalancer IP (svc_enable=true).
# Скористатися генератором з kube-vip:
docker run --rm ghcr.io/kube-vip/kube-vip:v0.8.0 manifest daemonset \
  --interface eth0 --services --inCluster --arp --leaderElection \
  | kubectl apply -f -

kubectl -n kube-system patch ds kube-vip-ds --type merge -p '{
  "spec":{"template":{"spec":{
    "nodeSelector":{"role":"lb"},
    "tolerations":[{"key":"dedicated","operator":"Equal","value":"lb","effect":"NoSchedule"}]
  }}}}'
```

### 4.7. Longhorn
```bash
helm repo add longhorn https://charts.longhorn.io && helm repo update
kubectl create ns longhorn-system
helm upgrade --install longhorn longhorn/longhorn -n longhorn-system \
  -f ./variables/iwis-ai/longhorn-values.yaml
kubectl -n longhorn-system get pods -w
```

### 4.8. Helm-addons через `k3ctl apply`

`k3ctl apply -c iwis-ai` бере секцію `spec.addons` з `variables/iwis-ai.yaml` і ставить HAProxy Ingress, cert-manager, ExternalDNS, monitoring. Підтримувані addon-и: `cert-manager`, `ingress`, `monitoring`, `external-dns` (`pkg/types/module.go`).

```bash
# Передумова: секрети для Cloudflare API у потрібних namespace.
export CF_API_TOKEN=...     # Cloudflare API Token (Zone:DNS:Edit для iwis.dev)
kubectl create ns cert-manager 2>/dev/null || true
kubectl -n cert-manager create secret generic cloudflare-api-token \
  --from-literal=api-token=$CF_API_TOKEN --dry-run=client -o yaml | kubectl apply -f -

# 1) Усі addons разом:
k3ctl apply -c iwis-ai

# Або окремо:
k3ctl apply -c iwis-ai cert-manager
k3ctl apply -c iwis-ai ingress         # → HAProxy з ingress-haproxy-values.yaml
k3ctl apply -c iwis-ai external-dns
k3ctl apply -c iwis-ai monitoring

kubectl -n ingress-haproxy get svc     # EXTERNAL-IP має бути 10.0.40.100
arping -c2 10.0.40.100                 # з будь-якої ноди в 10.0.40.0/24
```

### 4.9. ClusterIssuer Let's Encrypt (вручну)

`k3ctl apply ... cert-manager` ставить лише Helm chart. ClusterIssuer-и (ресурси cert-manager CRD) — вручну:

```bash
kubectl apply -f - <<'EOF'
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata: { name: letsencrypt-prod }
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@iwis.dev
    privateKeySecretRef: { name: letsencrypt-prod }
    solvers:
      - dns01:
          cloudflare:
            apiTokenSecretRef:
              name: cloudflare-api-token
              key: api-token
EOF
```

### 4.10. (зарезервовано)
Об’єднано з §4.8 — ExternalDNS розгортається через `k3ctl apply -c iwis-ai external-dns`.

### 4.11. Cloudflare Tunnel (адмін-сервіси)
```bash
# 1) У Cloudflare Zero Trust створити tunnel iwis-ai-admin, отримати TOKEN
kubectl create ns cloudflared
kubectl -n cloudflared create secret generic cloudflared-token \
  --from-literal=token=$CLOUDFLARE_TUNNEL_TOKEN
kubectl apply -f ./variables/iwis-ai/cloudflared.yaml
kubectl -n cloudflared get pods

# 2) У Cloudflare Zero Trust → Tunnels → iwis-ai-admin → Public Hostname:
#   longhorn.admin.iwis.dev -> http://longhorn-frontend.longhorn-system.svc.cluster.local:80
#   hubble.admin.iwis.dev   -> http://hubble-ui.kube-system.svc.cluster.local:80
```

### 4.12. Smoke-test
```bash
kubectl run web --image=nginx --port=80
kubectl expose pod web --port=80
kubectl create ingress web --class=haproxy \
  --rule="hello.iwis.dev/*=web:80,tls=hello-iwis-tls" \
  --annotation cert-manager.io/cluster-issuer=letsencrypt-prod
curl -kI https://hello.iwis.dev   # через DNAT 116.202.72.52
```

---

## 5. Масштабування до 3 master (HA etcd)

1. Додати master-2 / master-3 у `variables/iwis-ai.yaml` (розкоментувати).
2. На обох виконати `scripts/prepare-node.sh` (§3) та додати `config.yaml` (§9.3, без `cluster-init`, з правильним `node-ip`).
3. Приєднати (k3ctl додасть лише нові master за статусом):
   ```bash
   k3ctl cluster create iwis-ai -c ./variables/iwis-ai.yaml
   ```
4. Перевірити: `kubectl get nodes`, `kubectl -n kube-system get pod -l name=kube-vip -o wide`, `etcdctl member list`.
5. У `cilium-values.yaml` `operator.replicas: 2` → `helm upgrade`.

---

## 6. Бекап і обслуговування

- **etcd** — авто-snapshot (налаштовано у `/etc/rancher/k3s/config.yaml`: `etcd-snapshot-schedule-cron: "0 */6 * * *"`, `etcd-snapshot-retention: 10`): `/var/lib/rancher/k3s/server/db/snapshots/`.
  Для S3: додати у `config.yaml`: `etcd-s3: true`, `etcd-s3-bucket: ...`, `etcd-s3-endpoint: ...`.
- **Longhorn** — recurring backup до S3 (Hetzner Object Storage / Backblaze).
- **Velero** — окремо через `addons.backup` (за потреби).
- **Моніторинг** — `addons.monitoring` (grafana-agent → Prometheus/Loki).

---

## 7. Чек-ліст готовності

- [ ] Усі ноди `Ready`, версії k3s збігаються.
- [ ] `cilium status` зелений; `kubeProxyReplacement: True`.
- [ ] `kubectl -n kube-system get pod -l name=kube-vip -o wide` running на master-1; `ping 10.0.40.20` ОК.
- [ ] `kubectl -n ingress-haproxy get svc` → EXTERNAL-IP `10.0.40.100`; `arping` з 10.0.40.0/24 проходить.
- [ ] DNS `*.iwis.dev → 116.202.72.52`, DNAT працює; `curl https://hello.iwis.dev` → 200.
- [ ] cert-manager видає LE, ExternalDNS створює A-записи в Cloudflare.
- [ ] Longhorn: 3 репліки на worker-1..3, lb/master без дисків.
- [ ] etcd snapshot файл створено.
- [ ] Cloudflare Tunnel up, `longhorn.admin.iwis.dev` доступний.

---

## 8. Корисні команди

```bash
# Видалити ноду з кластера
./k3ctl node delete <node> -c iwis-ai

# Перевірити VIP
ip -br a | grep 10.0.40.20

# Cilium debug
cilium status
cilium hubble enable
cilium hubble ui   # local proxy

# Longhorn UI
kubectl -n longhorn-system port-forward svc/longhorn-frontend 8080:80

# Видалити kube-config locally (без зміни поточного контексту)
./k3ctl kubeconfig get iwis-ai --kubeconfig-switch-context=false
```

## 9. Troubleshooting

- **Нода не приєднується (passwd mismatch)**:
  ```bash
  kubectl -n kube-system delete secret <node>.node-password.k3s
  ```
- **MTU проблеми (фрагментація через VSwitch)**: перевірити `MTU: 1400` у Cilium values, `ping -M do -s 1372`.
- **VIP не пінгується**: перевірити `vip_interface` у kube-vip-cp.yaml, ARP-таблицю на сусідніх нодах.
- **HAProxy LB не отримує EXTERNAL-IP**: перевірити Kube-VIP DaemonSet (`svc_enable=true`) та ConfigMap `kubevip` (`range-global`).
- **API через VIP: `x509: certificate is valid for... not 10.0.40.20`** (див. §9.2 нижче).
- **DNAT не працює**: tcpdump на lb-1 `:80,:443`, перевірити SNAT/маршрутизацію на шлюзі.

### 9.2. Відсутній `10.0.40.20` у TLS SANs API-сервера (баг k3ctl + etcd-snapshot)

**Симптом:**
```
kubectl --server=https://10.0.40.20:6443 get nodes
Unable to connect... x509: certificate is valid for 10.0.40.10, 10.43.0.1, ... not 10.0.40.20
```

**Причина:** `k3ctl` генерує `/etc/systemd/system/k3s.service` із розбитим cron schedule на окремі аргументи:
```
'--etcd-snapshot-schedule-cron=0' \
'*/6' \
'*' \
'*' \
'*' \
```
Це ламає pflag parsing — всі `--tls-san` після `*/6` ігноруються API-сервером.

**Рішення (на master-1 як root):**
```bash
# 1. Перевірити порядок tls-san у service файлі
systemctl cat k3s.service | grep -A2 -B2 'tls-san'

# 2. Якщо tls-san після '*/6' — перенести ДО нього:
sed -i "s|'--tls-san=116.202.72.52'|'--tls-san=116.202.72.52' \\\
        '--tls-san=10.0.40.20' \\\
        '--tls-san=api.iwis.dev'|" /etc/systemd/system/k3s.service

# 3. Перезавантажити systemd
systemctl daemon-reload

# 4. Видалити старий cert API-сервера (згенерується новий з правильними SANs)
rm /var/lib/rancher/k3s/server/tls/serving-kube-apiserver.crt \
   /var/lib/rancher/k3s/server/tls/serving-kube-apiserver.key

# 5. Перезапустити k3s
systemctl restart k3s
sleep 35

# 6. Перевірити SANs
echo | openssl s_client -connect 10.0.40.20:6443 2>/dev/null \
  | openssl x509 -noout -text | grep -A3 'Subject Alternative'
# Має бути: IP Address:10.0.40.20

# 7. Перевірити kubectl через VIP
kubectl --server=https://10.0.40.20:6443 get nodes
```

**Довгострокове виправлення (вже застосовано у k3ctl):** аргументи зі spaces у value (зокрема `etcd-snapshot-schedule-cron`) автоматично пишуться у `/etc/rancher/k3s/config.yaml`, а не в `INSTALL_K3S_EXEC`.

### 9.3. Перенесення конфігурації з k3s.service у config.yaml (best practice)

K3s підтримує `/etc/rancher/k3s/config.yaml` — це чистіший спосіб ніж CLI flags у service файлі.

**На кожному master (як root):**
```bash
# 1. Записати всі параметри у config.yaml
cat > /etc/rancher/k3s/config.yaml << 'EOF'
disable:
  - servicelb
  - traefik
service-cidr: "10.43.0.0/16"
cluster-cidr: "10.42.0.0/16"
cluster-domain: "cluster.local"
flannel-backend: "none"
secrets-encryption: true
disable-kube-proxy: true
disable-network-policy: true
cluster-init: true
tls-san:
  - "116.202.72.52"
  - "10.0.40.20"
  - "api.iwis.dev"
  - "10.0.40.10"
node-ip: "10.0.40.10"          # замінити на IP ноди
advertise-address: "10.0.40.10" # замінити на IP ноди
etcd-snapshot-schedule-cron: "0 */6 * * *"
etcd-snapshot-retention: 10
EOF

# 2. Спростити ExecStart у service файлі
cat > /etc/systemd/system/k3s.service << 'SVCEOF'
[Unit]
Description=Lightweight Kubernetes
Documentation=https://k3s.io
Wants=network-online.target
After=network-online.target

[Install]
WantedBy=multi-user.target

[Service]
Type=notify
EnvironmentFile=-/etc/default/%N
EnvironmentFile=-/etc/sysconfig/%N
EnvironmentFile=-/etc/systemd/system/k3s.service.env
KillMode=process
Delegate=yes
User=root
LimitNOFILE=1048576
LimitNPROC=infinity
LimitCORE=infinity
TasksMax=infinity
TimeoutStartSec=0
Restart=always
RestartSec=5s
ExecStartPre=-/sbin/modprobe br_netfilter
ExecStartPre=-/sbin/modprobe overlay
ExecStart=/usr/local/bin/k3s server
SVCEOF

# 3. Перезапустити
systemctl daemon-reload
systemctl restart k3s
sleep 35

# 4. Перевірити
systemctl is-active k3s
echo | openssl s_client -connect 10.0.40.20:6443 2>/dev/null \
  | openssl x509 -noout -text | grep -A2 'Subject Alternative'
```

> **Для master-2, master-3**: `node-ip` та `advertise-address` замінити відповідно на `10.0.40.11`, `10.0.40.12`; `cluster-init` видалити — вони приєднуються через `server: https://10.0.40.20:6443`.
