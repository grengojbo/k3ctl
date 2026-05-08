# Addon: ingress

Встановлює Ingress Controller через Helm. Підтримувані провайдери: `haproxy` (рекомендовано), `nginx`.

## Конфігурація (`variables/<cluster>.yaml`)

```yaml
spec:
  addons:
    ingress:
      name: haproxy          # haproxy | nginx  (обов'язково)
      # namespace: haproxy-controller  # default для haproxy
      # version: 1.49.0               # optional — latest if omitted
      # disabled: false

      valuesFile: ./variables/<cluster>/ingress-haproxy-values.yaml

      # values: (inline Helm overrides, опційно)
      #   key: value

      # hostMode: false  # true = hostNetwork (не рекомендується з Cilium)
```

> **Примітка**: `name` визначає який controller встановлювати. Значення відповідають константам `IngressHaproxy="haproxy"` та `IngressNginx="nginx"` (`pkg/types/types.go`).

## Провайдери

### `haproxy` — HAProxy Kubernetes Ingress

- **Helm repo**: `haproxytech/kubernetes-ingress` (`https://haproxytech.github.io/helm-charts`)
- **Default namespace**: `haproxy-controller`
- **Default release name**: `haproxy`
- **Статистика**: `http://<LB-IP>:1024/` (HAProxy stats page)

### `nginx` — ingress-nginx

- **Helm repo**: `ingress-nginx/ingress-nginx` (`https://kubernetes.github.io/ingress-nginx`)
- **Default namespace**: `ingress-nginx`

## Потік виконання `k3ctl apply ingress`

```
k3ctl apply -c <cluster> ingress
  1. helm repo add haproxytech https://haproxytech.github.io/helm-charts
  2. helm upgrade --install haproxy haproxytech/kubernetes-ingress \
       -n haproxy-controller [--values <valuesFile>]
```

## Налаштування для Cilium + kube-vip (iwis-ai)

Кластер `iwis-ai` використовує Cilium з `kube-proxy-replacement=true` (tunnel mode) та kube-vip для Service VIP.

**Ключові параметри** (`variables/iwis-ai/ingress-haproxy-values.yaml`):

```yaml
controller:
  kind: DaemonSet
  daemonset:
    useHostNetwork: false   # НЕ потрібен — Cilium eBPF сам маршрутизує трафік
    useHostPort: false      # НЕ потрібен — конфліктує з Cilium при hostNetwork

  service:
    type: LoadBalancer
    annotations:
      kube-vip.io/loadbalancerIPs: "10.0.40.100"   # VIP від kube-vip
    externalTrafficPolicy: Cluster   # Cluster — обов'язково з Cilium tunnel mode

  nodeSelector:
    role: lb               # тільки на lb-нодах

  tolerations:
    - key: dedicated
      operator: Equal
      value: lb
      effect: NoSchedule   # lb-ноди мають taint dedicated=lb:NoSchedule

defaultBackend:
  enabled: false           # HAProxy має вбудований fallback
```

> **Чому `hostNetwork: false`?**
> При `hostNetwork: true` `hostPort` мусить збігатися з `containerPort`. HAProxy chart встановлює різні значення — це призводить до помилки валідації DaemonSet. З Cilium `kube-proxy-replacement=true` `hostNetwork` не потрібен — Cilium перехоплює трафік до Service VIP через eBPF на рівні ноди.

> **Чому `externalTrafficPolicy: Cluster`?**
> З Cilium у tunnel mode та `kube-proxy-replacement=true` значення `Local` може некоректно розподіляти трафік — Cilium сам балансує між podами через eBPF незалежно від topology. `Cluster` гарантує доступність з будь-якої ноди.

## Підготовка нод (lb-ноди)

lb-ноди повинні мати лейбл `role=lb` та taint `dedicated=lb:NoSchedule`.

Лейбл виставляється автоматично через `spec.nodes[].labels` у `variables/<cluster>.yaml`:
```yaml
nodes:
  - name: lb-1
    labels:
      role: "lb"
```

Taint виставляється вручну (k3ctl не підтримує taints):
```bash
kubectl taint node lb-1 lb-2 dedicated=lb:NoSchedule --overwrite
```

## Перевірка

```bash
# Поди DaemonSet (мають бути на lb-нодах)
kubectl -n haproxy-controller get pods -o wide

# Service — EXTERNAL-IP має бути VIP від kube-vip
kubectl -n haproxy-controller get svc

# Статистика HAProxy
curl http://10.0.40.100:1024/

# ARP-перевірка VIP з будь-якої ноди кластера
arping -c2 10.0.40.100
```

## Приклад Ingress ресурсу

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: my-app
  annotations:
    kubernetes.io/ingress.class: haproxy
spec:
  ingressClassName: haproxy
  rules:
    - host: my-app.iwis.dev
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: my-app
                port:
                  number: 80
  tls:
    - hosts:
        - my-app.iwis.dev
      secretName: my-app-tls   # Certificate від cert-manager
```
