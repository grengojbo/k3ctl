# k3ctl
management k3s kubernetes clusters

Example create cluster

```bash
./k3ctl cluster create sample --dry-run --verbosee
ssh login@ip -C "sudo cat /etc/rancher/k3s/k3s.yaml" > ~/.kube/cluster-name.yaml
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

и в се

```bash
kubeadm token generate
kubeadm certs certificate-key

sudo cat /var/lib/rancher/k3s/server/node-token
sudo kubectl get nodes --selector='node-role.kubernetes.io/worker' -o json | jq -r '.items[].status.nodeInfo.kubeletVersion' | sort -u | tr '+' '-'

```

```bash
curl -sfL https://get.k3s.io | K3S_URL='https://<IP>:6443' K3S_TOKEN='<TOKEN>' INSTALL_K3S_CHANNEL='v1.23' sh -s -
```

```bash
WORKER_NODE=$(kubectl get node -o jsonpath='{.items[*].metadata.name}' --selector='!node-role.kubernetes.io/master')
for n in ${WORKER_NODE}
  do
    echo "kubectl label node ${n} node-role.kubernetes.io/worker=worker"
  done
kubectl get nodes -o wide
```




```bash
./k3ctl cluster create --verbose -c sample -h
 4326  ./k3ctl cluster create -c sample --k3s-version v1.19
 4327  ./k3ctl cluster create --verbose -c sample --k3s-version v1.19
 4370  ./k3ctl cluster create --verbose -c sample
 4388  ./k3ctl cluster create -c sample -h
 4389  ./k3ctl cluster create -c sample --trace
 4396  ./k3ctl cluster create -c sample --dry-run --debug
 4397  ./k3ctl cluster create -c sample --dry-run --verbose --no-lb=true
 4398  ./k3ctl cluster create -c sample --dry-run --verbose --no-lb
 4400  ./k3ctl cluster create -c sample --dry-run --verbose --no-ingress
 4402  ./k3ctl cluster create -c sample --dry-run -h
 4403  ./k3ctl cluster create -c sample --dry-run --secrets-encryption
 4404  ./k3ctl cluster create -c sample --dry-run --rootless
 4405  ./k3ctl cluster create -c sample --dry-run --selinux
 4406  sudo ./k3ctl cluster create -c sample
 4409  ./k3ctl cluster create -c sample --dry-run
 4419  sudo ./k3ctl cluster create -c sample --verbose
 4420  ./k3ctl cluster create -c sample --verbose
 4421  ./k3ctl cluster create -c sample
 4575  ./k3ctl cluster create -c sample --dry-run --verbose
 4580  ./k3ctl cluster create --dry-run --verbose
 4581  ./k3ctl cluster create --dry-run --verbose sample
 4582  ./k3ctl cluster create sample --dry-run --verbose aaa
 4583  ./k3ctl cluster create --dry-run --verbose aaa
 4584  ./k3ctl cluster create sample --dry-run
 4585  git clone git@github.com:grengojbo/k3ctl.git
 4586  cd k3ctl
 4595  ./k3ctl cluster create aaa --dry-run --verbose
 4596  ./k3ctl cluster create samle --dry-run --verbose
 4597  ./k3ctl cluster create sample --dry-run --verbose
 4598  ./k3ctl cluster create developer --dry-run --verbose
 4599  ./k3ctl cluster create developer --dry-run
 4600  ./k3ctl cluster create aaa --dry-run
 4601  ./k3ctl cluster create noname --dry-run
 4602  ./k3ctl cluster create demo --dry-run
 4603  ./k3ctl cluster create demo -h
```