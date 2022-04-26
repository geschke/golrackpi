// Copyright 2022 Ralf Geschke. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"io"
	"os"
	"time"

	"github.com/geschke/golrackpi"
	"github.com/spf13/cobra"
)

func init() {

	settingsModuleCmd.Flags().BoolVarP(&outputCSV, "csv", "c", false, "Set output to CSV format")
	settingsModuleCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")
	settingsModuleCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write output to file [filename]")
	settingsModuleCmd.Flags().BoolVarP(&outputTimestamp, "timestamp", "t", false, "Add timestamp to output")
	settingsModuleCmd.Flags().BoolVarP(&outputAppend, "append", "a", false, "Append output to file (default: overwrite content)")
	settingsModuleCmd.Flags().BoolVarP(&outputNoHeaders, "no-headers", "", false, "Omit headline in CSV output")

	settingsModuleSettingCmd.Flags().BoolVarP(&outputCSV, "csv", "c", false, "Set output to CSV format")
	settingsModuleSettingCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")
	settingsModuleSettingCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write output to file [filename]")
	settingsModuleSettingCmd.Flags().BoolVarP(&outputTimestamp, "timestamp", "t", false, "Add timestamp to output")
	settingsModuleSettingCmd.Flags().BoolVarP(&outputAppend, "append", "a", false, "Append output to file (default: overwrite content)")
	settingsModuleSettingCmd.Flags().BoolVarP(&outputNoHeaders, "no-headers", "", false, "Omit headline in CSV output")

	settingsModuleSettingsCmd.Flags().BoolVarP(&outputCSV, "csv", "c", false, "Set output to CSV format")
	settingsModuleSettingsCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")
	settingsModuleSettingsCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write output to file [filename]")
	settingsModuleSettingsCmd.Flags().BoolVarP(&outputTimestamp, "timestamp", "t", false, "Add timestamp to output")
	settingsModuleSettingsCmd.Flags().BoolVarP(&outputAppend, "append", "a", false, "Append output to file (default: overwrite content)")
	settingsModuleSettingsCmd.Flags().BoolVarP(&outputNoHeaders, "no-headers", "", false, "Omit headline in CSV output")

	rootCmd.AddCommand(settingsCmd)
	settingsCmd.AddCommand(settingsListCmd)
	settingsCmd.AddCommand(settingsModuleCmd)
	settingsCmd.AddCommand(settingsModuleSettingCmd)
	settingsCmd.AddCommand(settingsModuleSettingsCmd)

}

var settingsCmd = &cobra.Command{
	Use: "settings",

	Short: "List settings content",
	//Long:  ``,
	Run: func(cmd *cobra.Command,
		args []string) {
		handleSettings()
	},
}

var settingsListCmd = &cobra.Command{
	Use: "list",

	Short: "List all modules with their list of settings identifiers.",
	//Long:  ``,

	Run: func(cmd *cobra.Command,
		args []string) {
		listSettings()
	},
}

var settingsModuleCmd = &cobra.Command{
	Use: "module <moduleid>",

	Short: "Get module settings values.",
	//Long:  ``,

	Run: func(cmd *cobra.Command,
		args []string) {
		getSettingsModule(args)
	},
}

var settingsModuleSettingCmd = &cobra.Command{
	Use: "setting <moduleid> <settingid>",

	Short: "Get module setting value.",
	//Long:  ``,

	Run: func(cmd *cobra.Command,
		args []string) {
		getSettingsModuleSetting(args)
	},
}

var settingsModuleSettingsCmd = &cobra.Command{
	Use: "settings <moduleid> <settingids>",

	Short: "Get module settings values. Use a comma-separated list of settingids.",
	//Long:  ``,

	Run: func(cmd *cobra.Command,
		args []string) {
		getSettingsModuleSettings(args)
	},
}

// listSettings prints a (huge) list of module ids with their corresponding setting ids
func listSettings() {
	var outErr io.Writer = os.Stderr

	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()
	defer lib.Logout()

	if err != nil {
		fmt.Fprintln(outErr, "An error occurred:", err)
		return
	}

	settings, err := lib.Settings()

	if err != nil {
		fmt.Fprintln(outErr, "An error occurred:", err)
		return
	}
	for _, s := range settings {
		fmt.Println(s.ModuleId)
		for _, data := range s.Settings {
			fmt.Println("\t", data.Id)
		}
	}

}

// getSettingsModule takes a module id as argument and prints setting ids and their current values
func getSettingsModule(args []string) {
	var outErr io.Writer = os.Stderr

	if len(args) < 1 {
		fmt.Fprintln(outErr, "Please submit a moduleid.")
		return
	} else if len(args) > 1 {
		fmt.Fprintln(outErr, "Please submit only one moduleid.")
		return
	}

	moduleId := args[0]

	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()
	if err != nil {
		fmt.Fprintln(outErr, "An error occurred:", err)
		return
	}
	defer lib.Logout()

	values, err := lib.SettingsModule(moduleId)

	if err != nil {
		fmt.Fprintln(outErr, "An error occurred:", err)
		return
	}
	writeSettingsValues(values)
}

// getSettingsModuleSetting takes a module id and a setting id as arguments and prints setting ids and their current value
func getSettingsModuleSetting(args []string) {
	var outErr io.Writer = os.Stderr

	if len(args) < 2 {
		fmt.Fprintln(outErr, "Please submit a moduleid and a settingid.")
		return
	} else if len(args) > 2 {
		fmt.Fprintln(outErr, "Please submit only one moduleid with its settingid.")
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
	if err != nil {
		fmt.Fprintln(outErr, "An error occurred:", err)
		return
	}
	defer lib.Logout()

	values, err := lib.SettingsModuleSetting(moduleId, settingId)

	if err != nil {
		fmt.Fprintln(outErr, "An error occurred:", err)
		return
	}

	writeSettingsValues(values)

}

// getSettingsModuleSettings takes a module id and one or more setting ids as arguments and prints setting ids and their current value
func getSettingsModuleSettings(args []string) {
	var outErr io.Writer = os.Stderr

	if len(args) < 2 {
		fmt.Fprintln(outErr, "Please submit a moduleid and one or more settingids")
		return
	}

	settingIds := args[1:]
	moduleId := args[0]

	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()
	if err != nil {
		fmt.Fprintln(outErr, "An error occurred:", err)
		return
	}
	defer lib.Logout()

	values, err := lib.SettingsModuleSettings(moduleId, settingIds...)

	if err != nil {
		fmt.Fprintln(outErr, "An error occurred:", err)
		return
	}
	writeSettingsValues(values)

}

// writeSettingValues is a helper function to print a slice of setting ids and their value
func writeSettingsValues(values []golrackpi.SettingsValues) {

	var outErr io.Writer = os.Stderr
	var w io.Writer

	f, err := getOutFile()
	if err != nil {
		fmt.Fprintln(outErr, "Could not open file ", outputFile)
		return
	}
	if f != nil {
		w = f
		defer closeOutFile(f)
	} else {
		w = os.Stdout
	}

	if outputCSV {
		if !outputNoHeaders {
			if outputTimestamp {
				fmt.Fprintf(w, "Timestamp%s", delimiter)
			}
			fmt.Fprintf(w, "Id%sValue\n", delimiter)
		}
		for _, v := range values {
			if outputTimestamp {
				fmt.Fprintf(w, "%s%s", time.Now().Format(time.RFC3339), delimiter)
			}
			fmt.Fprintf(w, "%s%s%s\n", v.Id, delimiter, v.Value)
		}

	} else {
		if outputTimestamp {
			fmt.Fprintf(w, "Timestamp\t")
		}
		fmt.Fprintln(w, "Id\tValue")
		for _, v := range values {
			if outputTimestamp {
				fmt.Fprintf(w, "%s\t", time.Now().Format(time.RFC3339))
			}
			fmt.Fprintf(w, "%s\t%s\n", v.Id, v.Value)
		}

		fmt.Fprintln(w)
	}

}

/*
* Handle settings-related commands
 */
func handleSettings() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi settings --help to show available commands.")
}
