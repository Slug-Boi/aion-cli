package cmd

import (
	"github.com/Slug-Boi/aion-cli/src/config"
	"github.com/Slug-Boi/aion-cli/src/html"
	"github.com/spf13/cobra"
)

// generateCmd represents the generate command
var generateCmd = &cobra.Command{
	Use:   "generate <FilePath>",
	Short: "Generates an HTML file that is populated with CSV data",
	Long: `This command reads from a given CSV file, and then generates an HTML file populated with the CSV data.
The FilePath refers to the designated path. An example would be: 'C://Program/MyCSVFile.csv'`,
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckConfig()

		conf, err := config.GetConfigFile()
		if err != nil {
			Sugar.Panicf("Error getting config file: %v", err)
		}

		// If the cal flag is set, save the solution as an iCal file (or if config file is set to save as iCal)
		if val, _ := cmd.Flags().GetBool("cal"); val || conf.Ical_save {
			go html.CreateICal()
		}

		var gurobiFlag, minCostFlag bool
		gurobiFlag, _ = cmd.Flags().GetBool("gurobi")
		minCostFlag, _ = cmd.Flags().GetBool("minCost")
		if gurobiFlag && minCostFlag {
			Sugar.Panicf("Cannot use both Gurobi and minCost flags")
		}

		// If the gurobi flag is set, use Gurobi as the solver
		if gurobiFlag {
			conf.DefaultSolver = "gurobi"
		}
		// If the minCost flag is set, use min_cost flow graph solver
		if minCostFlag {
			conf.DefaultSolver = "minCost"
		}

		html.GenerateHTML(args, conf.DefaultSolver, conf.DefaultSorter)
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
	generateCmd.Flags().Bool("cal", false, "Save the solution as an iCal file")
	generateCmd.Flags().Bool("gurobi", false, "Use Gurobi as the solver")
	generateCmd.Flags().Bool("minCost", false, "Use min_cost flow graph solver")
}
