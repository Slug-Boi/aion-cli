package cmd

import (
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
		html.GenerateHTML(args)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
