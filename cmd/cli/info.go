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

	ok, err := lib.Login()

	fmt.Println("Ok?", ok)
	if err != nil {
		fmt.Println(err)
		panic(err.Error())
	}

	info := lib.Version()
	for k, v := range info {
		fmt.Println("key, value: ", k, v)
	}

	//fmt.Println("returned: ", pd)

	/*
		moduleNames := make([]string, 0, len(pd))
		for mn := range pd {
			moduleNames = append(moduleNames, mn)
		}

		// sort the slice by keys
		sort.Strings(moduleNames)

		for _, moduleId := range moduleNames {
			fmt.Println("Module:", moduleId)
			fmt.Println("ProcessDataIds:")
			for _, processDataIds := range pd[moduleId].ProcessDataIds {
				fmt.Println("\t", processDataIds)

			}
			fmt.Println()
		}

		//fmt.Printf("%-"+fmt.Sprintf("%d", maxStrlen)+"s%-21s%-21s\n", domainname, dtCreated, dtUpdated)
	*/
}

/*
* Handle processdata-related commands
 */
func handleInfo() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi info --help to show available commands.")
}
