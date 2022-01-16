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
```

a