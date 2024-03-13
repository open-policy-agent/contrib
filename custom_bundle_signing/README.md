# Custom Bundle Signing Example

This example demonstrates how to implement custom bundle signing and verification.

Starting at `cmd/opa/main.go`, we hook into OPA's RootCommand and inject a PersistentPreRun 
hook for certain OPA commands. We call `bundle.RegisterSigner` and `bundle.RegisterVerifier` 
for our custom implementations of the `bundle.Signer` and `bundle.Verifier` interfaces, respectively.

Our [Signer](custom_bundle_signing/internal/signer.go) and [Verifier](custom_bundle_signing/internal/verification.go) 
implementations demonstrate how to implement the interfaces, without actually performing any real signing or verification. 
To create a functional Signer and Verifier, you'll want to look at the [DefaultSigner](https://github.com/godaddy/opa/blob/custom-sign-verify/bundle/sign.go#L46) 
and [DefaultVerifier](https://github.com/godaddy/opa/blob/custom-sign-verify/bundle/verify.go#L57) implementations in OPA's source.

Use `make build` and `make run` to run this example locally.
