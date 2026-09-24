#/usr/bin/env bash

# Deregister the 'OpenPolicyAgent::OPA::Hook' type hook and all of its versions

if [[ -z "$1" ]]; then
    echo "Usage:"
    echo
    echo "./deregister-hook.sh <arn>"
    echo
    echo "The ARN of a hook can be found with the command 'aws cloudformation list-types'"
fi

arn="$1"

versions=$(aws cloudformation list-type-versions --arn "$arn" | jq -r '.TypeVersionSummaries[].VersionId')

for version in $versions; do
    echo "Deregistering version: $version"
    aws cloudformation deregister-type \
        --type HOOK \
        --type-name "OpenPolicyAgent::OPA::Hook" \
        --version-id "$version"
done

echo "Deregistering hook OpenPolicyAgent::OPA::Hook"

aws cloudformation deregister-type --type HOOK --type-name "OpenPolicyAgent::OPA::Hook"
