package cmd

import (
	"fmt"
	"os"

	"github.com/Slug-Boi/aion-cli/config"
	"github.com/spf13/cobra"
)

// sorterCmd represents the sorter command
var sorterCmd = &cobra.Command{
	Use:   "sorter",
	Short: "This sub command will toggle which type of sorting to use when generating the html",
	Long: `This sub command will toggle which type of sorting is used when when generating the html.
	The default sorting is based on the time slots. The other type of sorting is based on the name of the groups.
	The group sorting is done using natural sorting. e.g. group1, group2, group 3, ... group10, group11, ...
	The config file is located in the user's config directory. Example: ` + config.UserConf() + `config.json`,
	Run: func(cmd *cobra.Command, args []string) {
		config.CheckConfig()

		fmt.Println("Reading current config file")

		conf, err := config.GetConfigFile()
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("Current sorting value:", conf.DefaultSorter)
		if conf.DefaultSorter == "timeslot" {
			conf.DefaultSorter = "group_number"
		} else {
			conf.DefaultSorter = "timeslot"
		}

		err = os.Truncate(config.UserConf()+"config.json", 0)
		if err != nil {
			fmt.Println(err)
		}

		f, err := os.OpenFile(config.UserConf()+"config.json", os.O_RDWR, 0644)
		if err != nil {
			fmt.Println(err)
		}

		config.WriteConfig(f, conf)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println("Sorting value updated")
	},
}

func init() {
	configCmd.AddCommand(sorterCmd)
}
