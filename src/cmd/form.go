package cmd

import (
	"fmt"
	"log"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/spf13/cobra"
)

// formCmd represents the form command
var formCmd = &cobra.Command{
	Use:   "form",
	Short: "A command to get back google form data as JSON",
	Long: `
	This command will retrieve a google form and return it as JSON. 
	You will need to provide the form ID as an argument.
	The command can be used to retrieve form responses incase you want to pipe it to another command or for testing purposes.
	`,
	Run: func(cmd *cobra.Command, args []string) {

		var conf forms.Config
		var err error

		if len(args) == 1 {
			// override formID from config file if formID is provided as an argument
			conf, err = forms.GetConfigFile()
			if err != nil {
				log.Fatal(err)
			}
			conf.FormID = args[0]
		} else {
			// get config file
			conf, err = forms.GetConfigFile()
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("Form is being processed with the following Form ID:", conf.FormID)

		form := forms.GetForm(conf)

		fmt.Println(form)
		
		},
}

func init() {
	rootCmd.AddCommand(formCmd)
}
