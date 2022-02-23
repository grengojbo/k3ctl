package types

// const HelmListCommand = "helm list -A -o json"
const HelmListCommand = "helm list -A -o json --kubeconfig %s"
const HelmDeleteCommand = "helm del %s -n %s --kubeconfig %s"

var HelmAddons = []string{"cert-manager"}

const IngressDefaultName = "ingress-nginx"

const NginxHelmRepo = "ingress-nginx/ingress-nginx"
const NginxHelmURL = "https://kubernetes.github.io/ingress-nginx"
const NginxDefaultName = "ingress-nginx"
const NginxDefaultNamespace = "ingress-nginx"
const NginxGetSvcCommand = "kubectl get svc ingress-nginx-controller -n ingress-nginx"

const (
	MetalLBVersion           = "v0.12.1"
	MetalLBNamespaceManifest = "https://raw.githubusercontent.com/metallb/metallb/%s/manifests/namespace.yaml"
	MetalLBManifest          = "https://raw.githubusercontent.com/metallb/metallb/%s/manifests/metallb.yaml"
)
