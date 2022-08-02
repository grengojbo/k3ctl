#!/bin/bash

[ ${UID} -ne 0 ] && {
  echo "RUN only ROOT user"
  exit 1;
}

kubeconfig_console=true
kubeconfig_ssm=false
# kubeconfig_ssm=true

CLUSTER_NAME=bbox
# APP_PROVIDER=aws
APP_PROVIDER=k3s
AWS_SSM=false
APISERVER_EXTERNAL="noname.bbox.kiev.ua"
K3S_SERVER_LB=true
export INSTALL_K3S_CHANNEL='v1.19'

export K3S_TOKEN=d46tgsretgseSAFszfdsfhfgd543q3 
export K3S_PROVIDER_ID=$(curl -s http://169.254.169.254/latest/meta-data/placement/availability-zone)/$(curl -s http://169.254.169.254/latest/meta-data/instance-id)
export K3S_KUBECONFIG_MODE=644
export K3S_NODE_NAME=$(hostname -f)

export K3S_PUBLIC_IP=$(curl -s http://169.254.169.254/latest/meta-data/public-ipv4)

UbuntuUpdate() {
  apt update -y && apt upgrade -y
}

UbuntuIstallPackage() {
  apt install -y awscli
  [ "${AWS_SSM}" = "true" ] && {
    echo "Install amazon-ssm-agent"
    # wget https://s3.amazonaws.com/ec2-downloads-windows/SSMAgent/latest/debian_amd64/amazon-ssm-agent.deb
    # dpkg -i ./amazon-ssm-agent.deb
    # rm ./amazon-ssm-agent.deb
  }
}

UbuntuClean() {
  apt autoremove -y
  apt-get -y clean
}

RemoveRsyslog() {
  echo "Stop rsyslog to clean up logs"
  systemctl stop rsyslog
  systemctl disable rsyslog
  apt purge -y rsyslog

}

RemoveSnap() {
  echo "Stop and remove snap"
  #snap list
  snap remove lxd
  snap remove core18
  snap remove snapd
  apt purge -y snapd
}

CreateNotRootUser() {
  local NEW_USER=${1:-user}
  local NEW_UID=${2:-1001}
  groupadd -g ${NEW_UID} ${NEW_USER}
  useradd -u ${NEW_UID} -r -g ${NEW_UID} -s /sbin/nologin -c "Default Application User" ${NEW_USER}
  usermod -aG ${NEW_UID} ${USER}
  #groupadd ${NEW_USER}
  #useradd -r -u 1001 -g 1001 ${NEW_USER}
  #usermod -aG ${NEW_USER} ${USER_NAME}
  ## usermod -aG docker ${NEW_USER} # скорее всего ненадо
  ## sudo usermod -aG docker user
  ## sudo usermod -aG user baz
}

#RemoveRsyslog
#RemoveSnap

#UbuntuUpdate
#UbuntuIstallPackage
#UbuntuClean
#CreateNotRootUser

#curl -sfL https://get.k3s.io | K3S_TOKEN=d46tgsretgseSAFszfdsfhfgd543q3 INSTALL_K3S_EXEC='server --cluster-init --tls-san 3.124.193.26 --tls-san bbox-master.bbox.kiev.ua --tls-san ip-172-18-1-209 --disable servicelb --disable traefik --secrets-encryption --disable-cloud-controller --node-name ip-172-18-1-209 --kubelet-arg cloud-provider=external --write-kubeconfig-mode 644 --kubelet-arg provider-id=aws:///eu-central-1a/i-0139a8a312d3e67f0' INSTALL_K3S_CHANNEL='v1.19' sh -

InstallK3sServer() {

  K3S_ARGS=" "
  if [[ -z ${K3S_SERVER_LB} ]]; then
    # echo "Add ${K3S_SERVER_LB}"
    K3S_ARGS="${K3S_ARGS} --disable servicelb"
  fi

  if [[ "${INGRESS_NAME}" != "traefik" ]]; then
    # echo "Add ${K3S_SERVER_LB}"
    K3S_ARGS="${K3S_ARGS} --disable traefik"
  fi

  [ "${APP_PROVIDER}" = "aws" ] && { 
    K3S_ARGS="${K3S_ARGS} --disable-cloud-controller --kubelet-arg cloud-provider=external --kubelet-arg provider-id=aws:///${K3S_PROVIDER_ID}"
  
    K3S_ARGS="${K3S_ARGS} --kube-apiserver-arg cloud-provider=external"
    K3S_ARGS="${K3S_ARGS} --kube-apiserver-arg allow-privileged=true"
    K3S_ARGS="${K3S_ARGS} --kube-apiserver-arg feature-gates=CSINodeInfo=true,CSIDriverRegistry=true,CSIBlockVolume=true,VolumeSnapshotDataSource=true"
    K3S_ARGS="${K3S_ARGS} --kube-controller-arg cloud-provider=external"
    K3S_ARGS="${K3S_ARGS} --kubelet-arg feature-gates=CSINodeInfo=true,CSIDriverRegistry=true,CSIBlockVolume=true"
  }

  ADD_TLS_SAN="--tls-san ${K3S_PUBLIC_IP}  --tls-san ${K3S_NODE_NAME} "
  [ ! -z ${APISERVER_EXTERNAL} ] && ADD_TLS_SAN="${ADD_TLS_SAN} --tls-san ${APISERVER_EXTERNAL} "
  export INSTALL_K3S_EXEC="server ${ADD_TLS_SAN}${K3S_ARGS}"
  echo ${INSTALL_K3S_EXEC}
  curl -sfL https://get.k3s.io | sh -
  # curl -sfL https://get.k3s.io | INSTALL_K3S_EXEC="server --tls-san ${K3S_PUBLIC_IP} --tls-san bbox-master.bbox.kiev.ua --tls-san ${K3S_NODE_NAME} --disable servicelb --disable traefik --secrets-encryption --disable-cloud-controller --kubelet-arg cloud-provider=external --kubelet-arg provider-id=aws:///${K3S_PROVIDER_ID}${K3S_ARGS}" INSTALL_K3S_CHANNEL='v1.19' sh -

  # Wait for k3s.yaml to appear
  echo -n "Wait /etc/rancher/k3s/k3s.yaml sleep 5sec"	
  while ! sudo [ -e /etc/rancher/k3s/k3s.yaml ]
  do
    echo -n "."	
    sleep 5
  done
  echo " "

  echo -n "Wait /var/lib/rancher/k3s/server/manifests sleep 5sec"	
  while ! sudo [ -e  /var/lib/rancher/k3s/server/manifests ]
  do
    echo -n "."	
    sleep 5
  done
  echo " "

  # curl -sfL h https://raw.githubusercontent.com/kmcgrath/k3s-terraform-modules/master/manifests/cloud-provider-aws.yaml > cloud-provider-aws.yaml
  # mv cloud-provider-aws.yaml /var/lib/rancher/k3s/server/manifests/

  #kubectl apply -f https://raw.githubusercontent.com/kubernetes/cloud-provider-aws/master/manifests/rbac.yaml
  #kubectl apply -f https://raw.githubusercontent.com/kubernetes/cloud-provider-aws/master/manifests/aws-cloud-controller-manager-daemonset.yaml

  #ProviderID=$(kubectl describe node $(hostname -f) | grep ProviderID)
  # Allow root to have superuser over the cluster:
  cp /etc/rancher/k3s/k3s.yaml ~root/.kube/config

} 

ShowStatus() {
  # Wait for all kube-system deployments to roll out
  for d in $(sudo kubectl get deploy -n kube-system --no-headers -o name)
  do
    sudo kubectl -n kube-system rollout status $d
  done
}

InstallK3sServer
ShowStatus

# Display kubeconfig to console.
#[ ${inst-id} -eq 0 -a "${kubeconfig-console}" = "true" ] && {
[ "${kubeconfig_console}" = "true" ] && {
  ss=/etc/cron.daily/kube-config-to-console.sh

  cat > $ss << EOF
#!/bin/bash
kc=~root/.kube/config
[ -r \$kc ] && {
  echo ===================================================================== > /dev/console
  echo Cluster kubeconfig: > /dev/console
  echo ===================================================================== > /dev/console
  cat \$kc > /dev/console
  echo ===================================================================== > /dev/console
  echo End cluster kubeconfig > /dev/console
}
EOF

  chmod 755 $ss
  $ss
}

# Save kubeconfig as an ssm parameter.
#[ ${inst-id} -eq 0 -a "${kubeconfig-ssm}" = "true" ] && {
[ "${kubeconfig_ssm}" = "true" ] && {
  aws ssm put-parameter --name ${CLUSTER_NAME}-kubeconfig --value "`cat ~root/.kube/config`" --type String
}

# nginx-ingress install. To replace traefik which broke (was deployed as a daemonset but overtime changed to deployment).
#sudo helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
#sudo helm repo update

# Force ingress-nginx install as a daemonset
#sudo helm install ingress-nginx ingress-nginx/ingress-nginx \
#  -n ingress-nginx \
#  --create-namespace \
#  --set kind=DaemonSet 

# finally how does k3s look?
sudo kubectl get nodes
sudo kubectl get all -A

kubectl describe node $(hostname -f) | grep ProviderID
