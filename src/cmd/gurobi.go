package cmd

import (
	"fmt"

	"github.com/Slug-Boi/aion-cli/src/config"
	"github.com/Slug-Boi/aion-cli/src/solvers/gurobi"
	"github.com/facette/natsort"
	"github.com/spf13/cobra"
)

// gurobiCmd represents the gurobi command
var gurobiCmd = &cobra.Command{
	Use:   "gurobi [formID]",
	Short: "The gurobi command will run a python solver program and will print the solution to the terminal",
	Long: `The gurobi solver command runs a python program that uses gurobipy to solve for minimum cost scheduling.
    The python program is embedded in the binary and uses the golang virtual file system to run the python program.
    It will print the solution to the terminal(this is mostly for debugging
                                               and testing purposes use the generate or root command to actually run the program).

    `,
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckConfig()

		if id, _ := cmd.Flags().GetBool("saveID"); id {
			config.CheckConfig()
			fmt.Println("\nSaving form ID to config file...")
			EditFormID(args[0])
			fmt.Println()
		}

		cost, Timeslots, wishLevels, groupTimeslotCost := gurobi.SolveGurobi(args)

		printSolutionGurobi(cost, Timeslots, wishLevels, groupTimeslotCost)
	},
}

func init() {
	solveCmd.AddCommand(gurobiCmd)
	gurobiCmd.Flags().Bool("saveID", false, "Save the formID to the config file")
}

func printSolutionGurobi(cost float64, Timeslots map[string]string, wishLevels map[string]string, groupTimeslotCost map[string]float64) {
	fmt.Println("Gurobi Solver")
	fmt.Println("Min Cost:", cost)
	fmt.Println("Timeslots:")

	// sort on keys to make sure the output is in order of group number
	keys := make([]string, 0, len(Timeslots))
	for k := range Timeslots {
		keys = append(keys, k)
	}
	natsort.Sort(keys)

	for _, group := range keys {
		fmt.Println(group, "->", Timeslots[group],
			"Wish Level:", wishLevels[group], "|", "Cost:", groupTimeslotCost[group+Timeslots[group]])
	}
}
