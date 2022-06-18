package types

var LoadBalancerList = []string{MetalLb, ServiceLb}
var AddonsList = []string{"cert-manager", "ingress", "external-dns", "monitoring", "backup"}
var PresetList = []string{PresetSingle, PresetOneMaster, PresetWorker, PresetMultyMaster}

// var IngressControllers = []string{IngressAmbassador, IngressAmbassadorAPI, IngressContour, IngressHaproxyName, IngressHaproxy, IngressKing, IngressNginx, IngressTraefik}
var IngressList = []string{IngressHaproxy, IngressNginx}

// const HelmListCommand = "helm list -A -o json"
const HelmRepoListCommand = "helm repo list --kubeconfig %s --kube-context %s -o json | jq -r '.[].name'"
const HelmListCommand = "helm list -A -o json --kubeconfig %s --kube-context %s"
const HelmRepoAddCommand = "helm repo add %s %s --kubeconfig %s --kube-context %s"
const HelmRepoUpdateCommand = "helm repo update --kubeconfig %s --kube-context %s"
const HelmDeleteCommand = "helm delete %s -n %s --kubeconfig %s --kube-context %s"

const SecretCreateCommand = "kubectl create secret %s %s %s --kubeconfig %s --context %s"
const SecretDeleteCommand = "kubectl delete secret %s -n %s --kubeconfig %s --context %s"
const SecretListCommand = "kubectl get secret -n %s --kubeconfig %s --context %s -o json | jq -r '.items[].metadata.name'"

const NamespaceGetCommand = "kubectl get ns -o=custom-columns='NAME:.metadata.name' --no-headers --kubeconfig %s --context %s"
const NamespaceCreateCommand = "kubectl create namespace %s --kubeconfig %s --context %s"

var HelmAddons = []string{"cert-manager"}

const IngressDefaultName = "ingress-nginx"

const DefaultAwsRegion = "eu-central-1"

// TODO: add Velero
const VeleroHelmRepoName = "vmware-tanzu"
const VeleroHelmRepo = "vmware-tanzu/velero"
const VeleroHelmURL = "https://vmware-tanzu.github.io/helm-charts"
const VeleroDefaultName = "velero"
const VeleroDefaultNamespace = "velero"
const VeleroPluginAwsImage = "velero/velero-plugin-for-aws:v1.4.1"

const BackupDefaultProvider = "aws"
const BackupDefaultName = "velero"
const BackupDefaultNamespace = "backup"

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

// https://www.haproxy.com/documentation/kubernetes/latest/configuration/custom-resources/backend/install-backend-crd/
const HaproxyCrdBackend = "https://cdn.haproxy.com/documentation/kubernetes/1.8/crd/backend.yaml"

// https://www.haproxy.com/documentation/kubernetes/latest/configuration/custom-resources/defaults/install-defaults-crd/
const HaproxyCrdDefaults = "https://cdn.haproxy.com/documentation/kubernetes/1.8/crd/defaults.yaml"

// https://www.haproxy.com/documentation/kubernetes/latest/configuration/custom-resources/global/install-global-crd/
const HaproxyCrdGlobal = "https://cdn.haproxy.com/documentation/kubernetes/1.8/crd/global.yaml"

// const HaproxyCrd = ""

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
