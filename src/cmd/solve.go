package cmd

import (
	"fmt"
	"log"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/Slug-Boi/aion-cli/graph"
	"github.com/spf13/cobra"
)

// solveCmd represents the solve command
var solveCmd = &cobra.Command{
	Use:   "solve",
	Short: "This command will run the selected solver and print the solution to the terminal",
	Long: `This is mostly a debugging tool to see the output of the solver.
	The two available solvers are the min_cost graph (min_cost) and gurobi (gurobi) solver.
	The solve command takes 0 or 1 argument. If no argument is provided, the form ID from the config file will be used.
	If an argument is provided, it will override the form ID from the config file.
	Example: aion-cli solve min_cost 1a2b3c4d5e6f7g8h9i0j
	`,
	Run: func(cmd *cobra.Command, args []string) {

		sink, users, cost, paths, nodeToTimeslot := Solve(args)

		println("Sink:", sink)
		println("Paths used")

		for j, path := range paths {
			println("Path:", j)
			i := sink
			println("Group Number:", users[path[path[i]]].GroupNumber)
			println("Timeslot:", nodeToTimeslot[path[i]])
			println(i)
			for i != 0 {
				println(path[i])
				i = path[i]
			}
		}

		//cost_new := ((cost - float64(groups)) - float64(len(paths)))

		println("Min cost:", int(cost), "≈", cost)

	},
}

func init() {
	rootCmd.AddCommand(solveCmd)
	solveCmd.Flags().Bool("save", false, "Save the solution as a CSV file")

}

func Solve(args []string) (int, map[int]forms.Form, float64, [][]int, map[int]string) {

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

	// Create a graph
	g, sink, users, nodeToTimeslot := graph.Translate(form)

	groups := len(form)

	cost, paths := graph.MinCostPath(len(g), groups, 0, sink, g)

	return sink, users, cost, paths, nodeToTimeslot
}

// This is a debugging graph
func debugGraphBuilder() []graph.Edge {
	// Edge values
	// From, To, Capacity, Cost
	g := []graph.Edge{}

	// 0 source
	// 8 sink
	// Groups -> 1, 2, 3
	// Timeslots -> 4, 5, 6, 7
	// Add edges to the graph
	g = append(g, graph.Edge{From: 0, To: 1, Capacity: 100, Cost: 1})
	g = append(g, graph.Edge{From: 0, To: 2, Capacity: 100, Cost: 1})
	g = append(g, graph.Edge{From: 0, To: 3, Capacity: 100, Cost: 1})

	g = append(g, graph.Edge{From: 1, To: 4, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 2, To: 4, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 2, To: 5, Capacity: 2, Cost: 2})
	g = append(g, graph.Edge{From: 3, To: 6, Capacity: 3, Cost: 3})
	g = append(g, graph.Edge{From: 3, To: 7, Capacity: 1, Cost: 1})

	g = append(g, graph.Edge{From: 4, To: 8, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 5, To: 8, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 6, To: 8, Capacity: 1, Cost: 1})
	g = append(g, graph.Edge{From: 7, To: 8, Capacity: 1, Cost: 1})

	return g
}

// TODO: Figure out what format the output will be finialized in and save it to a file
func SaveSolution() {

}
