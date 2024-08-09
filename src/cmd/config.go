package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/Slug-Boi/aion-cli/src/config"
	"github.com/Slug-Boi/aion-cli/src/tui"
	"github.com/inancgumus/screen"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config <subcommand>",
	Short: "Configure the config file",
	Long: `This command allows you to configure the config file.
	The config file is used to store the default solver that the solve and generate command uses.
	You can also store the form ID for the google form you want to use by default when no arguments are given.
	The configuration file is located at: ` + config.UserConf() + `config.json
	`,
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckConfig()
		var choices = defChoices
		var terminal_msg = "Which config setting would you like to edit?"
		var subcmd = ""
		conf, err := config.GetConfigFile()
		if err != nil {
			Sugar.Panic("Error getting the config file:", err)
		}
		for {
			screen.Clear()
			var ans = ""
			if len(choices) > 0 {
				ans = tui.RunConfigTea(choices, terminal_msg)
			} else {
				ans = tui.StartTextTUI(terminal_msg, "formID")
			}
			switch ans {
			// Config select options
			case "Default Solver":
				choices = []string{"min_cost", "gurobi"}
				terminal_msg = "Which default solver would you like to use?\n(Gurobi requires third-party setup)"
			case "formID":
				// Text input bubble tea handler
				terminal_msg = "Enter the google sheets ID you would like to use by default\nThis will be used when no argument is given with commands"
				choices = []string{}
			case "ICal Save":
				choices = []string{"Enable", "Disable"}
				terminal_msg = "Would you like automatic ICal ics file saving when running the generate command?"
				subcmd = "ICal"
			case "CSV Save":
				choices = []string{"Enable", "Disable"}
				terminal_msg = "Would you like to enable saving and use of local CSV files?\n(This requires manual deletion of CSV files use only for testing or repeat runs on same form data)"
				subcmd = "CSV"
			case "Default Sorter":
				choices = []string{"timeslot", "group_number"}
				terminal_msg = "Which default sorting method would you like to use in the generated HTML?\n(Time slots - Earliest to Latest)\n(Group Number = Natural sorting - group 1, group 2, ..., group 10, group 11, ...)"
			case "Remove Config":
				config.RemoveConfig()
				os.Exit(0)
			case "Display Config":
				tui.RunDisplayTea(StringConfig(conf))

			// Answer option section
			case "Enable", "Disable":
				//TODO: Change this to be seperate function calls maybe?
				// Check which enable option
				if ans == "Enable" {
					if subcmd != "" && subcmd == "ICal" {
						conf.Ical_save = true
						config.StartConfigEdit(conf)
					} else if subcmd != "" {
						conf.CsvSave = true
						config.StartConfigEdit(conf)
					}
				} else {
					if subcmd != "" && subcmd == "ICal" {
						conf.Ical_save = false
						config.StartConfigEdit(conf)
					} else if subcmd != "" {
						conf.CsvSave = false
						config.StartConfigEdit(conf)
					}
				}
				// Reset Options
				choices, terminal_msg = ResetOptions()
			case "min_cost", "gurobi":
				conf.DefaultSolver = ans
				config.StartConfigEdit(conf)

				choices, terminal_msg = ResetOptions()
			case "timeslot", "group_number":
				conf.DefaultSorter = ans
				config.StartConfigEdit(conf)

				choices, terminal_msg = ResetOptions()

			case "quit":
				// Exit the program
				os.Exit(0)

			case "reset":
				choices, terminal_msg = ResetOptions()
			default:
				// formID update and error handling
				formIDSlice := strings.Split(ans, ":")
				if len(formIDSlice) > 1 && formIDSlice[0] == "formID" {
					conf.FormID = formIDSlice[1]
					config.StartConfigEdit(conf)

					choices, terminal_msg = ResetOptions()
				} else {
					Sugar.Panic("Something went wrong during bubbleTea command:", ans)
				}
			}
		}
	},
}

var defChoices = []string{"Default Solver", "formID", "ICal Save", "CSV Save", "Default Sorter", "Display Config", "Remove Config"}

func ResetOptions() ([]string, string) {
	// Reset options
	var choices = defChoices
	var terminal_msg = "Which config setting would you like to edit?"
	return choices, terminal_msg
}

func init() {
	rootCmd.AddCommand(configCmd)

}

func StringConfig(conf config.Config) string {
	// String build current config
	sb := strings.Builder{}

	sb.WriteString("## Config file location:\n\t" + config.UserConf() + "config.json\n")
	sb.WriteString("# Current config file:\n")
	sb.WriteString("## Default solver:\t" + conf.DefaultSolver + "\n")
	if conf.FormID == "" {
		sb.WriteString("## Form ID:\t *\\*Not set*\\*\n")
	} else {
		sb.WriteString("## Form ID:\t" + conf.FormID + "\n")
	}
	sb.WriteString("## ICal save:\t" + fmt.Sprintf("%t", conf.Ical_save) + "\n")
	sb.WriteString("## Csv save:\t" + fmt.Sprintf("%t", conf.CsvSave) + "\n")
	sb.WriteString("## Default sorter:\t" + conf.DefaultSorter + "\n")

	return sb.String()
}
