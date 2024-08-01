package cmd

import (
	"fmt"
	"os"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/spf13/cobra"
)

// csvSaveCmd represents the csvSave command
var csvSaveCmd = &cobra.Command{
	Use:   "csvSave",
	Short: "This sub command edits the csv save value in the config file.",
	Long: `This command allows you to edit the default solver in the config file.
	The current choices for solvers are min_cost and gurobi 
	The config file is located in the user's config directory. Example: ` + UserConf() + `config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckConfig()

		fmt.Println("Reading current config file")

		conf, err := forms.GetConfigFile()
		if err != nil {
			fmt.Println(err)
		}

		conf.CsvSave = true

		err = os.Truncate(UserConf()+"config.json", 0)
		if err != nil {
			fmt.Println(err)
		}

		f, err := os.OpenFile(UserConf()+"config.json", os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
		}

		WriteConfig(f, conf)
		fmt.Println("Csv save updated")
	},
}

func init() {
	configCmd.AddCommand(csvSaveCmd)
}
