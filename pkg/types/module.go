package types

// const HelmListCommand = "helm list -A -o json"
const HelmListCommand = "helm list -A -o json --kubeconfig %s"
const HelmDeleteCommand = "helm del %s -n %s --kubeconfig %s"

var HelmAddons = []string{"cert-manager"}

const IngressDefaultName = "ingress-nginx"

const DefaultAwsRegion = "eu-central-1"

// TODO: add Velero
const VeleroHelmRepoName = "vmware-tanzu"
const VeleroHelmRepo = "vmware-tanzu/velero"
const VeleroHelmURL = "https://vmware-tanzu.github.io/helm-charts"
const VeleroDefaultName = "velero"
const VeleroDefaultNamespace = "velero"

const BackupDefaultProvider = "aws"
const BackupDefaultName = "velero"

// const BackupDefault = ""
// const BackupDefault = ""
// const BackupDefault = ""
// const BackupDefault = ""

const CertManagerHelmRepoName = "jetstack"
const CertManagerHelmRepo = "jetstack/cert-manager"
const CertManagerHelmURL = "https://charts.jetstack.io"
const CertManagerDefaultName = "cert-manager"
const CertManagerDefaultNamespace = "cert-manager"

const NginxHelmRepoNane = "ingress-nginx"
const NginxHelmRepo = "ingress-nginx/ingress-nginx"
const NginxHelmURL = "https://kubernetes.github.io/ingress-nginx"
const NginxDefaultName = "nginx"
const NginxDefaultNamespace = "ingress-nginx"
const NginxGetSvcCommand = "kubectl get svc ingress-nginx-controller -n ingress-nginx"

const HaproxyHelmRepoName = "haproxytech"
const HaproxyHelmRepo = "haproxytech/kubernetes-ingress"
const HaproxyHelmURL = "https://haproxytech.github.io/helm-charts"
const HaproxyDefaultNamespace = "haproxy-controller"
const HaproxyDefaultName = "haproxy"

const MonitoringDefaultNamespace = "monitoring"

const GrafanaAgentCloudHelmRepoName = "grengojbo"
const GrafanaAgentCloudHelmRepo = "grengojbo/grafana-agent-cloud"
const GrafanaAgentCloudHelmURL = "https://grengojbo.github.io/charts/"
const GrafanaAgentCloudDefaultNamespace = "monitoring"
const GrafanaAgentCloudDefaultName = "grafana-agent"

const (
	MetalLBVersion           = "v0.12.1"
	MetalLBNamespaceManifest = "https://raw.githubusercontent.com/metallb/metallb/%s/manifests/namespace.yaml"
	MetalLBManifest          = "https://raw.githubusercontent.com/metallb/metallb/%s/manifests/metallb.yaml"
)
