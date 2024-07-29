# The package path. This should be passed with the display_endpoint flag
# in the PAM configuration file.
package display

import rego.v1

# regal ignore:rule-name-repeats-package
display_spec := [
	{
		"message": "Welcome to the OPA-PAM demonstration.",
		"style": "info",
	},
	{
		"message": "Please enter your last name: ",
		"style": "prompt_echo_on",
		"key": "last_name",
	},
	{
		"message": "Please enter your secret: ",
		"style": "prompt_echo_off",
		"key": "secret",
	},
]
