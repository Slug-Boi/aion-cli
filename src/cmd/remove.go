package cmd

import (
	"fmt"

	"github.com/Slug-Boi/aion-cli/src/config"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Removes the configuration file",
	Long: `This command will remove the configuration file from the system.
	The configuration file is located at: ` + config.UserConf() + `config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Removing the configuration file")
		config.RemoveConfig()
	},
}

func init() {
	configCmd.AddCommand(removeCmd)
}
