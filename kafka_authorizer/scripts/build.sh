#!/usr/bin/env bash

DIR="$( cd "$(dirname "$0")" ; pwd -P )"

docker run -w /src -v $DIR/..:/src maven:3.5.3-jdk-8 mvn install
