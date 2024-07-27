package cmd

import (
	libfuncs "github.com/Slug-Boi/aion-cli/lib_funcs"
	"github.com/Slug-Boi/aion-cli/solvers/graph"
	"github.com/spf13/cobra"
)

// solveCmd represents the solve command
var solveCmd = &cobra.Command{
	Use:   "solve [formID]",
	Short: "This command will run the selected (or the default solver) solver and print the solution to the terminal",
	Long: `This is mostly a debugging tool to see the output of the solver.
	The two available solvers are the min_cost flow graph (minCost) and gurobi (gurobi) solver.
	The solve command takes 0 or 1 argument. If no argument is provided, the form ID from the config file will be used.
	If an argument is provided, it will override the form ID from the config file.
	Example: aion-cli solve min_cost 1a2b3c4d5e6f7g8h9i0j
	The base solve command will use the default solver specified in the config file. (by default it is the min_cost solver)
	You can change the default solver in the config file by using the config command with the solver subcommand.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the config file (TODO: slight redundancy here getting the config
		// twice might add a bypass with params later)
		conf := libfuncs.SetupConfig(args)

		// Check which solver is the default
		if conf.DefaultSolver == "min_cost" {
			// Run the min_cost solver
			sink, users, cost, paths, nodeToTimeslot := graph.SolveMin_Cost(args)
			printSolutionMinCost(sink, users, cost, paths, nodeToTimeslot)
		} else {
			// Run the gurobi solver
			SolveGurobi(args)
		}

	},
}

func init() {
	rootCmd.AddCommand(solveCmd)
	solveCmd.Flags().Bool("save", false, "Save the solution as a CSV file")

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

// TODO: Figure out what format the output will be finalized in and save it to a file
func SaveSolution() {

}
