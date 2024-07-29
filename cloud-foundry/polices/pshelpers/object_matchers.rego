package pshelpers

import rego.v1

path_value_match(search_object, search_path, search_value) if {
	walk(search_object, [search_path, search_value])
}
