/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/Slug-Boi/aion-cli/graph"
	"github.com/spf13/cobra"
)

// minCostCmd represents the minCost command
var minCostCmd = &cobra.Command{
	Use:   "minCost [formID]",
	Short: "This command will run the minCost graph solver and print the solution to the terminal",
	Long: `The min_cost command will run the min cost flow graph solver to solve for minimum cost scheduling.
	The solver is based on SPFA (Shortest Path Faster Algorithm) uses negative cycles to reduce cost.
	It will print the solution to the terminal (this is mostly for debugging 
	and testing purposes use the generate or root command to actually run the program).
	
	`,
	Run: func(cmd *cobra.Command, args []string) {

	sink, users, cost, paths, nodeToTimeslot := SolveMin_Cost(args)

	printSolutionMinCost(sink, users, cost, paths, nodeToTimeslot)
	},
}

func init() {
	solveCmd.AddCommand(minCostCmd)

}

func SolveMin_Cost(args []string) (int, map[int]forms.Form, float64, [][]int, map[int]string) {

	// Get the config file
	conf := SetupConfig(args)

	//TODO: Make this a hidden form id see if there is a way to make it display when clicked
	fmt.Println("Form is being processed with the following Form ID:", conf.FormID)

	form := forms.GetForm(conf)

	// Create a graph
	g, sink, users, nodeToTimeslot := graph.Translate(form)

	groups := len(form)

	cost, paths := graph.MinCostPath(len(g), groups, 0, sink, g)


	return sink, users, cost, paths, nodeToTimeslot
}

func printSolutionMinCost(sink int, users map[int]forms.Form, cost float64, paths [][]int, nodeToTimeslot map[int]string) {

		fmt.Println("Sink:", sink)
		fmt.Println("Paths used:")

		for j, path := range paths {
			fmt.Println("Path:", j)
			i := sink
			fmt.Println("Group Number:", users[path[path[i]]].GroupNumber)
			fmt.Println("Timeslot:", nodeToTimeslot[path[i]])
			fmt.Println(i)
			for i != 0 {
				fmt.Println(path[i])
				i = path[i]
			}
		}

		println("Min cost:", int(cost), "≈", cost)
}