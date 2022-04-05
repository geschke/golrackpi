// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/geschke/golrackpi"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golrackpi",
	Short: "Go Client for Kostal Plenticore Inverters",
	Long: `
 golrackpi is a small CLI application to read values from Kostal Plenticore Inverters.
 `,
}

var authData golrackpi.AuthClient
var csvOutput bool = false
var delimiter string = ","
var outputFile string = ""
var outputTimestamp bool = false
var outputAppend bool = false
var outputNoHeaders bool = false

func init() {
	rootCmd.PersistentFlags().StringVarP(&authData.Password, "password", "p", "", "Password (required)")
	rootCmd.PersistentFlags().StringVarP(&authData.Server, "server", "s", "", "Server (e.g. inverter IP address) (required)")
	rootCmd.PersistentFlags().StringVarP(&authData.Scheme, "scheme", "m", "", "Scheme (http or https, default http)")
	rootCmd.MarkPersistentFlagRequired("password")
	rootCmd.MarkPersistentFlagRequired("server")

}

func Exec() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	/*err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}*/
}
