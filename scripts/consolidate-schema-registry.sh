#!/usr/bin/env bash

# Exit on error. Append "|| true" if you expect an error.
set -o errexit
# Exit on error inside any functions or subshells.
set -o errtrace
# Do not allow use of undefined vars. Use ${VAR:-} to use an undefined VAR
set -o nounset
# Catch the error in case mysqldump fails (but gzip succeeds) in `mysqldump |gzip`
set -o pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${__dir}/libs/logs.sh"

GRAFANA_VERSION=${GRAFANA_VERSION:-"next"} # version of the schemas/grafana to consolidate.
KIND_REGISTRY_PATH=${KIND_REGISTRY_PATH:-'../kind-registry'} # Path to the kind-registry
COG_EMBEDDED_SCHEMAS_PATH=${COG_EMBEDDED_SCHEMAS_PATH:-'./schemas'} # Path to the kind-registry

#################
### Utilities ###
#################

function copy_prometheus_schema() {
  local kind_registry_path=${1}
  shift
  local target_version=${1}
  shift
  local embedded_schemas_path=${1}
  shift

  rm -rf "${kind_registry_path}/grafana/${target_version}/composable/prometheus"
  cp -R "${embedded_schemas_path}/composable/prometheus" "${kind_registry_path}/grafana/${target_version}/composable"
}

###########################
### Consolidation cases ###
###########################

if [ "${GRAFANA_VERSION}" == "next" ]; then
  debug "Consolidating schemas for ${GRAFANA_VERSION}"

  debug " → copy_prometheus_schema"
  copy_prometheus_schema "${KIND_REGISTRY_PATH}" "${GRAFANA_VERSION}" "${COG_EMBEDDED_SCHEMAS_PATH}"
fi
