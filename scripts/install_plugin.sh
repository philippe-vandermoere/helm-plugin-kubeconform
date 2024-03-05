#!/bin/sh

set -e

if [ -n "${HELM_KUBCONFORM_PLUGIN_NO_INSTALL_HOOK}" ]; then
    echo "Development mode: not downloading versioned release."
    exit 0
fi

if [ "${HELM_DEBUG:-}" = true ]; then
    set -x
fi

version="$(grep -oP 'version: *\K([0-9]+\.[0-9]+\.[0-9]+)' plugin.yaml)"
echo "Downloading and installing ${HELM_PLUGIN_NAME:-kubeconform} v${version}"

case $(uname -m) in
    x86_64)
        arch="amd64"
        ;;
    armv6*)
        arch="armv6"
        ;;
    armv7*)
        arch="armv7"
        ;;
    aarch64 | arm64)
        arch="arm64"
        ;;
    *)
        echo "Failed to detect target architecture"
        exit 1
        ;;
esac

if [ "$(uname)" = "Linux" ]; then
    os=linux
elif [ "$(uname)" = "Darwin" ] ; then
    os=darwin
else
    os=windows
fi

download_url="https://github.com/philippe-vandermoere/helm-plugin-kubeconform/releases/download/v${version}/helm-kubeconform_v${version}_${os}_${arch}.tar.gz"
TMP_DIR="$(mktemp -d)"
if [ -x "$(which curl 2>/dev/null)" ]; then
    curl -fsSL "${download_url}" -o "${TMP_DIR}/v${version}.tar.gz"
else
    wget -q "${download_url}" -O "${TMP_DIR}/v${version}.tar.gz"
fi

tar xzvf "${TMP_DIR}/v${version}.tar.gz" -C "${TMP_DIR}"
mv "${TMP_DIR}/bin/"* "${HELM_PLUGIN_DIR}/bin/"
rm -rf "${TMP_DIR}"
