package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golrackpi",
	Short: "Go Client for Kostal Plenticore Inverters",
	Long: `
 golrackpi is a small CLI application to read some values from Kostal Plenticore Inverters.
 `,
}

func Exec() {

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
