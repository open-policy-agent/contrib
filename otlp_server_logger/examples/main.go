// Copyright 2026 The OPA Authors.  All rights reserved.
// Use of this source code is governed by an Apache2
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"os"

	"github.com/open-policy-agent/opa/cmd"
	"github.com/open-policy-agent/opa/runtime"

	"github.com/open-policy-agent/contrib/otlp_server_logger"
)

func main() {
	runtime.RegisterPlugin(otlpserverlogger.PluginName, otlpserverlogger.Factory{})

	if err := cmd.RootCommand.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
