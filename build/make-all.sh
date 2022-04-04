#!/usr/bin/env bash

set -e

TARGET=$1

if [ -z "${TARGET}" ]; then
    echo "usage: $0 {target}"
    exit 1
fi

CONTRIB_ROOT=$(dirname "${BASH_SOURCE}")/..

# This works around the error:
#
#	Error while loading /usr/local/sbin/dpkg-split: No such file or directory
#
# Which can occur when using `docker buildx` to perform multi-architecture
# builds in GitHub Actions.
#
# See: https://github.com/docker/buildx/issues/495#issuecomment-761562905
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes
docker buildx create --name multiarch --driver docker-container --use
docker buildx inspect --bootstrap

for file in ${CONTRIB_ROOT}/*/Makefile; do

    dir=$(dirname "$file")
    base=$(basename "$dir")

    # TODO: Fix these two currently failing GH actions builds
    if [ "$base" == "gatekeeper_mtail_violations_exporter" ] || [ "$base" == "data_filter_mongodb" ] || [ "$base" == "data_filter_example" ]; then
        continue
    fi

    # Assume that contribution Makefiles are NOT written to be run from any location.
    make -C "$dir" "$TARGET"

done
