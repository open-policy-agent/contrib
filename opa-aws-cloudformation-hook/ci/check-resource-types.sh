#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

current=$(cat "$SCRIPT_DIR/../hooks/opa-hook.json")

new=$(echo "$current" \
| opa eval --stdin-input \
           --format pretty \
           --data "$SCRIPT_DIR/resourcetypes.rego" \
           data.aws.cloudformation.output)

if [[ "$current" != "$new" ]]; then
    echo "Resource types have been updated. Please run:"
    echo
    echo "cat hooks/opa-hook.json | opa eval -I -f pretty -d ci/resourcetypes.rego data.aws.cloudformation.output > hooks/opa-hook-new.json"
    echo
    echo "mv hooks/opa-hook-new.json hooks/opa-hook.json"
    echo
    echo "And commit the result"
    exit 1
fi
