package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "opa_mongo",
	Short: "A CLI tool to run OPA with mongo DB integration",
	Long: `This CLI tool integrates OPA with MongoDB. It leverages OPA's partial evaluation feature to translate Rego to a MongoDB query which can then be applied to the database.

Note: Before you run this CLI make sure you have working mongo database that is accessible through this cli.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
