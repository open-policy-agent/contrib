module github.com/open-policy-agent/contrib/custom_bundle_signing

go 1.15

require (
	github.com/open-policy-agent/opa v0.27.1
	github.com/spf13/cobra v1.1.3
)

replace github.com/open-policy-agent/opa v0.27.1 => github.com/godaddy/opa v0.23.1-0.20210406192957-5eb38445fa97
