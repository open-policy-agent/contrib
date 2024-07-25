package kafka.authz

import rego.v1

default allow := false

allow if {
	not deny
}

deny if {
	is_read_operation
	topic_contains_pii
	not consumer_is_whitelisted_for_pii
}

deny if {
	is_write_operation
	topic_has_large_fanout
	not producer_is_whitelisted_for_large_fanout
}

###############################################################################
# Example whitelists. For conciseness in the tutorial, the whitelists are
# hardcoded inside the policy. In real-world deployments, the whitelists could
# be loaded into OPA as raw JSON data.
###############################################################################

producer_whitelist := {"large-fanout": {"fanout_producer"}}

consumer_whitelist := {"pii": {"pii_consumer"}}

topic_metadata := {
	"click-stream": {"tags": ["large-fanout"]},
	"credit-scores": {"tags": ["pii"]},
}

###############################################################################
# Helper rules for checking whitelists.
###############################################################################

topic_contains_pii if {
	"pii" in topic_metadata[topic_name].tags
}

topic_has_large_fanout if {
	"large-fanout" in topic_metadata[topic_name].tags
}

consumer_is_whitelisted_for_pii if {
	principal.name in consumer_whitelist.pii
}

producer_is_whitelisted_for_large_fanout if {
	principal.name in producer_whitelist["large-fanout"]
}

###############################################################################
# Helper rules for input processing.
###############################################################################

is_write_operation if {
	input.operation.name == "Write"
}

is_read_operation if {
	input.operation.name == "Read"
}

is_topic_resource if {
	input.resource.resourceType.name == "Topic"
}

topic_name := input.resource.name if {
	is_topic_resource
}

principal := {"fqn": parsed.CN, "name": cn_parts[0]} if {
	parsed := parse_user(urlquery.decode(input.session.sanitizedUser))
	cn_parts := split(parsed.CN, ".")
}

parse_user(user) := {key: value |
	parts := split(user, ",")
	[key, value] := split(parts[_], "=")
}
