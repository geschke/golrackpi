package cmd

import (
	"fmt"

	"github.com/geschke/golrackpi"
	"github.com/spf13/cobra"

	//"log"
	//"os"
	"sort"
	"strings"
)

var authData golrackpi.AuthClient
var csvOutput bool = false
var delimiter string = ","

func init() {
	rootCmd.PersistentFlags().StringVarP(&authData.Password, "password", "p", "", "Password")
	rootCmd.PersistentFlags().StringVarP(&authData.Server, "server", "s", "", "Server (e.g. inverter IP address)")
	rootCmd.PersistentFlags().StringVarP(&authData.Scheme, "scheme", "m", "", "Scheme (http or https, default http)")

	processdataGetCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
	processdataGetCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")

	rootCmd.AddCommand(processdataCmd)
	processdataCmd.AddCommand(processdataListCmd)
	processdataCmd.AddCommand(processdataGetCmd)

}

var processdataCmd = &cobra.Command{
	Use: "processdata",

	Short: "List processdata values",
	//Long:  `Manage dynpower domain entries in database.`,
	Run: func(cmd *cobra.Command,
		args []string) {
		handleProcessdata()
	},
}

var processdataListCmd = &cobra.Command{
	Use: "list",

	Short: "List all available processdata and modules",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		listProcessdata()
	},
}

var processdataGetCmd = &cobra.Command{
	Use: "get [moduleid] [processdataid(s)] or get [moduleid|processdataid(s)] [moduleid|processdataid(s)] ... ",

	Short: "Get processdata values",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command,
		args []string) {
		getProcessdata(args)
	},
}

func listProcessdata() {
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

	defer lib.Logout()

	processData, err := lib.ProcessData()
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	for _, pdItem := range processData {
		fmt.Println("ModuleId:", pdItem.ModuleId)
		if len(pdItem.ProcessDataIds) > 0 {
			fmt.Println("ProcessDataIds:")
			for _, pdId := range pdItem.ProcessDataIds {
				fmt.Println("\t", pdId)
			}
		} else {
			fmt.Println("No ProcessDataId found.")
		}

	}

}

func getProcessdata(args []string) {

	// check format of submitted arguments
	var requestProcessData []golrackpi.ProcessData

	if strings.Contains(args[0], "|") { // search "|"" separator to request one or more modules with their processdataids
		for _, argModuleProcessdata := range args {
			moduleProcessdata := strings.Split(argModuleProcessdata, "|")
			if len(moduleProcessdata) != 2 {
				fmt.Println("Wrong format of moduleid and processdataid values.")
				return
			}
			argModuleId := moduleProcessdata[0]
			processdataIds := strings.Split(moduleProcessdata[1], ",")
			v := golrackpi.ProcessData{ModuleId: argModuleId, ProcessDataIds: processdataIds}
			requestProcessData = append(requestProcessData, v)

		}

	} else if len(args) == 2 { // else moduleid and processdataids must submitted separately
		moduleIds := strings.Split(args[0], ",")
		processdataIds := strings.Split(args[1], ",")
		fmt.Println("moduleids:", moduleIds)
		fmt.Println("processdataids:", processdataIds)

		if len(moduleIds) > 1 {
			fmt.Println("Please enter only one moduleid.")
			return
		}
		v := golrackpi.ProcessData{ModuleId: moduleIds[0], ProcessDataIds: processdataIds}
		requestProcessData = append(requestProcessData, v)

	} else {
		fmt.Println("Please submit module and processdata in an appropriate format.")
		return
	}

	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()
	if err != nil {
		fmt.Println("An error occurred:", err)
		//panic(err.Error())
		return
	}
	defer lib.Logout()

	pdv := lib.GetProcessDataValues(requestProcessData)
	fmt.Println("processDataValues:", pdv)

	moduleNames := make([]string, 0, len(pdv))
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
	}
}

/*
* Handle processdata-related commands
 */
func handleProcessdata() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi processdata --help to show available commands.")
}
