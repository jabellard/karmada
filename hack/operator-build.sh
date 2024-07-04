#!/usr/bin/env bash
# Script to build the Karmada operator image

set -o errexit
set -o nounset
set -o pipefail
set -o xtrace

# shellcheck disable=SC2155
export REPO_ROOT=$(realpath "$(dirname "$0")/..")
export VERSION="latest"
export REGISTRY="docker.io/karmada"
make image-karmada-operator GOOS="linux" --directory="${REPO_ROOT}"
