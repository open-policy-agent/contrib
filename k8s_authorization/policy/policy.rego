package k8s.authz

deny[reason] {
	input.spec.resourceAttributes.namespace == "kube-system"

	reason := "OPA: denied access to namespace kube-system"
}

deny[reason] {
	input.spec.resourceAttributes.namespace == "opa"

	required_groups := {"system:authenticated", "devops"}
	provided_groups := {group | group := input.spec[groups][_]}

	count(required_groups & provided_groups) != count(required_groups)

	reason := sprintf("OPA: provided groups (%v) does not include all required groups: (%v)", [
		concat(", ", provided_groups),
		concat(", ", required_groups),
	])
}

decision = {
	"apiVersion": input.apiVersion,
	"kind": "SubjectAccessReview",
	"status": {
		"allowed": count(deny) == 0,
		"reason": concat(" | ", deny),
	},
}

# Take into account the typo of missing "s" in groups
# in v1beta1 version (which is used by default)
# https://github.com/kubernetes/kubernetes/issues/32709
groups = {
	"authorization.k8s.io/v1": "groups",
	"authorization.k8s.io/v1beta1": "group",
}[input.apiVersion]
