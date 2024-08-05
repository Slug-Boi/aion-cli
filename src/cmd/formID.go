package cmd

import (
	"fmt"
	"os"

	"github.com/Slug-Boi/aion-cli/config"
	"github.com/spf13/cobra"
)

// formIDCmd represents the formID command
var formIDCmd = &cobra.Command{
	Use:   "formID <formID>",
	Short: "This command will edit the form ID in the config file.",
	Long: `This command allows you to edit the form ID in the config file.
	The formID is the ID of the strawpoll form that you want to get data from.
	This will be used by default if you call the form command without any arguments.
	Calling the form command with the --save flag will save the formID to the config file as well.
	The config file is located in the user's config directory. Example: ` + config.UserConf() + `config.json`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		EditFormID(args[0])
	},
}

func init() {
	configCmd.AddCommand(formIDCmd)
}

func EditFormID(id string) {
	config.CheckConfig()

	fmt.Println("Reading current config file")

	conf, err := config.GetConfigFile()
	if err != nil {
		fmt.Println(err)
	}

	conf.FormID = id

	// Truncate to delete everything in the file
	err = os.Truncate(config.UserConf()+"config.json", 0)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.OpenFile(config.UserConf()+"config.json", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}

	// Write the new config to the file
	config.WriteConfig(f, conf)
	fmt.Println("Form ID updated")
}
