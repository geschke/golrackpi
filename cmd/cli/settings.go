package cmd

import (
	"fmt"

	"github.com/geschke/golrackpi"
	"github.com/spf13/cobra"
	//"log"
	//"os"
	//"sort"
	//"strings"
)

func init() {
	/*	rootCmd.PersistentFlags().StringVarP(&authData.Password, "password", "p", "", "Password")
		rootCmd.PersistentFlags().StringVarP(&authData.Server, "server", "s", "", "Server (e.g. inverter IP address)")
		rootCmd.PersistentFlags().StringVarP(&authData.Scheme, "scheme", "m", "", "Scheme (http or https, default http)")

		processdataGetCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
		processdataGetCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")
	*/
	rootCmd.AddCommand(settingsCmd)
	settingsCmd.AddCommand(settingsListCmd)
	settingsCmd.AddCommand(settingsModuleCmd)

}

var settingsCmd = &cobra.Command{
	Use: "settings",

	Short: "List settings content",
	//Long:  `Manage dynpower domain entries in database.`,
	Run: func(cmd *cobra.Command,
		args []string) {
		handleSettings()
	},
}

var settingsListCmd = &cobra.Command{
	Use: "list",

	Short: "List all modules with their list of settings identifiers.",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		listSettings()
	},
}

var settingsModuleCmd = &cobra.Command{
	Use: "module",

	Short: "List ...",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		listSettingsModule()
	},
}

func listSettings() {
	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()
	defer lib.Logout()

	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	settings, err := lib.Settings()

	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}
	for _, s := range settings {
		fmt.Println(s.ModuleId)
		for _, data := range s.Settings {
			fmt.Println("\t", data.Id)
		}
	}

}

func listSettingsModule() {
	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()
	defer lib.Logout()

	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	values, err := lib.SettingsModule("scb:network")

	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}
	for _, v := range values {
		fmt.Println(v.Id, v.Value)
	}

}

/*
* Handle processdata-related commands
 */
func handleSettings() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi settings --help to show available commands.")
}
