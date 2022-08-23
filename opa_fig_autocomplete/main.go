package main

import (
	"fmt"

	"github.com/open-policy-agent/opa/cmd"
	genFigSpec "github.com/withfig/autocomplete-tools/packages/cobra"
)

func main() {
	root := cmd.RootCommand
	root.Use = "opa"
	s := genFigSpec.MakeFigSpec(root)
	fmt.Println(s.ToTypescript())
}
