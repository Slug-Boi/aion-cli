package cmd

import (
	"fmt"
	"os"

	"github.com/Slug-Boi/aion-cli/src/config"
	"github.com/spf13/cobra"
)

// solverCmd represents the api command
var solverCmd = &cobra.Command{
	Use:   "solver <default_solver>",
	Short: "This sub command edits the default solver in the config file.",
	Long: `This command allows you to edit the default solver in the config file.
	The current choices for solvers are min_cost and gurobi 
	The config file is located in the user's config directory. Example: ` + config.UserConf() + `config.json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckConfig()

		fmt.Println("Reading current config file")

		conf, err := config.GetConfigFile()
		if err != nil {
			fmt.Println(err)
		}

		conf.DefaultSolver = args[0]

		err = os.Truncate(config.UserConf()+"config.json", 0)
		if err != nil {
			fmt.Println(err)
		}

		f, err := os.OpenFile(config.UserConf()+"config.json", os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
		}

		config.WriteConfig(f, conf)
		fmt.Println("Default solver updated")

	},
}

func init() {
	configCmd.AddCommand(solverCmd)
}
