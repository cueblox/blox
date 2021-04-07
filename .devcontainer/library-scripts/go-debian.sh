#!/usr/bin/env bash
#-------------------------------------------------------------------------------------------------------------
# Copyright (c) Microsoft Corporation. All rights reserved.
# Licensed under the MIT License. See https://go.microsoft.com/fwlink/?linkid=2090316 for license information.
#-------------------------------------------------------------------------------------------------------------
#
# Docs: https://github.com/microsoft/vscode-dev-containers/blob/master/script-library/docs/go.md
# Maintainer: The VS Code and Codespaces Teams
#
# Syntax: ./go-debian.sh [Go version] [GOROOT] [GOPATH] [non-root user] [Add GOPATH, GOROOT to rc files flag] [Install tools flag]

TARGET_GO_VERSION=${1:-"latest"}
TARGET_GOROOT=${2:-"/usr/local/go"}
TARGET_GOPATH=${3:-"/go"}
USERNAME=${4:-"automatic"}
UPDATE_RC=${5:-"true"}
INSTALL_GO_TOOLS=${6:-"true"}

set -e

if [ "$(id -u)" -ne 0 ]; then
    echo -e 'Script must be run as root. Use sudo, su, or add "USER root" to your Dockerfile before running this script.'
    exit 1
fi

# Ensure that login shells get the correct path if the user updated the PATH using ENV.
rm -f /etc/profile.d/00-restore-env.sh
echo "export PATH=${PATH//$(sh -lc 'echo $PATH')/\$PATH}" > /etc/profile.d/00-restore-env.sh
chmod +x /etc/profile.d/00-restore-env.sh

# Determine the appropriate non-root user
if [ "${USERNAME}" = "auto" ] || [ "${USERNAME}" = "automatic" ]; then
    USERNAME=""
    POSSIBLE_USERS=("vscode" "node" "codespace" "$(awk -v val=1000 -F ":" '$3==val{print $1}' /etc/passwd)")
    for CURRENT_USER in ${POSSIBLE_USERS[@]}; do
        if id -u ${CURRENT_USER} > /dev/null 2>&1; then
            USERNAME=${CURRENT_USER}
            break
        fi
    done
    if [ "${USERNAME}" = "" ]; then
        USERNAME=root
    fi
elif [ "${USERNAME}" = "none" ] || ! id -u ${USERNAME} > /dev/null 2>&1; then
    USERNAME=root
fi

function updaterc() {
    if [ "${UPDATE_RC}" = "true" ]; then
        echo "Updating /etc/bash.bashrc and /etc/zsh/zshrc..."
        echo -e "$1" >> /etc/bash.bashrc
        if [ -f "/etc/zsh/zshrc" ]; then
            echo -e "$1" >> /etc/zsh/zshrc
        fi
    fi
}

export DEBIAN_FRONTEND=noninteractive

# Install curl, tar, git, other dependencies if missing
if ! dpkg -s curl ca-certificates tar git g++ gcc libc6-dev make pkg-config > /dev/null 2>&1; then
    if [ ! -d "/var/lib/apt/lists" ] || [ "$(ls /var/lib/apt/lists/ | wc -l)" = "0" ]; then
        apt-get update
    fi
    apt-get -y install --no-install-recommends curl ca-certificates tar git g++ gcc libc6-dev make pkg-config
fi

# Get latest version number if latest is specified
if [ "${TARGET_GO_VERSION}" = "latest" ] ||  [ "${TARGET_GO_VERSION}" = "current" ] ||  [ "${TARGET_GO_VERSION}" = "lts" ]; then
    TARGET_GO_VERSION=$(curl -sSL "https://golang.org/VERSION?m=text" | sed -n '/^go/s///p' )
fi

# Install Go
GO_INSTALL_SCRIPT="$(cat <<EOF
    set -e
    echo "Downloading Go ${TARGET_GO_VERSION}..."
    curl -sSL -o /tmp/go.tar.gz "https://golang.org/dl/go${TARGET_GO_VERSION}.linux-amd64.tar.gz"
    echo "Extracting Go ${TARGET_GO_VERSION}..."
    tar -xzf /tmp/go.tar.gz -C "${TARGET_GOROOT}" --strip-components=1
    rm -f /tmp/go.tar.gz
EOF
)"
if [ "${TARGET_GO_VERSION}" != "none" ] && ! type go > /dev/null 2>&1; then
    mkdir -p "${TARGET_GOROOT}" "${TARGET_GOPATH}" 
    chown -R ${USERNAME} "${TARGET_GOROOT}" "${TARGET_GOPATH}"
    su ${USERNAME} -c "${GO_INSTALL_SCRIPT}"
else
    echo "Go already installed. Skipping."
fi

# Install Go tools that are isImportant && !replacedByGopls based on
# https://github.com/golang/vscode-go/blob/0c6dce4a96978f61b022892c1376fe3a00c27677/src/goTools.ts#L188
# exception: golangci-lint is installed using their install script below.
GO_TOOLS="\
    golang.org/x/tools/gopls \
    honnef.co/go/tools/... \
    golang.org/x/lint/golint \
    github.com/mgechev/revive \
    github.com/uudashr/gopkgs/v2/cmd/gopkgs \
    github.com/ramya-rao-a/go-outline \
    github.com/go-delve/delve/cmd/dlv \
    github.com/golangci/golangci-lint/cmd/golangci-lint"
if [ "${INSTALL_GO_TOOLS}" = "true" ]; then
    echo "Installing common Go tools..."
    export PATH=${TARGET_GOROOT}/bin:${PATH}
    mkdir -p /tmp/gotools /usr/local/etc/vscode-dev-containers ${TARGET_GOPATH}/bin
    cd /tmp/gotools
    export GOPATH=/tmp/gotools
    export GOCACHE=/tmp/gotools/cache

    # Go tools w/module support
    export GO111MODULE=on
    (echo "${GO_TOOLS}" | xargs -n 1 go get -v )2>&1 | tee -a /usr/local/etc/vscode-dev-containers/go.log

    # Move Go tools into path and clean up
    mv /tmp/gotools/bin/* ${TARGET_GOPATH}/bin/
    rm -rf /tmp/gotools
    chown -R ${USERNAME} "${TARGET_GOPATH}"
fi

# Add GOPATH variable and bin directory into PATH in bashrc/zshrc files (unless disabled)
updaterc "$(cat << EOF
export GOPATH="${TARGET_GOPATH}"
if [[ "\${PATH}" != *"\${GOPATH}/bin"* ]]; then export PATH="\${PATH}:\${GOPATH}/bin"; fi
export GOROOT="${TARGET_GOROOT}"
if [[ "\${PATH}" != *"\${GOROOT}/bin"* ]]; then export PATH="\${PATH}:\${GOROOT}/bin"; fi
EOF
)"

echo "Done!"

