package cmd

import (
	"fmt"
	"os"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/spf13/cobra"
)

// icalSaveCmd represents the icalSave command
var icalSaveCmd = &cobra.Command{
	Use:   "icalSave",
	Short: "This sub command will toggle the saving of the ICal file when generating the html",
	Long: `This sub command will toggle the saving of the ICal file when generating the html.
	The ICal .ics file can be imported into a calendar application to show the time slots.
	By default, the ICal save value is set to false.
	The config file is located in the user's config directory. Example: ` + UserConf() + `config.json`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		CheckConfig()

		fmt.Println("Reading current config file")

		conf, err := forms.GetConfigFile()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Current ICal save value:", conf.Ical_save, "\nUpdating ICal save value to:", !conf.Ical_save)
		if conf.Ical_save {
			conf.Ical_save = false
		} else {
			conf.Ical_save = true
		}

		err = os.Truncate(UserConf()+"config.json", 0)
		if err != nil {
			fmt.Println(err)
		}

		f, err := os.OpenFile(UserConf()+"config.json", os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
		}

		WriteConfig(f, conf)
		fmt.Println("ICal save value updated")
	},
}

func init() {
	configCmd.AddCommand(icalSaveCmd)
}
