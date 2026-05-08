# Addon: external-dns

Встановлює [ExternalDNS](https://github.com/kubernetes-sigs/external-dns) через Helm (bitnami chart). Автоматично синхронізує DNS-записи у зовнішньому провайдері (Cloudflare, Route53, тощо) з `Ingress` та `Service` ресурсами кластера.

## Конфігурація (`variables/<cluster>.yaml`)

```yaml
spec:
  loadBalancer:
    domain: iwis.dev          # (обов'язково) базовий домен — додається до domainFilters
    # externalIP: 1.2.3.4    # зовнішній IP для A-записів (опційно)

  addons:
    externalDns:
      provider: cloudflare    # cloudflare | aws | azure | hetzner | native (обов'язково)
      # namespace: kube-system          # default
      # version: 8.3.4                  # optional — latest if omitted
      # disabled: false

      # Домени для фільтрації (крім spec.loadBalancer.domain, який додається автоматично)
      domains:
        - iwis.dev

      # Для AWS — фільтрація по Hosted Zone ID (опційно)
      # HostedZoneIds:
      #   - Z1234567890

      # Для AWS — регіон (якщо не вказано — береться з spec.providers.aws.region)
      # region: eu-central-1

      values:
        logLevel: info          # debug | info | warning | error

      # valuesFile: ./variables/<cluster>/external-dns-values.yaml  # Helm values (опційно)
```

> **Увага**: якщо `spec.loadBalancer.domain` або `externalDns.provider` не встановлені — ExternalDNS **не** буде встановлено (виводиться WARN).

## Провайдери

| `provider` | Опис | Потрібні credentials |
|---|---|---|
| `cloudflare` | Cloudflare DNS API | `CF_API_TOKEN` або `CF_API_KEY` + `CF_API_EMAIL` |
| `aws` | AWS Route53 | IAM credentials або instance role |
| `azure` | Azure DNS | `azure.json` secret |
| `hetzner` | Hetzner DNS | `HETZNER_TOKEN` |
| `native` | Вбудований (для тестів) | — |

### Cloudflare

`k3ctl apply -c <cluster> external-dns` **автоматично** читає `CF_API_TOKEN` в такому порядку пріоритету:

1. `variables/<cluster>/.env` — **кластерний** `.env` (найвищий пріоритет)
2. `.env` — глобальний у поточній директорії (fallback)
3. Змінна оточення shell (`export CF_API_TOKEN=...`) — якщо вже встановлена в shell

Якщо `CF_API_TOKEN` знайдено — передається як `cloudflare.apiToken` у Helm override.
Якщо не знайдено — виводиться WARN і ExternalDNS встановлюється **без** токена (DNS синхронізація не працюватиме).

**Рекомендований спосіб — кластерний `.env`:**
```bash
# variables/<cluster>/.env  (вже у .gitignore)
CF_API_TOKEN=your-cloudflare-zone-dns-edit-token
```

`cloudflare.proxied` контролюється полем `proxied` у конфігурації (default: `false` — DNS-only, обов'язково для cert-manager DNS-01 challenge):

```yaml
# variables/<cluster>.yaml
addons:
  externalDns:
    provider: cloudflare
    proxied: false   # false = DNS-only (default) | true = CF proxy (CDN/DDoS)
```

### AWS Route53

Якщо існує файл `variables/<cluster>/aws-credentials` — k3ctl автоматично створить secret `exrernal-dns-<cluster>-aws-creds` з ключем `credentials` і передасть його в Helm через `aws.credentials.secretName`.

```yaml
addons:
  externalDns:
    provider: aws
    region: eu-central-1      # або через spec.providers.aws.region
    HostedZoneIds:
      - Z1234567890ABC
```

## Потік виконання `k3ctl apply external-dns`

```
k3ctl apply -c <cluster> external-dns
  1. Перевірка: spec.loadBalancer.domain та externalDns.provider мусять бути встановлені
  2. (тільки install, не update) crd.create=true — встановлює DNSEndpoint CRD
  3. (тільки install, AWS) створює secret з AWS credentials якщо є variables/<cluster>/aws-credentials
  4. helm repo add bitnami https://charts.bitnami.com/bitnami
  5. helm upgrade --install external-dns bitnami/external-dns -n kube-system \
       --set provider=<provider> \
       --set domainFilters[0]=<loadBalancer.domain> [--set domainFilters[1]=...] \
       --set txtOwnerId=<clusterName> \
       --set policy=sync \
       --set triggerLoopOnEvent=true \
       --set interval=5m \
       [--set metrics.enabled=true якщо monitoring увімкнено] \
       [--values <valuesFile>]
```

## Автоматичні Helm overrides

| Override | Значення | Опис |
|---|---|---|
| `provider` | `cloudflare` / `aws` / … | з `externalDns.provider` |
| `domainFilters[i]` | `iwis.dev`, … | `spec.loadBalancer.domain` + `externalDns.domains[]` |
| `txtOwnerId` | `<clusterName>` | ідентифікатор TXT ownership record |
| `policy` | `sync` | синхронізація (видаляє застарілі записи) |
| `triggerLoopOnEvent` | `true` | реакція на зміни в реальному часі |
| `interval` | `5m` | інтервал повної синхронізації |
| `crd.create` | `true` | (тільки install) створює DNSEndpoint CRD |
| `metrics.enabled` | `true` | якщо `monitoring.disabled=false` |
| `aws.region` | з config | тільки для AWS |
| `aws.zoneType` | `public` | тільки для AWS |

## Перевірка

```bash
# Поди
kubectl -n kube-system get pods -l app.kubernetes.io/name=external-dns

# Логи (синхронізація)
kubectl -n kube-system logs -l app.kubernetes.io/name=external-dns --tail=50

# Перевірити створені DNS-записи (для Cloudflare — через API або UI)
# TXT ownership records мають вигляд: "heritage=external-dns,external-dns/owner=<clusterName>"
```

## Примітки

- `policy: sync` — ExternalDNS **видаляє** DNS-записи якщо відповідний `Ingress`/`Service` видалено. Для безпечнішого режиму: `policy: upsert-only` (не видаляє).
- `txtOwnerId` встановлено в `<clusterName>` — дозволяє кількох ExternalDNS у різних кластерах керувати одним доменом без конфліктів.
- `cloudflare.proxied: false` рекомендується для `cert-manager` DNS-01 challenge — Cloudflare proxy може блокувати ACME validation.
