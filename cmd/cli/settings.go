package cmd

import (
	"fmt"

	"github.com/geschke/golrackpi"
	"github.com/spf13/cobra"

	"strings"
)

func init() {

	settingsModuleCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
	settingsModuleCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")
	settingsModuleSettingCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
	settingsModuleSettingCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")
	settingsModuleSettingsCmd.Flags().BoolVarP(&csvOutput, "csv", "c", false, "Set output to CSV format")
	settingsModuleSettingsCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")

	rootCmd.AddCommand(settingsCmd)
	settingsCmd.AddCommand(settingsListCmd)
	settingsCmd.AddCommand(settingsModuleCmd)
	settingsCmd.AddCommand(settingsModuleSettingCmd)
	settingsCmd.AddCommand(settingsModuleSettingsCmd)

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
	Use: "module <moduleid>",

	Short: "Get module settings values.",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		getSettingsModule(args)
	},
}

var settingsModuleSettingCmd = &cobra.Command{
	Use: "setting <moduleid> <settingid>",

	Short: "Get module setting value.",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		getSettingsModuleSetting(args)
	},
}

var settingsModuleSettingsCmd = &cobra.Command{
	Use: "settings <moduleid> <settingids>",

	Short: "Get module settings values. Use a comma-seperated list of settingids.",
	//Long:  `List all domains in the dynpower database. If a DSN is submitted by the flag --dsn, this DSN will be used. If no DSN is provided, dynpower-cli tries to use the environment variables DBHOST, DBUSER, DBNAME and DBPASSWORD.`,

	Run: func(cmd *cobra.Command,
		args []string) {
		getSettingsModuleSettings(args)
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

func getSettingsModule(args []string) {

	if len(args) < 1 {
		fmt.Println("Please submit a moduleid.")
		return
	} else if len(args) > 1 {
		fmt.Println("Please submit only one moduleid.")
		return
	}

	moduleId := args[0]

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

	values, err := lib.SettingsModule(moduleId)

	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}
	writeSettingValues(values)
}

func writeSettingValues(values []golrackpi.SettingsValues) {

	if csvOutput {
		fmt.Printf("Id%sValue\n", delimiter)
		for _, v := range values {
			fmt.Printf("%s%s%s\n", v.Id, delimiter, v.Value)
		}
	} else {
		fmt.Println("Id\tValue")
		for _, v := range values {
			fmt.Printf("%s\t%s\n", v.Id, v.Value)
		}

	}

}

func getSettingsModuleSetting(args []string) {

	if len(args) < 2 {
		fmt.Println("Please submit a moduleid and a settingid.")
		return
	} else if len(args) > 2 {
		fmt.Println("Please submit only one moduleid with its settingid.")
		return
	}

	moduleId := args[0]
	settingId := args[1]

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

	values, err := lib.SettingsModuleSetting(moduleId, settingId)

	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}

	writeSettingValues(values)

}

func getSettingsModuleSettings(args []string) {

	if len(args) < 2 {
		fmt.Println("Please submit a moduleid and s comma-seperated list of settingids.")
		return
	} else if len(args) > 2 {
		fmt.Println("Please submit only one moduleid with its settingids as comma-seperated list.")
		return
	}

	settingIds := strings.Split(args[1], ",")
	for _, settingId := range settingIds {
		fmt.Println("SettingId:", settingId)
	}

	moduleId := args[0]

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

	values, err := lib.SettingsModuleSettings(moduleId, settingIds)

	if err != nil {
		fmt.Println("An error occurred:", err)
		return
	}
	writeSettingValues(values)

}

/*
* Handle settings-related commands
 */
func handleSettings() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi settings --help to show available commands.")
}
