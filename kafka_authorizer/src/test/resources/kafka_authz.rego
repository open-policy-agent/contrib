package kafka.authz

default allow = false

allow {
	not deny
}

deny {
	is_read_operation
	topic_contains_pii
	not consumer_is_whitelisted_for_pii
}

deny {
	is_write_operation
	topic_has_large_fanout
	not producer_is_whitelisted_for_large_fanout
}

###############################################################################
# Example whitelists. For conciseness in the tutorial, the whitelists are
# hardcoded inside the policy. In real-world deployments, the whitelists could
# be loaded into OPA as raw JSON data.
###############################################################################

producer_whitelist = {"large-fanout": {"fanout_producer"}}

consumer_whitelist = {"pii": {"pii_consumer"}}

topic_metadata = {
	"click-stream": {"tags": ["large-fanout"]},
	"credit-scores": {"tags": ["pii"]},
}

###############################################################################
# Helper rules for checking whitelists.
###############################################################################

topic_contains_pii {
	topic_metadata[topic_name].tags[_] == "pii"
}

topic_has_large_fanout {
	topic_metadata[topic_name].tags[_] == "large-fanout"
}

consumer_is_whitelisted_for_pii {
	consumer_whitelist.pii[_] == principal.name
}

producer_is_whitelisted_for_large_fanout {
	producer_whitelist["large-fanout"][_] == principal.name
}

###############################################################################
# Helper rules for input processing.
###############################################################################

is_write_operation {
	input.operation.name == "Write"
}

is_read_operation {
	input.operation.name == "Read"
}

is_topic_resource {
	input.resource.resourceType.name == "Topic"
}

topic_name = input.resource.name {
	is_topic_resource
}

principal = {"fqn": parsed.CN, "name": cn_parts[0]} {
	parsed := parse_user(urlquery.decode(input.session.sanitizedUser))
	cn_parts := split(parsed.CN, ".")
}

parse_user(user) = {key: value |
	parts := split(user, ",")
	[key, value] := split(parts[_], "=")
}
