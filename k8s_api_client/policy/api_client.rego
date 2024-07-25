package kubernetes.api.client

import rego.v1

# This information could be retrieved from the kubernetes API
# too, but would essentially require a request per API group,
# so for now use a lookup table for the most common resources.
resource_group_mapping := {
	"services": "api/v1",
	"pods": "api/v1",
	"configmaps": "api/v1",
	"secrets": "api/v1",
	"persistentvolumeclaims": "api/v1",
	"daemonsets": "apis/apps/v1",
	"deployments": "apis/apps/v1",
	"statefulsets": "apis/apps/v1",
	"horizontalpodautoscalers": "api/autoscaling/v1",
	"jobs": "apis/batch/v1",
	"cronjobs": "apis/batch/v1beta1",
	"ingresses": "api/extensions/v1beta1",
	"replicasets": "apis/apps/v1",
	"networkpolicies": "apis/networking.k8s.io/v1",
}

# Query for given resource/name in provided namespace
# Example: query_ns("deployments", "my-app", "default")
query_name_ns(resource, name, namespace) := http.send({
	"url": sprintf("https://kubernetes.default.svc/%v/namespaces/%v/%v/%v", [
		# regal ignore:external-reference
		resource_group_mapping[resource],
		namespace,
		resource,
		name,
	]),
	"method": "get",
	"headers": {"authorization": sprintf("Bearer %v", [opa.runtime().env.KUBERNETES_API_TOKEN])},
	"tls_ca_cert_file": "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
	"raise_error": false,
})

# Query for given resource type using label selectors
# https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#api
# Example: query_label_selector_ns("deployments", {"app": "opa-kubernetes-api-client"}, "default")
query_label_selector_ns(resource, selector, namespace) := http.send({
	"url": sprintf("https://kubernetes.default.svc/%v/namespaces/%v/%v?labelSelector=%v", [
		# regal ignore:external-reference
		resource_group_mapping[resource],
		namespace,
		resource,
		label_map_to_query_string(selector),
	]),
	"method": "get",
	"headers": {"authorization": sprintf("Bearer %v", [opa.runtime().env.KUBERNETES_API_TOKEN])},
	"tls_ca_cert_file": "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
	"raise_error": false,
})

# Query for all resources of type resource in all namespaces
# Example: query_all("deployments")
query_all(resource) := http.send({
	"url": sprintf("https://kubernetes.default.svc/%v/%v", [
		# regal ignore:external-reference
		resource_group_mapping[resource],
		resource,
	]),
	"method": "get",
	"headers": {"authorization": sprintf("Bearer %v", [opa.runtime().env.KUBERNETES_API_TOKEN])},
	"tls_ca_cert_file": "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt",
	"raise_error": false,
})

label_map_to_query_string(map) := concat(",", [str | val := map[key]; str := concat("%3D", [key, val])])
