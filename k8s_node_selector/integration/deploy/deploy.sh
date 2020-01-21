#! /bin/bash
cd "$(dirname "$0")"

set -e

kubectl create namespace opa

kubectl config set-context --current --namespace opa

openssl genrsa -out ca.key 2048
openssl req -x509 -new -nodes -key ca.key -days 100000 -out ca.crt -subj "/CN=admission_ca"

cat >server.conf <<EOF
[req]
req_extensions = v3_req
distinguished_name = req_distinguished_name
[req_distinguished_name]
[ v3_req ]
basicConstraints = CA:FALSE
keyUsage = nonRepudiation, digitalSignature, keyEncipherment
extendedKeyUsage = clientAuth, serverAuth
EOF

openssl genrsa -out server.key 2048
openssl req -new -key server.key -out server.csr -subj "/CN=opa.opa.svc" -config server.conf
openssl x509 -req -in server.csr -CA ca.crt -CAkey ca.key -CAcreateserial -out server.crt -days 100000 -extensions v3_req -extfile server.conf

kubectl create secret tls opa-server --cert=server.crt --key=server.key

kubectl apply -f admission-controller.yaml

kubectl label --overwrite ns kube-system openpolicyagent.org/webhook=ignore
kubectl label --overwrite ns opa openpolicyagent.org/webhook=ignore

# Setup Authentication. This ensures that only the kubemgmt sidecar can call the OPA server
# Without this setup others in the cluster could put new policies without authentication
OPA_AUTH_TOKEN=$(LC_CTYPE=C cat /dev/urandom | tr -dc 'a-z0-9' | fold -w 12 | head -n 1)
cp ../../authz.rego . && sed -i s/{TOKEN_HERE}/$OPA_AUTH_TOKEN/ authz.rego
kubectl create secret generic authz-policy -n opa --from-file=./authz.rego --from-literal=token="$OPA_AUTH_TOKEN"


cat > webhook-configuration.yaml <<EOF
kind: MutatingWebhookConfiguration
apiVersion: admissionregistration.k8s.io/v1beta1
metadata:
  name: opa-validating-webhook
webhooks:
  - name: validating-webhook.openpolicyagent.org
    namespaceSelector:
      matchExpressions:
      - key: openpolicyagent.org/webhook
        operator: NotIn
        values:
        - ignore
    rules:
      - operations: ["CREATE", "UPDATE"]
        apiGroups: ["*"]
        apiVersions: ["*"]
        resources: ["*"]
    clientConfig:
      caBundle: $(cat ca.crt | base64 | tr -d '\n')
      service:
        namespace: opa
        name: opa
EOF

kubectl apply -f webhook-configuration.yaml

# Deploy the rego policy
kubectl create configmap opa-default-system-main --from-file ./../../main.rego




