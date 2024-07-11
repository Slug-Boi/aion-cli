package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the configuration file",
	Long: `This command will remove the configuration file from the system.
	The configuration file is located at: ` + UserConf() + `config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Removing the configuration file")

		_,err := os.Stat(UserConf())
		if os.IsNotExist(err) {
			fmt.Println("Configuration file does not exist")
			os.Exit(0)
		}

		err = os.RemoveAll(UserConf())
		if err != nil {
			fmt.Println("Failed to remove the configuration file")
			fmt.Println(err)
		}
	},
}

func init() {
	configCmd.AddCommand(removeCmd)
}
