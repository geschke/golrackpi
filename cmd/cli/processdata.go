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

	"strings"
)

func init() {

	processdataModuleCmd.Flags().BoolVarP(&outputCSV, "csv", "c", false, "Set output to CSV format")
	processdataModuleCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")
	processdataModuleCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write output to file [filename]")
	processdataModuleCmd.Flags().BoolVarP(&outputTimestamp, "timestamp", "t", false, "Add timestamp to output")
	processdataModuleCmd.Flags().BoolVarP(&outputAppend, "append", "a", false, "Append output to file (default: overwrite content)")
	processdataModuleCmd.Flags().BoolVarP(&outputNoHeaders, "no-headers", "", false, "Omit headline in CSV output")

	processdataGetCmd.Flags().BoolVarP(&outputCSV, "csv", "c", false, "Set output to CSV format")
	processdataGetCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")
	processdataGetCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write output to file [filename]")
	processdataGetCmd.Flags().BoolVarP(&outputTimestamp, "timestamp", "t", false, "Add timestamp to output")
	processdataGetCmd.Flags().BoolVarP(&outputAppend, "append", "a", false, "Append output to file (default: overwrite content)")
	processdataGetCmd.Flags().BoolVarP(&outputNoHeaders, "no-headers", "", false, "Omit headline in CSV output")

	processdataMultCmd.Flags().BoolVarP(&outputCSV, "csv", "c", false, "Set output to CSV format")
	processdataMultCmd.Flags().StringVarP(&delimiter, "delimiter", "d", ",", "Set CSV delimiter (default \",\")")
	processdataMultCmd.Flags().StringVarP(&outputFile, "output-file", "o", "", "Write output to file [filename]")
	processdataMultCmd.Flags().BoolVarP(&outputTimestamp, "timestamp", "t", false, "Add timestamp to output")
	processdataMultCmd.Flags().BoolVarP(&outputAppend, "append", "a", false, "Append output to file (default: overwrite content)")
	processdataMultCmd.Flags().BoolVarP(&outputNoHeaders, "no-headers", "", false, "Omit headline in CSV output")

	rootCmd.AddCommand(processdataCmd)
	processdataCmd.AddCommand(processdataListCmd)
	processdataCmd.AddCommand(processdataModuleCmd)
	processdataCmd.AddCommand(processdataMultCmd)
	processdataCmd.AddCommand(processdataGetCmd)

}

var processdataCmd = &cobra.Command{
	Use: "processdata",

	Short: "List processdata values",
	//Long:  ``,
	Run: func(cmd *cobra.Command,
		args []string) {
		handleProcessdata()
	},
}

var processdataListCmd = &cobra.Command{
	Use: "list",

	Short: "List all available modules and processdata identifiers",
	//Long:  ``,

	Run: func(cmd *cobra.Command,
		args []string) {
		listProcessdata()
	},
}

var processdataModuleCmd = &cobra.Command{
	Use: "module [moduleid]",

	Short: "Get all processdata values of the specified moduleid",
	//Long:  ``,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command,
		args []string) {
		getModuleProcessdata(args)
	},
}

var processdataMultCmd = &cobra.Command{
	Use: "mult [moduleid] [processdataid(s)] or mult [moduleid|processdataid(s)] [moduleid|processdataid(s)] ... ",

	Short: "Get one or more modules with their processdata values",
	//Long:  ``,
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command,
		args []string) {
		getMultProcessdata(args)
	},
}

var processdataGetCmd = &cobra.Command{
	Use: "get [moduleid] [processdataid(s)]",

	Short: "Get module with one or more of its processdata values",
	//Long:  ``,
	Args: cobra.MinimumNArgs(2),
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

func getMultProcessdata(args []string) {

	// check format of submitted arguments
	var requestProcessData []golrackpi.ProcessData
	var errOut io.Writer = os.Stderr
	var w io.Writer

	f, errFile := getOutFile()
	if errFile != nil {
		fmt.Fprintln(errOut, "Could not open file ", outputFile)
		return
	}
	if f != nil {
		w = f
		defer closeOutFile(f)
	} else {
		w = os.Stdout
	}

	if strings.Contains(args[0], "|") { // search "|"" separator to request one or more modules with their processdataids
		for _, argModuleProcessdata := range args {
			moduleProcessdata := strings.Split(argModuleProcessdata, "|")
			if len(moduleProcessdata) != 2 {
				fmt.Fprintln(errOut, "Wrong format of moduleid and processdataid values.")
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

		if len(moduleIds) > 1 {
			fmt.Fprintln(errOut, "Please enter only one moduleid.")
			return
		}
		v := golrackpi.ProcessData{ModuleId: moduleIds[0], ProcessDataIds: processdataIds}
		requestProcessData = append(requestProcessData, v)

	} else {
		fmt.Fprintln(errOut, "Please submit module and processdata in an appropriate format.")
		return
	}

	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()
	if err != nil {
		fmt.Fprintln(errOut, "An error occurred:", err)

		return
	}
	defer lib.Logout()

	processDataValues, err := lib.ProcessDataValues(requestProcessData)
	if err != nil {
		fmt.Fprintln(errOut, "An error occurred:", err)
		return
	}

	if outputCSV {
		if !outputNoHeaders {
			if outputTimestamp {
				fmt.Fprintf(w, "Timestamp%s", delimiter)
			}
			fmt.Fprintf(w, "Module%sProcessdata Id%sProcessdata Unit%sProcessdata Value\n", delimiter, delimiter, delimiter)
		}
		for _, pdv := range processDataValues {
			for _, pd := range pdv.ProcessData {

				if outputTimestamp {
					fmt.Fprintf(w, "%s%s", time.Now().Format(time.RFC3339), delimiter)
				}
				fmt.Fprintf(w, "%s%s%s%s%s%s%v\n", pdv.ModuleId, delimiter, pd.Id, delimiter, pd.Unit, delimiter, pd.Value)

			}
		}

	} else {
		if outputTimestamp {
			fmt.Fprintf(w, "Timestamp:\t")
		}
		for _, pdv := range processDataValues {
			if outputTimestamp {
				fmt.Fprintf(w, "%s\n", time.Now().Format(time.RFC3339))
			}
			fmt.Fprintln(w, "Module:", pdv.ModuleId)
			for _, pd := range pdv.ProcessData {
				fmt.Fprintln(w, pd.Id, "\t", pd.Unit, "\t", pd.Value)
			}
			fmt.Fprintln(w)
		}

	}
}

func getProcessdata(args []string) {
	var errOut io.Writer = os.Stderr
	var w io.Writer

	f, errFile := getOutFile()
	if errFile != nil {
		fmt.Fprintln(errOut, "Could not open file ", outputFile)
		return
	}
	if f != nil {
		w = f
		defer closeOutFile(f)
	} else {
		w = os.Stdout
	}
	// submitted values: moduleid pdid pdid2 pdid3...

	moduleId := args[0]
	processDataIds := args[1:]

	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()
	if err != nil {
		fmt.Fprintln(errOut, "An error occurred:", err)

		return
	}
	defer lib.Logout()

	processDataValues, err := lib.ProcessDataModuleValues(moduleId, processDataIds...)
	if err != nil {
		fmt.Fprintln(errOut, "An error occurred:", err)
		return
	}

	if outputCSV {
		if !outputNoHeaders {
			if outputTimestamp {
				fmt.Fprintf(w, "Timestamp%s", delimiter)
			}
			fmt.Fprintf(w, "Module%sProcessdata Id%sProcessdata Unit%sProcessdata Value\n", delimiter, delimiter, delimiter)
		}
		for _, pdv := range processDataValues {
			for _, pd := range pdv.ProcessData {
				if outputTimestamp {
					fmt.Fprintf(w, "%s%s", time.Now().Format(time.RFC3339), delimiter)
				}
				fmt.Fprintf(w, "%s%s%s%s%s%s%v\n", pdv.ModuleId, delimiter, pd.Id, delimiter, pd.Unit, delimiter, pd.Value)

			}
		}

	} else {
		if outputTimestamp {
			fmt.Fprintf(w, "Timestamp:\t")
		}
		for _, pdv := range processDataValues {
			if outputTimestamp {
				fmt.Fprintf(w, "%s\n", time.Now().Format(time.RFC3339))
			}
			fmt.Fprintln(w, "Module:", pdv.ModuleId)
			for _, pd := range pdv.ProcessData {
				fmt.Fprintln(w, pd.Id, "\t", pd.Unit, "\t", pd.Value)
			}
			fmt.Fprintln(w)
		}

	}
}

func getModuleProcessdata(args []string) {
	var errOut io.Writer = os.Stderr
	var w io.Writer

	f, errFile := getOutFile()
	if errFile != nil {
		fmt.Fprintln(errOut, "Could not open file ", outputFile)
		return
	}
	if f != nil {
		w = f
		defer closeOutFile(f)
	} else {
		w = os.Stdout
	}
	// submitted values: moduleid pdid pdid2 pdid3...

	moduleId := args[0]

	lib := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Scheme:   authData.Scheme,
		Server:   authData.Server,
		Password: authData.Password,
	})

	_, err := lib.Login()
	if err != nil {
		fmt.Fprintln(errOut, "An error occurred:", err)

		return
	}
	defer lib.Logout()

	processDataValues, err := lib.ProcessDataModule(moduleId)
	if err != nil {
		fmt.Fprintln(errOut, "An error occurred:", err)
		return
	}

	if outputCSV {
		if !outputNoHeaders {
			if outputTimestamp {
				fmt.Fprintf(w, "Timestamp%s", delimiter)
			}
			fmt.Fprintf(w, "Module%sProcessdata Id%sProcessdata Unit%sProcessdata Value\n", delimiter, delimiter, delimiter)
		}
		for _, pdv := range processDataValues {
			for _, pd := range pdv.ProcessData {
				if outputTimestamp {
					fmt.Fprintf(w, "%s%s", time.Now().Format(time.RFC3339), delimiter)
				}
				fmt.Fprintf(w, "%s%s%s%s%s%s%v\n", pdv.ModuleId, delimiter, pd.Id, delimiter, pd.Unit, delimiter, pd.Value)

			}
		}

	} else {
		if outputTimestamp {
			fmt.Fprintf(w, "Timestamp:\t")
		}
		for _, pdv := range processDataValues {
			if outputTimestamp {
				fmt.Fprintf(w, "%s\n", time.Now().Format(time.RFC3339))
			}
			fmt.Fprintln(w, "Module:", pdv.ModuleId)
			for _, pd := range pdv.ProcessData {
				fmt.Fprintln(w, pd.Id, "\t", pd.Unit, "\t", pd.Value)
			}
			fmt.Fprintln(w)
		}

	}
}

// getOutFile returns a pointer to an opened file if the corresponding flags are set.
// If the return value is nil, output should be sent to os.Stdout
func getOutFile() (*os.File, error) {
	var f *os.File
	var err error
	if len(outputFile) > 0 {
		if outputAppend {
			f, err = os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				return nil, err
			}
		} else {
			f, err = os.Create(outputFile)
			if err != nil {
				return nil, err
			}
		}
		return f, nil
	}
	return nil, nil
}

// closeOutFile closes the file and handles errors
func closeOutFile(f *os.File) {
	err := f.Close()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

// Handle processdata-related commands
func handleProcessdata() {
	fmt.Println("\nUnknown or missing command.\nRun golrackpi processdata --help to show available commands.")
}
