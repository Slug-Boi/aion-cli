package cmd

import (
	"strings"

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
		// TODO: Retrieve args after Strawpoll feat
		//specs := Reader.CsvToString(args[0])

		var webData []html.WebData

		sink, users, _, paths, nodeToTimeslot := SolveMin_Cost(args)

		// TODO: Convert path slice to an actual linear path to reduce nested calls as it looks disgusting
		// Might be as simple as using the uncommented code below try later
		for _, path := range paths {
			// i := sink
			// pathData := []int{sink}
			// for i != 0 {
			// 	// split the group number from the name
			// 	i = path[i]
			// 	pathData = append(pathData, i)
			// }
			timeslot := nodeToTimeslot[path[sink]]
			Date_Day_Timeslot := strings.Split(timeslot, " ")
			trimmedTimeslot := strings.Trim(Date_Day_Timeslot[2], "[]")
			day := Date_Day_Timeslot[1]
			date := Date_Day_Timeslot[0]
			wishLevel := users[path[path[sink]]].Votes[timeslot]
			
			webData = append(webData, html.WebData{GroupNumber: users[path[path[sink]]].GroupNumber, Timeslot: trimmedTimeslot, Day: day,Date: date, WishLevel: wishLevel})
		}

		html.GenerateHTML(webData)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
