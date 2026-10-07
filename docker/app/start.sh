#!/bin/bash
set -euo pipefail

lock=/go/pkg/mod/.download.lock
mkdir -p /go/pkg/mod
while ! mkdir "$lock" 2>/dev/null; do
  sleep 0.2
done

cd "$SERVICE"
go mod download

rmdir "$lock"
exec go run .