package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/spf13/cobra"
)

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure the config file",
	Long: `This command allows you to configure the config file.
	`,
	Run: func(cmd *cobra.Command, args []string) {
		CheckConfig()

		fmt.Println("Config file exists")
		fmt.Println("Please use sub commands to modify the config file. A list of these can be found by using -h or --help.")
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

}

func UserConf() string {
	userConfig, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config directory exiting")
		os.Exit(1)
	}
	// Append aion-cli to the user config directory
	return userConfig + "/aion-cli/"
}

func CheckConfig() {

	// Check if config file exists
	// If it does not exist, call creator function
	userConfig := UserConf()

	_, err := os.Stat(userConfig + "config.json")
	if os.IsNotExist(err) {
		fmt.Println("Config file does not exist")
		CreateConfig()
	}
}

func CreateConfig() {
	fmt.Println("Would you like to create a new config file? (y/n)")
	var response string
	fmt.Scanln(&response)
	if response == "y" || response == "Y" {

		//TODO: Probably move this to a seperate function
		fmt.Println("Creating config file")
		// Create config file
		userConfig := UserConf()

		// Create config directory
		err := os.MkdirAll(userConfig, 0755)
		if err != nil {
			fmt.Printf("Error creating config directory at: %s \nExiting", userConfig)
			os.Exit(1)
		}

		// Create and open config file
		f, err := os.OpenFile(userConfig+"config.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

		if err != nil {
			fmt.Printf("Error creating config file at: %s \nExiting", userConfig)
			os.Exit(1)
		}

		// Ask for user input for API key
		fmt.Println("Enter the API key for your strawpoll account: \n(it can be found here https://strawpoll.com/account/settings/ under the API section)")
		fmt.Scanln(&response)

		// Call writer to write to config file
		WriteConfig(f, forms.Config{Apikey: response})

		os.Exit(0)

	} else {
		fmt.Println("Exiting")
		os.Exit(0)
	}
}

func WriteConfig(f *os.File, conf forms.Config) {

	// Write to config file
	writer := bufio.NewWriter(f)

	writer.WriteString("{\n")

	// Write API key for strawpoll to config file
	fmt.Println("Writing API key to config file...")
	writer.WriteString(fmt.Sprintf("\t\"spAPI\": \"%s\",\n", conf.Apikey))

	// Write form ID field to config file (empty for now)
	fmt.Println("Writing form ID field to config file...")
	writer.WriteString(fmt.Sprintf("\t\"formID\": \"%s\"", conf.FormID))

	writer.WriteString("\n}")

	writer.Flush()

}
