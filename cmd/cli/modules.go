package cmd

import (
	"fmt"

	"github.com/geschke/golrackpi"
	"github.com/spf13/cobra"

	//"log"
	//"os"
	"sort"
	//"strings"
)

func init() {

	rootCmd.AddCommand(modulesCmd)
	modulesCmd.AddCommand(modulesListCmd)

}

var modulesCmd = &cobra.Command{
	Use: "modules",

	Short: "List modules content",
	//Long:  `Manage dynpower domain entries in database.`,
	Run: func(cmd *cobra.Command,
		args []string) {
		handleModules()
	},
}

var modulesListCmd = &cobra.Command{
	Use: "list",

	Short: "List all modules and their type",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

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

	moduleNames := make([]string, 0, len(modules))
	for mn := range modules {
		moduleNames = append(moduleNames, mn)
	}

	// sort the slice by keys
	sort.Strings(moduleNames)

	fmt.Println("Moduleid, Type")
	for _, moduleId := range moduleNames {
		fmt.Printf("%s %s\n", moduleId, modules[moduleId].Type)

	}

}

/*
* Handle processdata-related commands
 */
func handleModules() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi modules --help to show available commands.")
}
