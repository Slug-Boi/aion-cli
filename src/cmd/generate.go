package cmd

import (
	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/Slug-Boi/aion-cli/html"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate <FilePath>",
	Short: "Generates an HTML file that is populated with CSV data",
	Long: `This command reads from a given CSV file, and then generates an HTML file populated with the CSV data.
The FilePath refers to the designated path. An example would be: 'C://Program/MyCSVFile.csv'`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckConfig()

		conf, err := forms.GetConfigFile()	
		if err != nil {
			Sugar.Panicf("Error getting config file: %v", err)
		}

		// If the cal flag is set, save the solution as an iCal file (or if config file is set to save as iCal)
		if val, _ := cmd.Flags().GetBool("cal"); val || conf.Ical_save {
			go html.CreateICal()
		}

		html.GenerateHTML(args, conf.DefaultSolver)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().Bool("cal", false, "Save the solution as an iCal file")
}
