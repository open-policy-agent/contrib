#!/usr/bin/env bash

TARGET=$1

if [ -z "${TARGET}" ]; then
    echo "usage: $0 {target}"
    exit 1
fi

CONTRIB_ROOT=$(dirname "${BASH_SOURCE}")/..

for file in ${CONTRIB_ROOT}/*/Makefile; do

    # Assume that contribution Makefiles are NOT written to be run from any location.
    make -C $(dirname $file) $TARGET

done
