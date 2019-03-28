#!/usr/bin/env bash
set -euo pipefail

if [[ -f "/tmp/retried" ]]; then
    rm "/tmp/retried"
    exit 0
fi

touch /tmp/retried
exit 1
