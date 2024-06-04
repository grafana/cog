#!/usr/bin/env bash

# Exit on error. Append "|| true" if you expect an error.
set -o errexit
# Exit on error inside any functions or subshells.
set -o errtrace
# Do not allow use of undefined vars. Use ${VAR:-} to use an undefined VAR
set -o nounset
# Catch the error in case mysqldump fails (but gzip succeeds) in `mysqldump |gzip`
set -o pipefail

bold=$(tput bold)
normal=$(tput sgr0)

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
source "${__dir}/libs/logs.sh"

versions=("v10.1.x" "v10.2.x" "v10.3.x" "v10.4.x" "v11.0.x" "next")

for version in "${versions[@]}"; do
  info "🪧 Releasing ${bold}${version}${normal}"
  $__dir/release-version.sh "${version}"
done
