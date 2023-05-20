#!/bin/sh

set -e

kong_host="$1"
opa_host="$2"
httpbin_host="$3"
opaproxy_host="$4"
shift 4
cmd="$@"

until $(curl --output /dev/null --silent --head --fail http://$kong_host/status); do
  >&2 echo "Kong is unavailable - sleeping"
  sleep 5
done

until $(curl --output /dev/null --silent --fail "http://$opa_host/health?bundle=true"); do
  >&2 echo "Opa is unavailable - sleeping"
  sleep 3
done

until $(curl --output /dev/null --silent --head --fail http://$httpbin_host/status/200); do
  >&2 echo "Httpbin is unavailable - sleeping"
  sleep 1
done

until $(curl --output /dev/null --silent --fail http://$opaproxy_host/proxies); do
  >&2 echo "Toxiproxy is unavailable - sleeping"
  sleep 1
done

>&2 echo "All services are up - executing command"
# use newman as the default container command
exec newman $cmd
