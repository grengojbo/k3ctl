# cert-manager addon

cert-manager встановлюється через `k3ctl apply -c <cluster> cert-manager` і налаштовується у секції `spec.addons.certManager` конфіг-файлу кластера.

## Конфігурація (`variables/<cluster>.yaml`)

```yaml
spec:
  addons:
    certManager:
      # name: cert-manager           # default
      # namespace: cert-manager      # default
      # version: v1.16.3             # optional — latest if omitted
      # disabled: false              # set true to skip install

      provider: cloudflare           # http | cloudflare | route53
      email: admin@example.com       # ACME email для Let's Encrypt

      # valuesFile: ./manifests/<cluster>/cert-manager-values.yaml  # Helm values (опційно)

      # manifests: (опційно — якщо потрібен кастомний ClusterIssuer замість автогенерованого)
      #   - ./manifests/<cluster>/my-clusterissuer.yaml
```

> **Примітка**:
> - `provider` → k3ctl **автоматично генерує** ClusterIssuer `letsencrypt-prod` без жодних файлів.
> - `manifests` → список файлів для `kubectl apply -f` (перебиває автогенерацію якщо вказано).
> - `valuesFile` → Helm values (`helm upgrade --values`). Не плутати з `manifests`.

### Провайдери

- `http` — HTTP-01 (не рекомендується для production)
- `cloudflare` — DNS-01 через Cloudflare API Token (потребує `CF_API_TOKEN` у `.env`; якщо не встановлено — ClusterIssuer **не** буде застосовано, виводиться ERROR)
- `route53` — DNS-01 через AWS Route53 (потребує `AWS_REGION` та IAM прав)

## Потік виконання `k3ctl apply cert-manager`

```
k3ctl apply -c <cluster> cert-manager
  1. LoadDotEnv variables/<cluster>/.env  → CF_API_TOKEN → create secret cloudflare-api-token
  2. helm repo add jetstack https://charts.jetstack.io
  3. helm upgrade --install cert-manager jetstack/cert-manager \
       -n cert-manager --set installCRDs=true [--values <valuesFile>]
  4a. якщо manifests[] → kubectl apply -f <manifests[i]>  (кастомний ClusterIssuer)
  4b. якщо provider=cloudflare|route53 → kubectl apply -f - (автогенерований ClusterIssuer)
  4c. якщо provider=http або не вказано → ClusterIssuer не створюється
```

## Передумова: секрет Cloudflare API Token

`k3ctl apply -c <cluster> cert-manager` **автоматично** читає `CF_API_TOKEN` в такому порядку пріоритету:

1. `variables/<cluster>/.env` — **кластерний** `.env` (найвищий пріоритет)
2. `.env` — глобальний у поточній директорії (fallback)
3. Змінна оточення shell (`export CF_API_TOKEN=...`) — перебиває `.env` якщо вже встановлена в shell

Якщо `CF_API_TOKEN` знайдено — `k3ctl` сам створить namespace `cert-manager` і секрет `cloudflare-api-token` **перед** Helm install.

**Рекомендований спосіб — кластерний `.env`:**
```bash
# variables/iwis-ai/.env  (вже у .gitignore)
CF_API_TOKEN=your-cloudflare-zone-dns-edit-token
```

Тоді достатньо:
```bash
k3ctl apply -c iwis-ai cert-manager
```

**Або через shell export:**
```bash
export CF_API_TOKEN=your-cloudflare-zone-dns-edit-token
k3ctl apply -c iwis-ai cert-manager
```

> Якщо `CF_API_TOKEN` не встановлено — секрет **не** буде створено (Helm install відбудеться). Секрет можна створити вручну:
> ```bash
> kubectl -n cert-manager create secret generic cloudflare-api-token \
>   --from-literal=api-token=<TOKEN> --dry-run=client -o yaml | kubectl apply -f -
> ```

## ClusterIssuer (Cloudflare DNS-01)

Файл: `manifests/<cluster>/clusterissuer-cloudflare.yaml`

```yaml
apiVersion: cert-manager.io/v1
kind: ClusterIssuer
metadata:
  name: letsencrypt-prod
spec:
  acme:
    server: https://acme-v02.api.letsencrypt.org/directory
    email: admin@example.com
    privateKeySecretRef:
      name: letsencrypt-prod
    solvers:
      - dns01:
          cloudflare:
            apiTokenSecretRef:
              name: cloudflare-api-token
              key: api-token
```

## Встановлення

```bash
# 1. Створити секрет (якщо ще не зроблено)
export CF_API_TOKEN=...
kubectl create ns cert-manager 2>/dev/null || true
kubectl -n cert-manager create secret generic cloudflare-api-token \
  --from-literal=api-token=$CF_API_TOKEN --dry-run=client -o yaml | kubectl apply -f -

# 2. Встановити cert-manager + ClusterIssuer
k3ctl apply -c <cluster> cert-manager
```

## Перевірка

```bash
# Поди cert-manager
kubectl -n cert-manager get pods

# Статус ClusterIssuer
kubectl get clusterissuer letsencrypt-prod
kubectl get clusterissuer letsencrypt-prod -o jsonpath='{.status.conditions[0].message}'
# Очікувано: "The ACME account was registered with the ACME server"

# Тест — створити тестовий Certificate
kubectl apply -f - <<'EOF'
apiVersion: cert-manager.io/v1
kind: Certificate
metadata:
  name: test-cert
  namespace: default
spec:
  secretName: test-cert-tls
  issuerRef:
    name: letsencrypt-prod
    kind: ClusterIssuer
  dnsNames:
    - test.example.com
EOF
kubectl describe certificate test-cert
kubectl delete certificate test-cert
```

## Troubleshooting

```bash
# Логи cert-manager
kubectl -n cert-manager logs -l app=cert-manager --tail=50

# Статус Challenge (DNS-01)
kubectl get challenge -A
kubectl describe challenge -A

# Перевірити що секрет існує
kubectl -n cert-manager get secret cloudflare-api-token
```

## Підтримувані поля `CertManager` struct

| Поле | Тип | Опис |
|---|---|---|
| `name` | string | Helm release name (default: `cert-manager`) |
| `namespace` | string | Namespace (default: `cert-manager`) |
| `version` | string | Версія Helm chart (optional) |
| `disabled` | bool | Пропустити встановлення |
| `valuesFile` | string | Шлях до Helm values YAML файлу |
| `values` | map[string]string | Helm `--set` overrides |
| `manifests` | []string | Список YAML-файлів для `kubectl apply` після install |
| `repo` | HelmRepo | Кастомний Helm repo (override за замовчуванням `jetstack`) |
