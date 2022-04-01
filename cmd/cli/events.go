// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"

	"github.com/geschke/golrackpi"
	"github.com/spf13/cobra"
)

var language string
var max int

func init() {

	eventsCustomCmd.Flags().StringVarP(&language, "language", "l", "", "Language identifier, e.g. en-gb, de-de, fr-fr, ...")
	eventsCustomCmd.Flags().IntVarP(&max, "max", "x", 0, "Maximum number of events to return (default: 10)")

	eventsCustomCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
	eventsCustomCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")

	eventsLatestCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
	eventsLatestCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")

	rootCmd.AddCommand(eventsCmd)
	eventsCmd.AddCommand(eventsCustomCmd)
	eventsCmd.AddCommand(eventsLatestCmd)

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

var eventsCustomCmd = &cobra.Command{
	Use: "custom",

	Short: "Get the latest events, customized by language and number of returned events",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		latestCustomEvents()
	},
}

func latestCustomEvents() {

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

	//fmt.Println("language", language)
	//fmt.Println("max:", max)

	events, err := lib.EventsWithParam(language, max)

	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	if csvOutput {
		fmt.Printf("Description%sCategory%sLongDescription%sStartTime%sGroup%sEndTime%sCode%sIsActive\n", delimiter, delimiter, delimiter, delimiter, delimiter, delimiter, delimiter)
		for _, event := range events {
			fmt.Printf("%s%s%s%s%s%s%s%s%s%s%s%s%d%s%t\n", event.Description, delimiter, event.Category, delimiter, event.LongDescription, delimiter, event.StartTime, delimiter, event.Group, delimiter, event.EndTime, delimiter, event.Code, delimiter, event.IsActive)
		}
	} else {
		fmt.Println("Description\tCategory\tLongDescription\tStartTime\tGroup\tEndTime\tCode\tIsActive")
		for _, event := range events {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%t\n", event.Description, event.Category, event.LongDescription, event.StartTime, event.Group, event.EndTime, event.Code, event.IsActive)
		}

	}
}

func latestEvents() {

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

	events, err := lib.Events()
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	if csvOutput {
		fmt.Printf("Description%sCategory%sLongDescription%sStartTime%sGroup%sEndTime%sCode%sIsActive\n", delimiter, delimiter, delimiter, delimiter, delimiter, delimiter, delimiter)
		for _, event := range events {
			fmt.Printf("%s%s%s%s%s%s%s%s%s%s%s%s%d%s%t\n", event.Description, delimiter, event.Category, delimiter, event.LongDescription, delimiter, event.StartTime, delimiter, event.Group, delimiter, event.EndTime, delimiter, event.Code, delimiter, event.IsActive)
		}
	} else {
		fmt.Println("Description\tCategory\tLongDescription\tStartTime\tGroup\tEndTime\tCode\tIsActive")
		for _, event := range events {
			fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%d\t%t\n", event.Description, event.Category, event.LongDescription, event.StartTime, event.Group, event.EndTime, event.Code, event.IsActive)
		}
	}
}

/*
* Handle events-related commands
 */
func handleEvents() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi events --help to show available commands.")
}
