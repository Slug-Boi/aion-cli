/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Slug-Boi/aion-cli/forms"
	libfuncs "github.com/Slug-Boi/aion-cli/lib_funcs"
	"github.com/Slug-Boi/aion-cli/solvers/gurobi"
	"github.com/spf13/cobra"
)

// gurobiCmd represents the gurobi command
var gurobiCmd = &cobra.Command{
	Use:   "gurobi [formID]",
	Short: "The gurobi command will run a python solver program and will print the solution to the terminal",
	Long: `The gurobi solver command runs a python program that uses gurobipy to solve for minimum cost scheduling.
	The python program is embedded in the binary and uses the golang virtual file system to run the python program.
	It will print the solution to the terminal (this is mostly for debugging 
	and testing purposes use the generate or root command to actually run the program).
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckConfig()

		if id, _ := cmd.Flags().GetBool("saveID"); id {
			CheckConfig()
			fmt.Println("\nSaving form ID to config file...")
			EditFormID(args[0])
			fmt.Println()
		}

		SolveGurobi(args)
	},
}

func init() {
	solveCmd.AddCommand(gurobiCmd)
	gurobiCmd.Flags().Bool("saveID", false, "Save the formID to the config file")
}

func SolveGurobi(args []string) {

	// Get the config file
	conf := libfuncs.SetupConfig(args)

	Sugar.Debugln("Form is being processed with the following Form ID:", conf.FormID)

	// Get the form data
	data := forms.GetForm(conf)

	// Run the gurobi python program
	out, err := gurobi.RunGurobi(data)
	if err != nil {
		Sugar.Panicf("Error running gurobi: %v", err)
	}

	//TODO: add return values for the gurobi solver
	fmt.Println(out)
}
