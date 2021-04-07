package main

import (
	"fmt"
	"os"

	"github.com/open-policy-agent/contrib/custom_bundle_signing/setup"
)

func main() {
	cmd := setup.SetupRootCommand(nil)

	if err := cmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
