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

	rootCmd.AddCommand(infoCmd)
	infoCmd.AddCommand(infoVersionCmd)
	infoCmd.AddCommand(infoMeCmd)
	infoCmd.AddCommand(checkLoginLogoutCmd)

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

var infoMeCmd = &cobra.Command{
	Use: "me",

	Short: "Returns information about the user",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		infoMe()
	},
}

func infoMe() {
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

	info, err := lib.Me()
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	for k, v := range info {
		fmt.Printf("%s: %v\n", k, v)
	}

}

var checkLoginLogoutCmd = &cobra.Command{
	Use: "checklog",

	Short: "Check login and logout process",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		checkLoginLogout()
	},
}

func checkLoginLogout() {
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

	fmt.Println("logged in!")

	info, err := lib.Me()
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	for k, v := range info {
		fmt.Printf("%s: %v\n", k, v)
	}

	loggedOut, err := lib.Logout()
	fmt.Println("logout? ", loggedOut)
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}
	fmt.Println("Logged out!")

	info, err = lib.Me()
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
