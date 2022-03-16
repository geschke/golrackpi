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

	eventsLatestCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
	eventsLatestCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")

	rootCmd.AddCommand(eventsCmd)
	eventsCmd.AddCommand(eventsLatestCmd)
	//processdataCmd.AddCommand(processdataGetCmd)

}

var eventsCmd = &cobra.Command{
	Use: "events",

	Short: "Get the latest events",
	//Long:  `Manage dynpower domain entries in database.`,
	Run: func(cmd *cobra.Command,
		args []string) {
		handleEvents()
	},
}

var eventsLatestCmd = &cobra.Command{
	Use: "latest",

	Short: "Get the latest events",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		latestEvents()
	},
}

func latestEvents() {

	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	ok, err := lib.Login()
	defer lib.Logout()

	fmt.Println("Ok?", ok)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	events, err := lib.Events()

	fmt.Println("events:", events)

	/*	moduleNames := make([]string, 0, len(pdv))
		for mn := range pdv {
			moduleNames = append(moduleNames, mn)
		}

		// sort the slice by keys
		sort.Strings(moduleNames)

		if csvOutput {
			fmt.Printf("Module%sProcessdata Id%sProcessdata Unit%sProcessdata Value\n", delimiter, delimiter, delimiter)
			for _, moduleId := range moduleNames {

				for _, processData := range pdv[moduleId].ProcessData {
					fmt.Printf("%s%s%s%s%s%s%v\n", moduleId, delimiter, processData.Id, delimiter, processData.Unit, delimiter, processData.Value)

				}

			}
		} else {

			for _, moduleId := range moduleNames {
				fmt.Println("Module:", moduleId)
				fmt.Println("ProcessDataValues (Id\tUnit\tValue):")
				for _, processData := range pdv[moduleId].ProcessData {
					fmt.Println(processData.Id, "\t", processData.Unit, "\t", processData.Value)
					// todo: add better formatting

				}
				fmt.Println()
			}
		}*/
}

/*
* Handle events-related commands
 */
func handleEvents() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi events --help to show available commands.")
}
