Розгорання кластера k3s

Домен iwis.dev
  - master-1 ip: 10.0.40.10/24
  - worker-1 ip: 10.0.40.101/24
  - worker-2 ip: 10.0.40.102/24
	- worker-3 ip: 10.0.40.103/24
	- lb-1 ip: 10.0.40.98/24
	- lb-2 ip: 10.0.40.99/24


k3ctl не покриває (треба руками):

Cilium (CNI)
Kube-VIP (control-plane static-pod + cloud-provider)
Longhorn
Cloudflare Tunnel
ClusterIssuer (Helm chart cert-manager ставиться, але самі CRD-ресурси ClusterIssuer — вручну)
