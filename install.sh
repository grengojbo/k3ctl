#!/bin/bash

PROJECT_NAME=k3ctl
GITHUB_ACCOUNT=grengojbo
APP_CHANNEL=latest
GITHUB_URL=https://github.com/${GITHUB_ACCOUNT}/${PROJECT_NAME}/releases

APP_ARHIVE_NAME=k3ctl_0.1.1_linux_amd64.tar.gz
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
    # if [ -n "${INSTALL_K3S_COMMIT}" ]; then
    #     VERSION_K3S="commit ${INSTALL_K3S_COMMIT}"
    # elif [ -n "${INSTALL_K3S_VERSION}" ]; then
    #     VERSION_K3S=${INSTALL_K3S_VERSION}
    # else
  # info "Finding release for channel ${INSTALL_K3S_CHANNEL}"
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
    # fi
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
  HASH_URL=${GITHUB_URL}/download/${APP_VERSION}/sha256sum-${ARCH}.txt
	https://github.com/grengojbo/k3ctl/releases/download/v0.1.1/${PROJECT_NAME}_0.1.1_checksums.txt
    # info "Downloading hash ${HASH_URL}"
    # download ${TMP_HASH} ${HASH_URL}
    # HASH_EXPECTED=$(grep " k3s${SUFFIX}$" ${TMP_HASH})
    # HASH_EXPECTED=${HASH_EXPECTED%%[[:blank:]]*}
}
# --- download and verify k3s ---
download_and_verify() {
  setup_verify_arch
	echo "ARCH: ${ARCH}"
  verify_downloader curl1 || verify_downloader wget || fatal 'Can not find curl or wget for downloading files'
  setup_tmp
  get_release_version
    # download_hash

    # if installed_hash_matches; then
    #     info 'Skipping binary downloaded, installed k3s matches hash'
    #     return
    # fi

    # download_binary
    # verify_binary
    # setup_binary
}

InstallApp() {
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

download_and_verify