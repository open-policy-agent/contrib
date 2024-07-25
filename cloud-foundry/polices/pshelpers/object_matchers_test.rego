package pshelpers_test

import rego.v1

import data.pshelpers

test_path_value_match_with_correct_path_value if {
	input_object = {"random": {"stuff": "moo"}}
	path := ["random", "stuff"]
	value := "moo"

	pshelpers.path_value_match(input_object, path, value)
}

test_path_value_match_incorrect_path if {
	input_object = {"random": {"stuff": "moo"}}
	path := ["incorrect", "path"]
	value := "moo"

	not pshelpers.path_value_match(input_object, path, value)
}

test_path_value_match_incorrect_value if {
	input_object = {"random": {"stuff": "moo"}}
	path := ["random", "stuff"]
	value := "incorrect"

	not pshelpers.path_value_match(input_object, path, value)
}
