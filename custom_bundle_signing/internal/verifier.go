package internal

import (
	"fmt"
	"strings"

	"github.com/open-policy-agent/opa/bundle"
)

// CustomVerifier demonstrates a custom bundle verification implementation.
type CustomVerifier struct{}

// VerifyBundleSignature demonstrates how to implement the bundle.Verifier interface,
// for the purpose of creating custom bundle verification. Note: In this example,
// no actual verification is taking place, it simply demonstrates how one could
// begin a custom verification implementation.
func (v *CustomVerifier) VerifyBundleSignature(sc bundle.SignaturesConfig, bvc *bundle.VerificationConfig) (map[string]bundle.FileInfo, error) {
	files := make(map[string]bundle.FileInfo)

	if len(sc.Signatures) == 0 {
		return files, fmt.Errorf(".signatures.json: missing signature (expected exactly one)")
	}

	if len(sc.Signatures) > 1 {
		return files, fmt.Errorf(".signatures.json: multiple sgnatures not supported (expected exactly one)")
	}

	for _, signature := range sc.Signatures {
		if !strings.HasPrefix(signature, "some_signature") {
			return files, fmt.Errorf("unexpected signature")
		}
		parts := strings.Split(signature, ";")
		for _, file := range parts[1:] {
			// Normally, the file info is a part of the signature, in the
			// form of JWT data, but we are skipping that for this example.
			// Instead, we're just generating the FileInfo manually
			// and skipping the hash verification for all the files in the
			// example bundle.
			files[file] = bundle.FileInfo{
				Name:      file,
				Algorithm: "SHA-256",
			}
		}
	}

	return files, nil
}
