#!/usr/bin/env bash

set -euo pipefail

__dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

DESTINATION_DIR="${1}"
REPO_URL="${2}"

function usage() {
  echo "Usage: ${BASH_SOURCE[0]} DESTINATION_DIR REPO_URL"
}

if [ "${DESTINATION_DIR}" == "" ]; then
  usage
  exit 1
fi

if [ "${REPO_URL}" == "" ]; then
  usage
  exit 1
fi

if [ -d $DESTINATION_DIR/.git ]; then
  echo "Updating $REPO_URL..."
  git -C "$DESTINATION_DIR" fetch --all --quiet
  git -C "$DESTINATION_DIR" reset --hard origin/main --quiet
else
  echo "Cloning $REPO_URL..."
  git clone --depth=1 "$REPO_URL" "$DESTINATION_DIR"
fi

