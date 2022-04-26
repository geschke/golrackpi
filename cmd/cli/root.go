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
 golrackpi is a small CLI application to read different values from Kostal Plenticore Inverters.
 `,
}

var authData golrackpi.AuthClient
var delimiter string = ","

var (
	outputCSV       bool   = false
	outputFile      string = ""
	outputTimestamp bool   = false
	outputAppend    bool   = false
	outputNoHeaders bool   = false
)

// init sets the global flags and their options.
func init() {
	rootCmd.PersistentFlags().StringVarP(&authData.Password, "password", "p", "", "Password (required)")
	rootCmd.PersistentFlags().StringVarP(&authData.Server, "server", "s", "", "Server (e.g. inverter IP address) (required)")
	rootCmd.PersistentFlags().StringVarP(&authData.Scheme, "scheme", "m", "", "Scheme (http or https, default http)")
	rootCmd.MarkPersistentFlagRequired("password")
	rootCmd.MarkPersistentFlagRequired("server")

}

// Exec is the entrypoint of the Cobra CLI library.
func Exec() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}

// getOutFile returns a pointer to an opened file if the corresponding flags are set.
// If the return value is nil, output should be sent to os.Stdout
func getOutFile() (*os.File, error) {
	var f *os.File
	var err error
	if len(outputFile) > 0 {
		if outputAppend {
			f, err = os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
		} else {
			f, err = os.Create(outputFile)
			if err != nil {
				return nil, err
			}
		}
		return f, nil
	}
	return nil, nil
}

// closeOutFile closes a file and handles errors
func closeOutFile(f *os.File) {
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
