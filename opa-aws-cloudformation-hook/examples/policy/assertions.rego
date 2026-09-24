# METADATA
# description: Utility functions for working with Rego unit tests
# authors:
#   - Anders Eknert <anders@eknert.com>
#
package assertions

import rego.v1

# METADATA
# description: Assert expected is equal to result, or fail while printing both inputs to console
assert_equals(expected, result) if expected == result

assert_equals(expected, result) := false if {
	expected != result
	print("expected equals:", _quote_if_string(expected), "got:", result)
}

# METADATA
# description: Assert expected is not equal to result, or fail while printing both inputs to console
assert_not_equals(expected, result) if expected != result

assert_not_equals(expected, result) := false if {
	expected == result
	print("expected not equals:", _quote_if_string(expected), "got:", result)
}

# METADATA
# description: Assert item is in coll, or fail while printing the collection to console
assert_in(item, coll) if item in coll

assert_in(item, coll) := false if {
	not item in coll
	print("expected", type_name(item), _quote_if_string(item), "in", type_name(coll), "got:", coll)
}

# METADATA
# description: Assert item is not in coll, or fail while printing the collection to console
assert_not_in(item, coll) if not item in coll

assert_not_in(item, coll) := false if {
	item in coll
	print("expected", type_name(item), _quote_if_string(item), "not in", type_name(coll), "got:", coll)
}

# METADATA
# description: Assert provided collection is empty, or fail while printing the collection to console
assert_empty(coll) if count(coll) == 0

assert_empty(coll) := false if {
	count(coll) != 0
	print("expected empty", type_name(coll), "got:", coll)
}

# METADATA
# description: Fail with provided message
fail(msg) if {
	print(msg)
	false # regal ignore:constant-condition
}

_quote_if_string(x) := concat("", [`"`, x, `"`]) if is_string(x)

_quote_if_string(x) := x if not is_string(x)
