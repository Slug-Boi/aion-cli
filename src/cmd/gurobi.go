/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/Slug-Boi/aion-cli/gurobi"
	"github.com/spf13/cobra"
)

// gurobiCmd represents the gurobi command
var gurobiCmd = &cobra.Command{
	Use:   "gurobi [formID]",
	Short: "The gurobi command will run a python program that will solve the scheduling problem",
	Long: `The gurobi solver command runs a python program that uses gurobipy to solve for minimum cost scheduling.
	The python program is embedded in the binary and uses the golang virtual file system to run the python program.
	
	`,
	Run: func(cmd *cobra.Command, args []string) {
		// Get the config file
		conf := SetupConfig(args)

		fmt.Println("Form is being processed with the following Form ID:", conf.FormID)

		// Get the form data
		data := forms.GetForm(conf)

		// Translate the form data to gurobi syntax
		groups, timeslots := gurobi.TranslateGurobi(data)

		// Run the gurobi python program
		out, err := gurobi.RunGurobi(groups, timeslots)
		if err != nil {
			Sugar.Panicf("Error running gurobi: %v", err)
		}

		fmt.Println(out)
	},
}

func init() {
	solveCmd.AddCommand(gurobiCmd)
	
}

