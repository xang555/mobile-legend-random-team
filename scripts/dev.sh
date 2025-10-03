#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR=$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)
CONFIG_PATH=${1:-${ROOT_DIR}/configs/config.yaml}

cd "${ROOT_DIR}"

go run ./cmd/server -config "${CONFIG_PATH}"
