#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

# shellcheck disable=SC2155
readonly ROOT_DIR=$(realpath "$(dirname "$0")/..")
readonly KARMADA_CONFIG_DIR="$HOME/.karmadaoperator"
readonly MANAGEMENT_CLUSTER_NAME="mc"
readonly PORTS_FILE="$KARMADA_CONFIG_DIR/available_ports.txt"
readonly KARMADARC_FILE="$KARMADA_CONFIG_DIR/karmadarc"
readonly CLUSTERS_FILE="$KARMADA_CONFIG_DIR/clusters.txt"
# shellcheck disable=SC2155
readonly ARCH=$(uname -m)
TRACE_FLAG_SET=false

while [[ $# -gt 0 ]]; do
    case "$1" in
        -t|--trace)
            TRACE_FLAG_SET=true
            shift
            ;;
        *)
            break
            ;;
    esac
done

if [[ "${KARMADA_TRACE:-}" || "$TRACE_FLAG_SET" == true ]]; then
    set -o xtrace
fi

export VERSION="v1.9.0"
export REGISTRY="docker.io/karmada"

ensure_prerequisites_are_met() {
    echo "Ensuring prerequisites are met..."
    if [[ $(uname -s) != "Darwin" ]]; then
        echo "Error: This script is only supported on macOS."
        exit 1
    fi

    if [[ "$ARCH" != "x86_64" && "$ARCH" != "arm64" ]]; then
        echo "Error: Unsupported architecture. This script is only supported on amd64 and arm64 architectures."
        exit 1
    fi

    for tool in docker kind kubectl helm; do
        if ! command -v "$tool" &> /dev/null; then
            echo "Error: $tool is not installed. Please ensure it is installed and try again."
            exit 1
        fi
    done

    echo "Prerequisites met."
}

init() {
    echo "Initializing the infrastructure..."
    ensure_prerequisites_are_met

    if [[ -d "$KARMADA_CONFIG_DIR" ]]; then
        echo "Infrastructure is already initialized. If you'd like to start over, first run the destroy command and try again."
        exit 1
    fi

    mkdir -p "$KARMADA_CONFIG_DIR"
    touch "$PORTS_FILE"

    for port in $(seq 30000 32767); do
        if ! lsof -i:"$port" >/dev/null; then
            echo "$port" >> "$PORTS_FILE"
        fi
        [[ $(wc -l < "$PORTS_FILE") -ge 20 ]] && break
    done

    local nic_ip
    nic_ip=$(ipconfig getifaddr en0)
    workspace=$(mktemp -d)
    trap 'rm -rf "$workspace"' EXIT

    # shellcheck disable=SC2155
    local management_cluster_config=$(cat <<EOF
kind: Cluster
apiVersion: "kind.x-k8s.io/v1alpha4"
networking:
  apiServerAddress: ${nic_ip}
nodes:
  - role: control-plane
    image: kindest/node:v1.26.4
    extraPortMappings:
EOF
)

    while read -r port; do
        management_cluster_config+="\n      - containerPort: $port\n"
        management_cluster_config+="        hostPort: $port\n"
        management_cluster_config+="        protocol: TCP\n"
        management_cluster_config+="        listenAddress: $nic_ip\n"
    done < "$PORTS_FILE"

    echo -e "$management_cluster_config" > "$workspace/management.yaml"

    local kubeconfig_path="$KARMADA_CONFIG_DIR/$MANAGEMENT_CLUSTER_NAME/kubeconfig"
    mkdir -p "$(dirname "$kubeconfig_path")"
    echo "Creating management cluster..."
    kind create cluster --name "$MANAGEMENT_CLUSTER_NAME" --config "$workspace/management.yaml" --kubeconfig "$kubeconfig_path"

    echo "alias $MANAGEMENT_CLUSTER_NAME='kubectl --kubeconfig=$kubeconfig_path'" >> "$KARMADARC_FILE"

    git tag -d "${VERSION}" || git tag "${VERSION}"
    IMAGES=(
        "karmada-operator"
        "karmada-controller-manager"
        "karmada-scheduler"
        "karmada-webhook"
        "karmada-aggregated-apiserver"
        "karmada-metrics-adapter"
    )

    for image in "${IMAGES[@]}"; do
        echo "Building ${image} image..."
        make image-"${image}" GOOS="linux" --directory="${ROOT_DIR}"
    done

    for image in "${IMAGES[@]}"; do
        echo "loading ${image} image into the management cluster..."
        kind load docker-image "${REGISTRY}/${image}:${VERSION}" --name="${MANAGEMENT_CLUSTER_NAME}"
    done


    echo "Installing Karmada operator..."
    #helm dependency update "$ROOT_DIR/charts/karmada-operator"
    helm install karmada-operator  "$ROOT_DIR/charts/karmada-operator" \
        --create-namespace \
        --namespace karmada-system \
        --set operator.image.tag="$VERSION" \
        --kubeconfig "$kubeconfig_path"

    touch "$CLUSTERS_FILE"
    echo "$MANAGEMENT_CLUSTER_NAME" > "$CLUSTERS_FILE"

    echo "Infrastructure has been successfully initialized."
    echo "Source the rc file to activate aliases: source $KARMADARC_FILE"
}

add_tenant() {
    if [[ -z "$1" ]]; then
        echo "Error: Tenant name not provided."
        echo "Usage: $0 add-tenant <tenant-name>"
        exit 1
    fi

    local tenant_name="$1"
    echo "Adding tenant '$tenant_name'..."

    if [[ $(wc -l < "$PORTS_FILE") -lt 1 ]]; then
        echo "Error: No available ports left to bind tenant's API server."
        exit 1
    fi

    local kubeconfig_path="$KARMADA_CONFIG_DIR/$MANAGEMENT_CLUSTER_NAME/kubeconfig"
    if kubectl --kubeconfig="$kubeconfig_path" get namespace "$tenant_name" >/dev/null 2>&1; then
        echo "Error: Tenant '$tenant_name' already exists."
        exit 1
    fi

    workspace=$(mktemp -d)
    trap 'rm -rf "$workspace"' EXIT

    local nic_ip
    nic_ip=$(ipconfig getifaddr en0)
    # shellcheck disable=SC2155
    local tenant_manifest=$(cat <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: ${tenant_name}
---
apiVersion: operator.karmada.io/v1alpha1
kind: Karmada
metadata:
  name: ${tenant_name}
  namespace: ${tenant_name}
spec:
  components:
    karmadaAPIServer:
      certSANs:
        - ${nic_ip}
      serviceType: NodePort
    karmadaAggregatedAPIServer:
      imageRepository: ${REGISTRY}/karmada/karmada-aggregated-apiserver
      imageTag: ${VERSION}
    karmadaControllerManager:
      imageRepository: ${REGISTRY}/karmada/karmada-controller-manager
      imageTag: ${VERSION}
    karmadaMetricsAdapter:
      imageRepository: ${REGISTRY}/karmada/karmada-metrics-adapter
      imageTag: ${VERSION}
    karmadaScheduler:
      imageRepository: ${REGISTRY}/karmada/karmada-scheduler
      imageTag: ${VERSION}
    karmadaWebhook:
      imageRepository: ${REGISTRY}/karmada/karmada-webhook
      imageTag: ${VERSION}
EOF
)

    echo -e "$tenant_manifest" > "$workspace/tenant.yaml"

    echo "Creating tenant..."
    kubectl --kubeconfig="$kubeconfig_path" apply -f "$workspace/tenant.yaml"

    echo "Waiting for tenant's control plane to reach ready state..."
    kubectl wait --for=condition=Ready --kubeconfig="$kubeconfig_path" -n "$tenant_name" karmada/"$tenant_name" --timeout 5m

    echo "Binding tenant's API server node port..."
    local available_port
    available_port=$(head -n 1 "$PORTS_FILE")
    kubectl --kubeconfig="$kubeconfig_path" -n "$tenant_name" patch service "$tenant_name-karmada-apiserver" -p "{\"spec\":{\"ports\":[{\"port\":5443,\"nodePort\":$available_port,\"protocol\":\"TCP\"}]}}"

    sed -i "" "1d" "$PORTS_FILE"

    local kubeconfig_secret="$tenant_name-karmada-admin-config"
    kubectl --kubeconfig="$kubeconfig_path" -n "$tenant_name" get secret "$kubeconfig_secret" -o jsonpath='{.data.kubeconfig}' | base64 --decode > "$workspace/kubeconfig"
    sed -i "" "s|server:.*|server: https://$nic_ip:$available_port|" "$workspace/kubeconfig"
    local tenant_kubeconfig_dir="$KARMADA_CONFIG_DIR/tenants/$tenant_name"
    mkdir -p "$tenant_kubeconfig_dir"
    mv "$workspace/kubeconfig" "$tenant_kubeconfig_dir/kubeconfig"

    echo "alias ${tenant_name}='kubectl --kubeconfig=$KARMADA_CONFIG_DIR/tenants/${tenant_name}/kubeconfig'" >> "$KARMADARC_FILE"

    echo "Successfully added '$tenant_name' tenant."
    echo "Source the rc file to activate aliases: source $KARMADARC_FILE"
}

rm_tenant() {
    if [[ -z "$1" ]]; then
        echo "Error: Tenant name not provided."
        echo "Usage: $0 rm-tenant <tenant-name>"
        exit 1
    fi

    if [[ ! -f "$CLUSTERS_FILE" ]]; then
        echo "Error: Environment not initialized. Please run the init command first."
        exit 1
    fi

    local tenant_name="$1"
    echo "Removing tenant '$tenant_name'..."

    local kubeconfig_path="$KARMADA_CONFIG_DIR/$MANAGEMENT_CLUSTER_NAME/kubeconfig"
    if ! kubectl --kubeconfig="$kubeconfig_path" get namespace "$tenant_name" >/dev/null 2>&1; then
        echo "Error: Tenant '$tenant_name' does not exist."
        exit 1
    fi

    local tenant_kubeconfig_dir="$KARMADA_CONFIG_DIR/tenants/$tenant_name"
    echo "Deleting tenant's control plane instance..."
    kubectl --kubeconfig="$kubeconfig_path" delete namespace "$tenant_name"

    echo "Removing tenant's config directory..."
    rm -rf "$tenant_kubeconfig_dir"

    echo "Removing tenant aliases from rc file..."
    sed -i "" "/alias ${tenant_name}/d" "$KARMADARC_FILE"

    echo "Successfully removed tenant '$tenant_name'."
}

destroy() {
    echo "Destroying the infrastructure's environment..."
    if [[ -f "$CLUSTERS_FILE" ]]; then
        while read -r cluster; do
            echo "Deleting kind cluster '$cluster'..."
            kind delete cluster --name "$cluster"
            echo "Deleted kind cluster '$cluster'"
        done < "$CLUSTERS_FILE"
    else
        echo "No clusters.txt file found. Nothing to delete."
    fi
    rm -rf "$KARMADA_CONFIG_DIR"
    echo "Successfully destroyed infrastructure's environment."
}

ls_tenants() {
    if [[ ! -f "$CLUSTERS_FILE" ]]; then
        echo "Error: Environment not initialized. Please run the init command first."
        exit 1
    fi

    local kubeconfig_path="$KARMADA_CONFIG_DIR/$MANAGEMENT_CLUSTER_NAME/kubeconfig"
    local tenants
    tenants=$(kubectl --kubeconfig="$kubeconfig_path" get karmadas --all-namespaces 2>/dev/null || true)

    if [[ -z "$tenants" ]]; then
        echo "No Karmada tenancies exist. To add a tenant, run: $0 add-tenant <tenant-name>"
    else
        echo "Tenancies==========================="
        echo "$tenants"
    fi
}

help() {
    echo "Usage: $0 [global options] {init|destroy|add-tenant <tenant-name>|rm-tenant <tenant-name>|ls-tenants|help}"
    echo
    echo "Commands:"
    echo "  init             Initialize the infrastructure"
    echo "  destroy          Destroy the infrastructure's environment"
    echo "  add-tenant       Add a tenant"
    echo "  rm-tenant        Remove a tenant"
    echo "  ls-tenants       List all tenants"
    echo "  help             Display this help message"
    echo
    echo "Global Options:"
    echo "  -t, --trace      Enable tracing for debugging purposes"
}

if [[ $# -lt 1 ]]; then
    help
    exit 1
fi

command="$1"
shift

case "$command" in
    init)
        init
        ;;
    destroy)
        destroy
        ;;
    add-tenant)
        add_tenant "$@"
        ;;
    rm-tenant)
        rm_tenant "$@"
        ;;
    ls-tenants)
        ls_tenants
        ;;
    help)
        help
        ;;
    *)
        echo "Error: Unknown command: $command"
        help
        exit 1
        ;;
esac
