package cmd

import (
	"fmt"
	"os"

	"github.com/Slug-Boi/aion-cli/src/config"
	"github.com/spf13/cobra"
)

// csvSaveCmd represents the csvSave command
var csvSaveCmd = &cobra.Command{
	Use:   "csvSave",
	Short: "This sub command toggles the csv save value in the config file.",
	Long: `This command allows you to toggle the csv save value in the config file.
	This value will determine if the program will save the output to a csv file.
	The csv file will be used for the solver when the program is run. Manual deletion of this file
	is required if you want to run a the program on a different form.
	By default, the csv save value is set to false.
	The config file is located in the user's config directory. Example: ` + config.UserConf() + `config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckConfig()

		fmt.Println("Reading current config file")

		conf, err := config.GetConfigFile()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Current csv save value:", conf.CsvSave, "\nUpdating csv save value to:", !conf.CsvSave)
		if conf.CsvSave {
			conf.CsvSave = false
		} else {
			conf.CsvSave = true
		}

		err = os.Truncate(config.UserConf()+"config.json", 0)
		if err != nil {
			fmt.Println(err)
		}

		f, err := os.OpenFile(config.UserConf()+"config.json", os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
		}

		config.WriteConfig(f, conf)
		fmt.Println("Csv save updated")
	},
}

func init() {
	configCmd.AddCommand(csvSaveCmd)
}
