#!/usr/bin/env bash

set -exo pipefail

curl -vL https://api.github.com/repos/open-policy-agent/opa/releases/latest |
    grep 'browser_download_url.*linux_amd64' |
    cut -d : -f 2,3 |
    tr -d '"' |
    wget -qi -

chmod u+x opa_linux_amd64
ln -s opa_linux_amd64 opa
