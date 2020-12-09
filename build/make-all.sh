#!/usr/bin/env bash

set -e

TARGET=$1

if [ -z "${TARGET}" ]; then
    echo "usage: $0 {target}"
    exit 1
fi

CONTRIB_ROOT=$(dirname "${BASH_SOURCE}")/..

for file in ${CONTRIB_ROOT}/*/Makefile; do

    dir=$(dirname "$file")
    base=$(basename "$dir")

    # TODO: Fix these two currently failing GH actions builds
    if [ "$base" == "data_filter_mongodb" ] || [ "$base" == "data_filter_example" ]; then
        continue
    fi

    # Assume that contribution Makefiles are NOT written to be run from any location.
    make -C "$dir" "$TARGET"

done
