package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/contrib/opa-iptables/pkg/version"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show version of CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("commit:",version.GitCommit)
		fmt.Println("version:",version.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
