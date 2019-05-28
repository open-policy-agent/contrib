package cmd

import (
	"os"
	"fmt"

	"github.com/spf13/cobra"
)

// Flags that are to be added to subset of commands.
var (
	opaURL    string
	opaAuth   string
	logLevel  string
	logFormat string
)

const banner  = `
                          _       _        _     _           
   ___  _ __   __ _      (_)_ __ | |_ __ _| |__ | | ___  ___ 
  / _ \| '_ \ / _' |_____| | '_ \| __/ _' | '_ \| |/ _ \/ __|
 | (_) | |_) | (_| |_____| | |_) | || (_| | |_) | |  __/\__ \
  \___/| .__/ \__,_|     |_| .__/ \__\__,_|_.__/|_|\___||___/
       |_|                 |_|                               
`

var rootCmd = &cobra.Command{
	Use:   "opa-iptables",
	Short: "opa-iptables enables to control iptables rules using OPA",
	Run: func(cmd *cobra.Command, args []string) {
		println(banner)
		cmd.Help()
	},
}

//Execute TODO
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
