package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

// This is an example of how to setup the logger for any CMD command you can then use it when doing calls.
// A similar logger can be setup in any other file that requires it by importing CMD and calling the SetupLogger function
var sugar = SetupLogger()

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "aion", 
	Short: "A scheduling tool that takes in wishes and outputs a schedule",
//TODO: CHANGE TO AION ACSII ART LATER
	Long: `                                                     
 _____ _                 _ _ 
|  _  |_|___ ___ ___ ___| |_|
|     | | . |   |___|  _| | |
|__|__|_|___|_|_|   |___|_|_|
                               
                        
This is a CLI tool scheduling tool that takes in timeslot wishes and outputs a schedule. 
The tool was designed around google forms format and should therefore work with any format that is similar to the google forms format.
                                              `,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	
}

// The setup logger function will live in the root command as most logging should be propagated up to the CMD commands
// This allows us to create a local logger for each command that can be used to log errors and info messages
func SetupLogger() *zap.SugaredLogger {
	// setup suggered zap logger
	// https://github.com/uber-go/zap
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()

	return sugar
}
