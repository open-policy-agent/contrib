#!/usr/bin/env bash

set -e

DIR="$( cd "$(dirname "$0")" ; pwd -P )"

docker run --rm -v $DIR/..:/src mribeiro/xmllint --xpath "//*[local-name()='project']/*[local-name()='version']/text()" /src/pom.xml
