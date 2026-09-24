package aws.eks.cluster

import rego.v1

deny contains msg if {
	not cluster_logging_enabled
	msg := sprintf("no logging types are enabled for cluster: %s", [input.resource.id])
}

cluster_logging_enabled if count(input.resource.properties.Logging.ClusterLogging.EnabledTypes) > 0
