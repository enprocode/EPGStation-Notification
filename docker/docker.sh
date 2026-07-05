#!/bin/bash

set -euo pipefail

docker exec EPGStation-Notification bash -c \
  'cd /opt/src && GOOS=linux GOARCH=amd64 go build -o bin/epgst-notify main.go'
