#!/usr/bin/env bash

set -exo pipefail

try_download_opa() {
    curl -vL https://api.github.com/repos/open-policy-agent/opa/releases/latest |
        grep 'browser_download_url.*linux_amd64' |
        cut -d : -f 2,3 |
        tr -d '"' |
        wget -i -
}

download_opa() {
    local attempt=0
    local max_attempts=60
    local delay=1
    local max_delay=30
    while [ 1 ]; do
        attempt=$((attempt + 1))
        rc=0
        try_download_opa || rc=$?
        if [[ $rc -eq 0 ]]; then
            break
        fi
        if (( attempt > max_attempts )); then
            echo "Failed to download OPA binary from GitHub."
            exit 1
        fi
        sleep $delay
        delay=$((delay * 2))
        if (( delay > max_delay )); then
            delay=$max_delay
        fi
    done
}

download_opa
chmod u+x opa_linux_amd64
ln -s opa_linux_amd64 opa
