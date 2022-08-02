package userdata

import (
	"fmt"
)

var DefaultDockerComposeVersion = "1.29.2"

type UserData struct {
	Podman        bool
	PodmanRoot    bool
	Docker        bool
	DockerCompose bool
	Provider      string
	KvmAgent      bool
	VmWareAgent   bool
	K3sMaster     bool
	K3sAgent      bool
	UserName      string
	GitLabRunner  bool
	TimeZone      string
	Locale        string
	PackageMinmal bool
}

var preData = `#!/usr/bin/env bash

SERVER_LOCALE=%s
TZ_INSTALL=%s
kvm=%v
vmware=%v
USER_NAME=%s

INSTALL_DOCKER=%v
INSTALL_COMPOSER=%v
DOCKER_COMPOSE_VERSION=%s

INSTALL_PODMAN=%v
PODMAN_ROOT=%v
INSTALL_GITLAB_RUNNER=%v
`

var postData = `if [[ ! -z ${INSTALL_DOCKER} ]]; then
if [[ "${INSTALL_DOCKER}" == "true" ]]; then
	echo "START install Docker"
	# UpdateGrub
	InstallDocker
	echo "Finish install Docker"
	echo " "
fi
fi

if [[ "${INSTALL_PODMAN}" == "true" ]]; then
# echo "START install Podman"
# UpdateGrub
InstallPodman
# echo "Finish install Podman"
if [[ "${PODMAN_ROOT}" == "true" ]]; then
	echo "enable podman socket"
	# systemctl enable --now podman.socket
	# systemctl start podman.socket
	# ln -s /var/run/podman/podman.sock /var/run/docker.sock
fi
	
	# mkdir -p ~/bin
	# curl -sSL https://github.com/rootless-containers/rootlesskit/releases/download/v0.14.2/rootlesskit-$(uname -m).tar.gz | tar Cxzv ~/bin
	
	# echo "ROUTING PING PACKETS"
	# echo "net.ipv4.ping_group_range=0   2147483647" > /etc/sysctl.d/ping_group_range.conf
	
	# echo 'kernel.unprivileged_userns_clone=1' > /etc/sysctl.d/userns.conf
	# echo "user.max_user_namespaces=28633" >> /etc/sysctl.d/userns.conf
	# [ -f /proc/sys/kernel/unprivileged_userns_clone ] && sudo sh -c 'echo "kernel.unprivileged_userns_clone=1" >> /etc/sysctl.d/userns.conf'
	# sudo sysctl --system
	# /etc/cni/net.d/87-podman-bridge.conflist

	# https://kubernetes.io/docs/tasks/configure-pod-container/translate-compose-kubernetes/
	# curl -L https://github.com/kubernetes/kompose/releases/download/v1.22.0/kompose-linux-amd64 -o kompose

echo " "
fi

if [[ ! -z ${INSTALL_COMPOSER} ]]; then
if [[ "${INSTALL_COMPOSER}" == "true" ]]; then
	echo "START install Docker composer"
	InstallDockerCompose
	echo "Finish install Docker composer"
	echo " "
fi
fi

if [[ ! -z ${INSTALL_GITLAB_RUNNER} ]]; then
if [[ "${INSTALL_GITLAB_RUNNER}" == "true" ]]; then
	echo "START install Docker Gitlab Runner"
	echo " "
	InstallGitlabRunner
fi
fi

echo "Закройте терминал или выйдие из ssh сессии"
if [[ "${INSTALL_PODMAN}" == "true" ]]; then
echo "Затем подключитесь опять и проверте работает ли Podman"
echo "[RUN] podman info"
fi

%s
`

var ubuntuData = `OS=$(awk '/DISTRIB_ID=/' /etc/*-release | sed 's/DISTRIB_ID=//' | tr '[:upper:]' '[:lower:]')
UpdateGrub() {
  sudo sed -i 's/GRUB_CMDLINE_LINUX=""/GRUB_CMDLINE_LINUX="cgroup_enable=memory swapaccount=1"/' /etc/default/grub
  sudo update-grub
  sudo update-grub2
  cat << EOF > /tmp/kubespray-br_netfilter.conf
br_netfilter
EOF
  #sudo mv /tmp/kubespray-br_netfilter.conf /etc/modules-load.d/kubespray-br_netfilter.conf
#   sudo cat << EOF > /etc/modules-load.d/kube_proxy-ipvs.conf
# ip_vs
# ip_vs_rr
# ip_vs_wrr
# ip_vs_sh
# nf_conntrack_ipv4
# EOF
  # modprobe nf_conntrack_ipv4
  sudo modprobe br_netfilter
}

InstallVmAgent() {
  echo "ставим пакеты"

  if [[ ${vmware} == true ]]; then
    echo "Install vmWare tools"
    apt install -y open-vm-tools
  fi

  if [[ ${kvm} == true ]]; then
    echo "Install kvm agent"
    apt install -y qemu-guest-agent
  fi
}

InstallPackageMinimal() {
  echo "обновляем систему минимум пакетов"
  apt update
  apt upgrade -y
  # aptitude install -y open-vm-tools
  apt install -y git sudo iftop curl bzip2 dnsutils keychain jq ntp mc vim socat libseccomp2 make
  apt install -y util-linux wget ca-certificates iotop mtr ipvsadm lsof lvm2 net-tools
  
  if [ $OS == "debian" ]; then
    apt install -y firmware-linux-free 
  fi
  
  apt autoremove -y
}

InstallPackage() {
  echo "обновляем систему"
  apt update
  apt upgrade -y
  apt install -y git sudo iftop curl bzip2 dnsutils hdparm keychain parted jq
  apt install -y build-essential libssl-dev libcap2-bin gcc g++ make
  apt install -y ntp mc neovim vim neovim locales-all socat libseccomp2 jq tree htop tmux unzip tar apt-transport-https
  apt install -y util-linux wget ca-certificates iotop mtr ipvsadm usbmount pmount lsof

  if [ $OS == "debian" ]; then
    apt install -y firmware-linux-free 
  fi
  
  apt install -y lvm2 net-tools && apt autoremove -y
}

RemoveRsyslog() {
  echo "Stop rsyslog to clean up logs"
  systemctl stop rsyslog
  systemctl disable rsyslog
  apt purge -y rsyslog

}

RemoveSnap() {
  echo "Stop and remove snap"
  # snap list 
  snap remove lxd
  snap remove core18
  snap remove snapd
  apt purge -y snapd
  apt autoremove -y
}

CreateNotRootUser() {
  local NEW_USER=${1:-user}
  local NEW_UID=${2:-1001}
  local NEW_SHELL=${3:-/sbin/nologin}
  local ADD_HOME_DIR=""

  if [[ "${NEW_SHELL}" != "/sbin/nologin" ]]; then
    ADD_HOME_DIR=" -d /home/${NEW_USER} -m"
  fi

  groupadd -g ${NEW_UID} ${NEW_USER}
  useradd -u ${NEW_UID} -r -g ${NEW_UID} -s ${NEW_SHELL}${ADD_HOME_DIR} -c "Default Application User" ${NEW_USER}
}

CleanupUbuntu() {
  echo "Clean Ubuntu..."

  # apt update -y && apt upgrade -y && apt install -y lvm2 net-tools && apt autoremove -y

  # Cleanup all logs
  #cat /dev/null > /var/log/audit/audit.log
  cat /dev/null > /var/log/wtmp
  cat /dev/null > /var/log/lastlog

  #cleanup persistent udev rules
  rm -rf /etc/udev/rules.d/70-persistent-net.rules

  #cleanup /tmp directories
  rm -rf /tmp/*
  rm -rf /var/tmp/*

  #cleanup current ssh keys
  # rm -f /etc/ssh/ssh_host_*
  # sed -i -e 's|exit 0||' /etc/rc.local
  # sed -i -e 's|.*test -f /etc/ssh/ssh_host_dsa_key.*||' /etc/rc.local
  # bash -c 'echo "test -f /etc/ssh/ssh_host_dsa_key || dpkg-reconfigure openssh-server" >> /etc/rc.local'
  # bash -c 'echo "exit 0" >> /etc/rc.local'

  # Clear hostname
  cat /dev/null > /etc/hostname

  # Cleanup apt
  apt autoremove -y
  apt-get -y clean

  #cleanup shell history
  #history -w
  #history -c
}

# AddSshKey /root/.ssh ubuntu
AddSshKey() {
  local USER_FROM=${1}
  local USER_TO=${2}
  echo "Copy authorized_keys"
  mkdir -p /home/${USER_TO}/.ssh
  chmod 0700 /home/${USER_TO}/.ssh
  chown ${USER_TO}:${USER_TO} /home/${USER_TO}/.ssh
  
  cp ${USER_FROM}/authorized_keys /home/${USER_TO}/.ssh/authorized_keys
  chmod 0600 /home/${USER_TO}/.ssh/authorized_keys
  chown ${USER_TO}:${USER_TO} /home/${USER_TO}/.ssh/authorized_keys
}

InstallDocker () {
	sudo apt remove -y docker docker.io
	sudo apt update -y
	sudo apt install -y apt-transport-https ca-certificates wget software-properties-common curl gnupg2
	
	echo "Ubuntu"
	curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo apt-key add -
	sudo add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"

	echo "доступные пакеты"
	sudo apt update -y && sudo apt-cache policy docker-ce

	sudo apt install -y docker-ce docker-ce-cli containerd.io
	sudo usermod -aG docker $USER_NAME
}

InstallPodman() {
  echo "Start intall podman"
  . /etc/os-release
  sudo sh -c "echo 'deb http://download.opensuse.org/repositories/devel:/kubic:/libcontainers:/stable/xUbuntu_${VERSION_ID}/ /' > /etc/apt/sources.list.d/devel:kubic:libcontainers:stable.list"
  wget -nv https://download.opensuse.org/repositories/devel:kubic:libcontainers:stable/xUbuntu_${VERSION_ID}/Release.key -O- | sudo apt-key add -
  sudo apt-get update -qq
  sudo apt-get -qq --yes install libapparmor-dev podman buildah skopeo slirp4netns uidmap fuse-overlayfs containernetworking-plugins dbus-user-session
}

InstallDockerCompose () {
  echo "Install docker compose"
  sudo curl -L "https://github.com/docker/compose/releases/download/${DOCKER_COMPOSE_VERSION}/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
  sudo chmod 0777 /usr/local/bin/docker-compose
  sudo ln -s /usr/local/bin/docker-compose /usr/bin/docker-compose
}

InstallGitlabRunner () {
  echo "Install Gitlab Runner"
  curl -L https://packages.gitlab.com/install/repositories/runner/gitlab-runner/script.deb.sh | sudo bash
  sudo apt install -y gitlab-runner
  sudo usermod -aG docker gitlab-runne
}

%s

%s
`

// NewUserData - Prepare for UserData
//   provider (aws, azure, hetzner, proxmox, vmware)
//   mode (k3sMaster, k3sWorker, docker, podman)
func NewUserData(mode string, provider string) (result *UserData) {
	result = &UserData{}

	result.Provider = provider
	result.Locale = "en_US.UTF-8"
	result.TimeZone = "Etc/UTC"
	result.UserName = "ubuntu"
	result.GitLabRunner = false
	result.PackageMinmal = false

	if result.Provider == "proxmox" {
		result.KvmAgent = true
		result.VmWareAgent = false
	} else if result.Provider == "vmware" {
		result.KvmAgent = false
		result.VmWareAgent = true
	} else {
		result.KvmAgent = false
		result.VmWareAgent = false
	}

	if mode == "docker" {
		result.Podman = true
		result.PodmanRoot = false
		result.Docker = true
		result.DockerCompose = true
		result.K3sMaster = false
		result.K3sAgent = false
	} else if mode == "podman" {
		result.Podman = true
		result.PodmanRoot = true
		result.Docker = false
		result.DockerCompose = false
		result.K3sMaster = false
		result.K3sAgent = false
	}
	return result
}

// GetUserData - Get UserData string
func (u *UserData) GetUserData() (res string) {
	resPreData := fmt.Sprintf(preData, u.Locale, u.TimeZone, u.KvmAgent, u.VmWareAgent, u.UserName, u.Docker, u.DockerCompose, DefaultDockerComposeVersion, u.Podman, u.PodmanRoot, u.GitLabRunner)
	resPostData := fmt.Sprintf(postData, "reboot;")

	strPackage := "InstallPackage"
	if u.PackageMinmal {
		strPackage = "InstallPackageMinimal"
	}
	str := ""
	if u.Provider == "vmware" || u.Provider == "proxmox" {
		str = fmt.Sprintf("%s\n%s", str, "InstallVmAgent")
	}
	// RemoveRsyslog
	// RemoveSnap
	resData := fmt.Sprintf(ubuntuData, str, strPackage)

	res = fmt.Sprintf("%s\n%s\n%s", resPreData, resData, resPostData)
	// res = fmt.Sprintf("%s\n", resData)
	return res
}
