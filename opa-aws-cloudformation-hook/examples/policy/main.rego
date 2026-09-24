# METADATA
# description: |
#   Dynamic routing to policy based in input.resource.type,
#   aggregating the deny rules found in all policies with a
#   matching package name
#
package system

import rego.v1

# METADATA
# entrypoint: true
main := {
	"allow": count(violations) == 0,
	"violations": violations,
}

# METADATA
# description: |
#   Main routing logic, simply converting input.resource.type, e.g.
#   AWS::S3::Bucket to data.aws.s3.bucket and returning that document.
#
#   By default, only input.action == "CREATE" | "UPDATE" will be routed
#   to the data.aws.s3.bucket document. If handling "DELETE" actions is
#   desirable, one may create a special policy for that by simply appending
#   "delete" to the package name, e.g. data.aws.s3.bucket.delete
#
route := document(lower(component), lower(type)) if ["AWS", component, type] = split(input.resource.type, "::")

violations contains msg if {
	# Aggregate all deny rules found in routed document
	some msg in route.deny
}

#
# Basic input validation to avoid having to do this in each resource policy
#

violations contains "Missing input.resource" if not input.resource

violations contains "Missing input.resource.type" if not input.resource.type

violations contains "Missing input.resource.id" if not input.resource.id

violations contains "Missing input.action" if not input.action

#
# Helpers
#

document(component, type) := data.aws[component][type] if input.action != "DELETE"

document(component, type) := data.aws[component][type].delete if input.action == "DELETE"
