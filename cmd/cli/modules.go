// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/geschke/golrackpi"
	"github.com/spf13/cobra"
)

func init() {
	modulesListCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
	modulesListCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")

	rootCmd.AddCommand(modulesCmd)
	modulesCmd.AddCommand(modulesListCmd)

}

var modulesCmd = &cobra.Command{
	Use: "modules",

	Short: "List modules content",
	//Long:  `...`,
	Run: func(cmd *cobra.Command,
		args []string) {
		handleModules()
	},
}

var modulesListCmd = &cobra.Command{
	Use: "list",

	Short: "List all modules and their type",
	//Long:  ``,

	Run: func(cmd *cobra.Command,
		args []string) {
		listModules()
	},
}

func listModules() {
	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	ok, err := lib.Login()
	defer lib.Logout()

	fmt.Println("Ok?", ok)
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	modules, err := lib.Modules()
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	if csvOutput {
		fmt.Printf("ModuleId%sType\n", delimiter)
		for _, module := range modules {
			fmt.Printf("%s%s%s\n", module.Id, delimiter, module.Type)

		}

	} else {

		fmt.Println("Moduleid\tType")
		for _, module := range modules {
			fmt.Printf("%s\t%s\n", module.Id, module.Type)

		}
	}
}

/*
* Handle processdata-related commands
 */
func handleModules() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi modules --help to show available commands.")
}
