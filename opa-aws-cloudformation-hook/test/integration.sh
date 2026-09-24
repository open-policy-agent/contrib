#!/usr/bin/env bash

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

opa run --server "$SCRIPT_DIR/../examples/policy/" &>/dev/null &

pid=$(echo $!)

# Provide some time for OPA to start
sleep 1

"$SCRIPT_DIR"/validate.py --test=true "$SCRIPT_DIR"/../examples/templates

status=$(echo $?)

kill "$pid"

exit "$status"
