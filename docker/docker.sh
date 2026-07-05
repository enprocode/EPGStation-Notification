#!/bin/bash

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

docker compose -f "${SCRIPT_DIR}/docker-compose.yml" exec -T go \
  bash -c 'GOOS=linux GOARCH=amd64 go build -o bin/epgst-notify ./main.go'
