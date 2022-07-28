#!/bin/bash

PROJECT_NAME=k3ctl
GITHUB_ACCOUNT=grengojbo
APP_CHANNEL=latest
GITHUB_URL=https://github.com/${GITHUB_ACCOUNT}/${PROJECT_NAME}/releases
STORAGE_URL=STORAGE_URL=https://storage.googleapis.com/k3s-ci-builds

# https://github.com/grengojbo/k3ctl/releases/download/v0.1.1/${APP_ARHIVE_NAME}
# pause
# echo "Finding download url"
# downloadUrl=$(echo ${downloadPage} | grep -Po '(?<=href=")[^"]*(?=")' | grep hub-linux-amd64)
# echo ${downloadUrl}
# echo "Getting installation file"
# wget "https://github.com${downloadUrl}" --output-document=github-hub.tar.gz
# echo "Untarring..."
# tar -zxvf github-hub.tar.gz
# rm github-hub.tar.gz
# cd hub-linux*
# ls -la
# echo "Installing"
# sudo bash install
# cd ..
# rm -rf hub-linux*

DOWNLOADER=

# --- helper functions for logs ---
info()
{
  echo '[INFO] ' "$@"
}
warn()
{
  echo '[WARN] ' "$@" >&2
}
fatal()
{
  echo '[ERROR] ' "$@" >&2
  exit 1
}

GetOS() {
  OS=`uname -s`
  case "$OS" in
    Darwin) OS="darwin" ;;
    Linux) OS="linux" ;;
    FreeBSD) OS="FreeBSD" ;;
  esac
}

# --- set arch and suffix, fatal if architecture not supported ---
setup_verify_arch() {
    if [ -z "$ARCH" ]; then
        ARCH=$(uname -m)
    fi
    case $ARCH in
        amd64)
            ARCH=amd64
            SUFFIX=
            ;;
        x86_64)
            ARCH=amd64
            SUFFIX=
            ;;
        arm64)
            ARCH=arm64
            SUFFIX=-${ARCH}
            ;;
        s390x)
            ARCH=s390x
            SUFFIX=-${ARCH}
            ;;
        aarch64)
            ARCH=arm64
            SUFFIX=-${ARCH}
            ;;
        arm*)
            ARCH=arm
            SUFFIX=-${ARCH}hf
            ;;
        *)
            fatal "Unsupported architecture $ARCH"
    esac
}

# --- verify existence of network downloader executable ---
verify_downloader() {
    # Return failure if it doesn't exist or is no executable
    [ -x "$(command -v $1)" ] || return 1

    # Set verified executable as our downloader program and return success
    DOWNLOADER=$1
    return 0
}

# --- create temporary directory and cleanup when done ---
setup_tmp() {
    TMP_DIR=$(mktemp -d -t k3ctl-install.XXXXXXXXXX)
    TMP_HASH=${TMP_DIR}/k3ctl.hash
    TMP_BIN=${TMP_DIR}/k3ctl.bin
    cleanup() {
        code=$?
        set +e
        trap - EXIT
        rm -rf ${TMP_DIR}
        exit $code
    }
    trap cleanup INT EXIT
}

# --- use desired k3s version if defined or find version from channel ---
get_release_version() {
	echo "Getting download page for latest release"
 	info "Finding release for channel ${APP_CHANNEL}"
    #     version_url="${INSTALL_K3S_CHANNEL_URL}/${INSTALL_K3S_CHANNEL}"
  case $DOWNLOADER in
    curl)
			APP_VERSION=$(curl -L -s https://api.github.com/repos/${GITHUB_ACCOUNT}/${PROJECT_NAME}/releases/${APP_CHANNEL} | jq -r .tag_name)
      # VERSION_K3S=$(curl -w '%{url_effective}' -L -s -S ${version_url} -o /dev/null | sed -e 's|.*/||')
      ;;
    wget)
			APP_VERSION=$( wget -q -O - https://api.github.com/repos/${GITHUB_ACCOUNT}/${PROJECT_NAME}/releases/${APP_CHANNEL} | jq -r .tag_name)
      # VERSION_K3S=$(wget -SqO /dev/null ${version_url} 2>&1 | grep -i Location | sed -e 's|.*/||')
      ;;
    *)
      fatal "Incorrect downloader executable '$DOWNLOADER'"
      ;;
  esac
	APP_ARHIVE_NAME=${PROJECT_NAME}_${OS}_${ARCH}_${APP_VERSION}
	info "Using ${APP_VERSION} as release"
}

# --- download from github url ---
download() {
  [ $# -eq 2 ] || fatal 'download needs exactly 2 arguments'

  case $DOWNLOADER in
    curl)
      curl -o $1 -sfL $2
      ;;
    wget)
      wget -qO $1 $2
      ;;
    *)
      fatal "Incorrect executable '$DOWNLOADER'"
      ;;
  esac

  # Abort if download command failed
  [ $? -eq 0 ] || fatal 'Download failed'
}

# --- download hash from github url ---
download_hash() {
    # if [ -n "${INSTALL_K3S_COMMIT}" ]; then
    #     HASH_URL=${STORAGE_URL}/k3s${SUFFIX}-${INSTALL_K3S_COMMIT}.sha256sum
    # else
    #     HASH_URL=${GITHUB_URL}/download/${VERSION_K3S}/sha256sum-${ARCH}.txt
    # fi
  HASH_URL=${GITHUB_URL}/download/${APP_VERSION}/checksums.txt
    info "Downloading hash ${HASH_URL}"
    download ${TMP_HASH} ${HASH_URL}
		# cat ${TMP_HASH}
    HASH_EXPECTED=$(grep " ${APP_ARHIVE_NAME}" ${TMP_HASH})
    HASH_EXPECTED=${HASH_EXPECTED%%[[:blank:]]*}
		# echo "HASH_EXPECTED: ${HASH_EXPECTED}"
}

# --- check hash against installed version ---
installed_hash_matches() {
  if [ -x ${BIN_DIR}/${PROJECT_NAME} ]; then
    SHASUM_BIN=$(command -v sha256sum 2> /dev/null)
    if [[ "${SHASUM_BIN}" == "" ]]; then
      SHASUM_BIN=$(command -v shasum 2> /dev/null)
      if [[ "${SHASUM_BIN}" == "" ]]; then
        fatal "sha256sum or shasum not found"
      else
        HASH_INSTALLED=$(shasum -a 256 ${BIN_DIR}/${PROJECT_NAME})
      fi
    else
      HASH_INSTALLED=$(sha256sum ${BIN_DIR}/${PROJECT_NAME})
    fi
    HASH_INSTALLED=${HASH_INSTALLED%%[[:blank:]]*}
    if [ "${HASH_EXPECTED}" = "${HASH_INSTALLED}" ]; then
      return
    fi
  fi
  return 1
}

# --- download binary from github url ---
download_binary() {
  # if [ -n "${INSTALL_K3S_COMMIT}" ]; then
  #   BIN_URL=${STORAGE_URL}/k3s${SUFFIX}-${INSTALL_K3S_COMMIT}
  # else
  #   BIN_URL=${GITHUB_URL}/download/${VERSION_K3S}/k3s${SUFFIX}
  # fi
  BIN_URL=${GITHUB_URL}/download/${APP_VERSION}/${APP_ARHIVE_NAME}
  info "Downloading binary ${BIN_URL}"
  download ${TMP_BIN} ${BIN_URL}
}

# --- verify downloaded binary hash ---
verify_binary() {
  info "Verifying binary download"
  SHASUM_BIN=$(command -v sha256sum 2> /dev/null)
  if [[ "${SHASUM_BIN}" == "" ]]; then
    SHASUM_BIN=$(command -v shasum 2> /dev/null)
    if [[ "${SHASUM_BIN}" == "" ]]; then
      fatal "sha256sum or shasum not found"
    else
      HASH_BIN=$(shasum -a 256 ${TMP_BIN})
    fi
  else
    HASH_BIN=$(sha256sum ${TMP_BIN})
  fi
  # SHASUM_BIN=$(sha256sum ${TMP_BIN} | awk '{print $1}')
  HASH_BIN=${HASH_BIN%%[[:blank:]]*}
  if [ "${HASH_EXPECTED}" != "${HASH_BIN}" ]; then
    fatal "Download sha256 does not match ${HASH_EXPECTED}, got ${HASH_BIN}"
  fi
}

# --- setup permissions and move binary to system directory ---
setup_binary() {
	# echo "TMP_BIN: ${TMP_BIN}"
	# echo "Untarring..."
	# tar -zxvf ${TMP_BIN}
	# ls -l ${TMP_BIN}
	chmod 755 ${TMP_BIN}
  info "Installing ${PROJECT_NAME} to ${BIN_DIR}/${PROJECT_NAME}"
  $SUDO chown root:root ${TMP_BIN}
  $SUDO mv -f ${TMP_BIN} ${BIN_DIR}/${PROJECT_NAME}
}

# --- download and verify k3s ---
download_and_verify() {
  setup_verify_arch
	GetOS
	echo "ARCH: ${ARCH} OS: ${OS}"
  verify_downloader curl || verify_downloader wget || fatal 'Can not find curl or wget for downloading files'
  setup_tmp
  get_release_version
  download_hash

  if installed_hash_matches; then
    info 'Skipping binary downloaded, installed k3ctl matches hash'
    return
  fi

  download_binary
  verify_binary
  setup_binary
}

# --- define needed environment variables ---
setup_env() {
	# --- use sudo if we are not already root ---
	SUDO=sudo
  if [ $(id -u) -eq 0 ]; then
    SUDO=
  fi

	# --- use /usr/local/bin if root can write to it, otherwise use /opt/bin if it exists
	BIN_DIR=/usr/local/bin
  if ! $SUDO sh -c "touch ${BIN_DIR}/k3s-ro-test && rm -rf ${BIN_DIR}/k3s-ro-test"; then
    if [ -d /opt/bin ]; then
      BIN_DIR=/opt/bin
    fi
  fi

	# --- setup channel values
  #INSTALL_K3S_CHANNEL_URL=${INSTALL_K3S_CHANNEL_URL:-'https://update.k3s.io/v1-release/channels'}
  #INSTALL_K3S_CHANNEL=${INSTALL_K3S_CHANNEL:-'stable'}

}

InstallArkade() {
  echo "install arkade..."
  rm -rf ${HOME}/.arkade
  
  verify_downloader curl || verify_downloader wget || fatal 'Can not find curl or wget for downloading files'
  
  case $DOWNLOADER in
    curl)
      curl -SLfs https://get.arkade.dev | $SUDO sh
      ;;
    wget)
      wget -q -O - https://get.arkade.dev | $SUDO sh
      ;;
    *)
      fatal "Incorrect downloader executable '$DOWNLOADER'"
      ;;
  esac
  

  # arkade --help
  # ark --help  # a handy alias

  # Windows users with Git Bash
  # curl -sLS https://get.arkade.dev | sh

  #arkade get porter
  #$SUDO mv ${HOME}/.arkade/bin/porter ${BIN_DIR}/

  arkade get jq
  $SUDO mv ${HOME}/.arkade/bin/jq ${BIN_DIR}/

  arkade get yq
  $SUDO mv ${HOME}/.arkade/bin/yq ${BIN_DIR}/

  #arkade get k3d
  #$SUDO mv ${HOME}/.arkade/bin/k3d ${BIN_DIR}/

  arkade get kubectl
  $SUDO mv ${HOME}/.arkade/bin/kubectl ${BIN_DIR}/
  
  arkade get kustomize
  $SUDO mv ${HOME}/.arkade/bin/kustomize ${BIN_DIR}/

  # HELM_VERSION=`curl https://api.github.com/repos/helm/helm/releases -s | grep -Po '"tag_name": "\K.*?(?=")' -m1`
  case $DOWNLOADER in
    curl)
      HELM_VERSION=`curl -L -s https://api.github.com/repos/helm/helm/releases/latest | jq -r .tag_name`
      ;;
    wget)
      HELM_VERSION=`wget -q -O - https://api.github.com/repos/helm/helm/releases/latest | jq -r .tag_name`
      ;;
    *)
      fatal "Incorrect downloader executable '$DOWNLOADER'"
      ;;
  esac

  arkade get helm -v ${HELM_VERSION} 
  $SUDO mv ${HOME}/.arkade/bin/helm ${BIN_DIR}/

  arkade get krew
  $SUDO mv ${HOME}/.arkade/bin/krew ${BIN_DIR}/

  arkade get stern
  $SUDO mv ${HOME}/.arkade/bin/stern ${BIN_DIR}/

}

InstallVelero() {
  # VELERO_VERSION=v1.8.1
  # VELERO_VERSION=`curl https://api.github.com/repos/vmware-tanzu/velero/releases -s | grep -Po '"tag_name": "\K.*?(?=")' -m1`
  case $DOWNLOADER in
    curl)
      VELERO_VERSION=`curl -L -s https://api.github.com/repos/vmware-tanzu/velero/releases/latest | jq -r .tag_name`
      curl https://github.com/vmware-tanzu/velero/releases/download/${VELERO_VERSION}/velero-${VELERO_VERSION}-${OS}-${ARCH}.tar.gz -o ./velero-${VELERO_VERSION}-${OS}-${ARCH}.tar.gz
      ;;
    wget)
      VELERO_VERSION=`wget -q -O - https://api.github.com/repos/vmware-tanzu/velero/releases/latest | jq -r .tag_name`
      wget https://github.com/vmware-tanzu/velero/releases/download/${VELERO_VERSION}/velero-${VELERO_VERSION}-${OS}-${ARCH}.tar.gz
      ;;
    *)
      fatal "Incorrect downloader executable '$DOWNLOADER'"
      ;;
  esac
  tar -xzvf velero-${VELERO_VERSION}-${OS}-${ARCH}.tar.gz
  rm -f velero-${VELERO_VERSION}-${OS}-${ARCH}.tar.gz
  chmod 0777 velero-${VELERO_VERSION}-${OS}-${ARCH}/velero
  $SUDO mv velero-${VELERO_VERSION}-${OS}-${ARCH}/velero /usr/local/bin/velero
}


setup_env
InstallArkade
download_and_verify
InstallVelero