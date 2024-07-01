package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/Slug-Boi/gopher_scheduler/forms" 
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
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		//TODO: Change to something more meaningful
		fmt.Println("form called")

		formService, err :=forms.CreateService()

		if err != nil {
			//TODO: Change to zap logger later
			fmt.Println(err)
		}

		//TODO: This could maybe be handled as a config file as well to minimize redundancy for repeated calls
		formID := args[0]

		form, err := forms.GetForm(formService, formID)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println(string(form))
		
		},
}

func init() {
	rootCmd.AddCommand(formCmd)
}
