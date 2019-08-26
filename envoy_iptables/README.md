# Proxy Init

This directory contains the Istio proxy init script and a Dockerfile for
building an image for the init container that installs iptables rules to
redirect all container traffic through the Envoy proxy sidecar.
