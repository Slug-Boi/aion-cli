package cmd

import (
	"fmt"

	"github.com/Slug-Boi/aion-cli/config"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config <subcommand>",
	Short: "Configure the config file",
	Long: `This command allows you to configure the config file.
	The config file is used to store the default solver that the solve and generate command uses.
	You can also store the form ID for the google form you want to use by default when no arguments are given.
	The configuration file is located at: ` + config.UserConf() + `config.json
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckConfig()

		fmt.Println("Config file exists")
		fmt.Println("Please use sub commands to modify the config file.\nA list of these can be found by using -h or --help.")

		conf, err := config.GetConfigFile()
		if err != nil {
			Sugar.Panicf("Error reading config file: %v", err)
		}

		fmt.Println()
		PrintConfig(conf)
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

}

func PrintConfig(conf config.Config) {
	fmt.Println("Config file location:", config.UserConf()+"config.json")
	fmt.Println("Current config file:")
	fmt.Println("\tDefault solver:", conf.DefaultSolver)
	fmt.Println("\tForm ID:", conf.FormID)
	fmt.Println("\tICal save:", conf.Ical_save)
	fmt.Println("\tCsv save:", conf.CsvSave)
	fmt.Println("\tDefault sorter:", conf.DefaultSorter)
}
