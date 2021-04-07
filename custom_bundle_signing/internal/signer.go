package internal

import (
	"fmt"
	"strings"

	"github.com/open-policy-agent/opa/bundle"
)

// CustomSigner demonstrates a custom bundle signing implementation.
type CustomSigner struct{}

// GenerateSignedToken demonstrates how to implement the bundle.Signer interface,
// for the purpose of creating custom bundle signing. Note: In this example,
// no actual signing is taking place, it simply demonstrates how one could begin
// a custom signing implementation.
func (s *CustomSigner) GenerateSignedToken(files []bundle.FileInfo, sc *bundle.SigningConfig, keyID string) (string, error) {
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name)
	}
	// Generate a signature for the bundle files...
	// OPA uses JWTs to create bundle signatures, but we're skipping that for this example.
	return fmt.Sprintf("some_signature;%s", strings.Join(fileNames, ";")), nil
}
