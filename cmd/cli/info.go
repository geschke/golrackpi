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

	//processdataGetCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
	//processdataGetCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")

	rootCmd.AddCommand(infoCmd)
	infoCmd.AddCommand(infoVersionCmd)

}

var infoCmd = &cobra.Command{
	Use: "info",

	Short: "Returns miscellaneous information",
	//Long:  `Manage dynpower domain entries in database.`,
	Run: func(cmd *cobra.Command,
		args []string) {
		handleInfo()
	},
}

var infoVersionCmd = &cobra.Command{
	Use: "version",

	Short: "Returns information about the API",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		infoVersion()
	},
}

func infoVersion() {
	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()

	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	info, err := lib.Version()
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	for k, v := range info {
		fmt.Printf("%s: %v\n", k, v)
	}

}

/*
* Handle processdata-related commands
 */
func handleInfo() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi info --help to show available commands.")
}
