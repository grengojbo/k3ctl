# k3ctl
management k3s kubernetes clusters

Example create cluster

```bash
./k3ctl cluster create -c sample --dry-run --verbose
```

## For developers

```bash
kubeadm token generate
kubeadm certs certificate-key
```

a