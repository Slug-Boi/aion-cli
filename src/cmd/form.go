package cmd

import (
	"fmt"
	"log"

	"github.com/Slug-Boi/aion-cli/src/config"
	"github.com/Slug-Boi/aion-cli/src/forms"
	"github.com/spf13/cobra"
)

// formCmd represents the form command
var formCmd = &cobra.Command{
	Use:   "form [formID]",
	Short: "A command to get back google form data as JSON",
	Long: `
	This command will retrieve a google form and return it as JSON. 
	You will need to provide the form ID as an argument.
	The command can be used to retrieve form responses incase you want to pipe it to another command or for testing purposes.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckConfig()

		var conf config.Config
		var err error

		if len(args) == 1 {
			// override formID from config file if formID is provided as an argument
			conf, err = config.GetConfigFile()
			if err != nil {
				log.Fatal(err)
			}
			conf.FormID = args[0]

			// If save flag is provided, save the formID to the config file
			if id, _ := cmd.Flags().GetBool("save"); id {
				config.CheckConfig()
				fmt.Println("\nSaving form ID to config file...")
				EditFormID(args[0])
				fmt.Println()
			}

		} else {
			// get config file
			conf, err = config.GetConfigFile()
			if err != nil {
				log.Fatal(err)
			}
		}

		Sugar.Debugln("Form is being processed with the following Form ID:", conf.FormID)

		form := forms.GetForm(conf)

		fmt.Println(form)

	},
}

func init() {
	rootCmd.AddCommand(formCmd)
	formCmd.Flags().Bool("save", false, "Save the form ID to the config file")

}
