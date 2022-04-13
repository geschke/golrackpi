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

	rootCmd.AddCommand(infoCmd)
	infoCmd.AddCommand(infoVersionCmd)
	infoCmd.AddCommand(infoMeCmd)
	infoCmd.AddCommand(checkLoginLogoutCmd)

}

var infoCmd = &cobra.Command{
	Use: "info",

	Short: "Returns miscellaneous information",
	//Long:  ``,
	Run: func(cmd *cobra.Command,
		args []string) {
		handleInfo()
	},
}

var infoVersionCmd = &cobra.Command{
	Use: "version",

	Short: "Returns information about the API",
	//Long:  ``,

	Run: func(cmd *cobra.Command,
		args []string) {
		infoVersion()
	},
}

var infoMeCmd = &cobra.Command{
	Use: "me",

	Short: "Returns information about the user",
	//Long:  ``,

	Run: func(cmd *cobra.Command,
		args []string) {
		infoMe()
	},
}

var checkLoginLogoutCmd = &cobra.Command{
	Use: "checklog",

	Short: "Check login and logout process, prints information about the user after login and logout",
	//Long:  ``,

	Run: func(cmd *cobra.Command,
		args []string) {
		checkLoginLogout()
	},
}

// infoVersion prints information about the API (i.e. hostname, api version...)
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
	defer lib.Logout()

	info, err := lib.Version()
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	for k, v := range info {
		fmt.Printf("%s: %v\n", k, v)
	}

}

// infoMe prints information about the user
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
	defer lib.Logout()

	info, err := lib.Me()
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	for k, v := range info {
		fmt.Printf("%s: %v\n", k, v)
	}

}

// checkLoginLogout checks login and logout process. It prints information from the "me" request with values about the user after login and logout.
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

	fmt.Println("Logged in!")

	info, err := lib.Me()
	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	for k, v := range info {
		fmt.Printf("%s: %v\n", k, v)
	}

	_, err = lib.Logout()

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

// Handle info-related commands
func handleInfo() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi info --help to show available commands.")
}
