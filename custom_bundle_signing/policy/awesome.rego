package awesome

import rego.v1

# regal ignore:rule-name-repeats-package
default awesome := false

# regal ignore:rule-name-repeats-package
awesome if {
	input.awesome
}
