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
	Args: cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		// TODO: Retrieve args after Strawpoll feat
		//specs := Reader.CsvToString(args[0])
		html.GenerateHTML()
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
