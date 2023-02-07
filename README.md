# k3ctl management k3s kubernetes clusters


Download and install

```bash
curl -sfL https://raw.githubusercontent.com/grengojbo/k3ctl/main/install.sh | sh -
```

```bash
export AWS_ACCESS_KEY_ID=<YOUR_ACCESS_KEY_ID>
export AWS_SECRET_ACCESS_KEY=<YOUR_SECRET_ACCESS_KEY>
```


```bash
export ARM_CLIENT_ID="WWWWWWWW-WWWW-WWWW-WWWW-WWWWWWWWWWWW" && \
export ARM_CLIENT_SECRET="XXXXXXXX-XXXX-XXXX-XXXX-XXXXXXXXXXXX" && \
export ARM_TENANT_ID="YYYYYYYY-YYYY-YYYY-YYYY-YYYYYYYYYYYY" && \
export ARM_SUBSCRIPTION_ID="ZZZZZZZZ-ZZZZ-ZZZZ-ZZZZ-ZZZZZZZZZZZZ"
```

```bash
export HCLOUD_TOKEN=XXXXXXXXXXXXXX
```

или .env

```bash
DB_PASSWORD=XXX
```

Example create cluster

```bash
./k3ctl cluster create sample
./k3ctl apply -c iwisops
```

## For developers

если bastion.Name == "local" то выполняем команду локально

Для работы через бастион необходимо указать

``` yaml
  bastions:
    - name: mybastion
      user: noname
      address: 192.168.0.2
      # sshPort: 2222
      # sshAuthorizedKey: ~/keys/my_rsa
  nodes:
    - name: k3-master
      user: nonameTwo
      bastion: mybastion
```

Неизменять текущий контест

```bash
./k3ctl kubeconfig get sample --kubeconfig-switch-context=false
```

Удаление ноды из кластера

```bash
./k3ctl node delete <node name> -c <cluster name> 
```

### Gitlab Agent



```bash
helm repo add gitlab https://charts.gitlab.io
helm repo update
helm upgrade --install gitlab-agent gitlab/gitlab-agent \
    --namespace gitlab-agent \
    --create-namespace \
    --set config.kasAddress=wss://gitlab.iwis.com.ua/-/kubernetes-agent/
    # --set image.tag=v15.2.0 \
    # --set config.token=<your_token> \
    
```


## проблемы

При подключении узла

```bash
kubectl -n kube-system delete secrets <agent-node-name>.node-password.rke2
kubectl -n kube-system delete secrets <agent-node-name>.node-password.k3s
```


```bash
/var/lib/rancher/k3s/server/cred/passwd newer than datastore and could cause a cluster outage. Remove the file(s) from disk and restart to be recreated from datastore.
```

```bash

```