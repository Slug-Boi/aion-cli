package cmd

import (
	"fmt"
	"os"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/spf13/cobra"
)

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api <API_Key>",
	Short: "This sub command edits the API key in the config file.",
	Long: `This command allows you to edit the API key in the config file. 
	The config file is located in the user's config directory. Example: ` + UserConf() + `config.json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		CheckConfig()

		fmt.Println("Reading current config file")

		conf, err := forms.GetConfigFile()
		if err != nil {
			fmt.Println(err)
		}

		conf.Apikey = args[0]

		err = os.Truncate(UserConf()+"config.json", 0)
		if err != nil {
			fmt.Println(err)
		}

		f, err := os.OpenFile(UserConf()+"config.json", os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
		}

		WriteConfig(f, conf)
		fmt.Println("API key updated")

	},
}

func init() {
	configCmd.AddCommand(apiCmd)
}
