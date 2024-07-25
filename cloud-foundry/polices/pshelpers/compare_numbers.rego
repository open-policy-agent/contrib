package pshelpers

import rego.v1

limit_exceeded(limit, set) if {
	some item in set
	item > limit
}
