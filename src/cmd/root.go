/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gopher_sched", //TODO: change this to the name of the tool (if we come up with something better)
	Short: "A scheduling tool that takes in wishes and outputs a schedule",
	Long: `                                   
             _                     
 ___ ___ ___| |_ ___ ___           
| . | . | . |   | -_|  _|          
|_  |___|  _|_|_|___|_|            
|___|   |_|                        
                                   
         _         _     _         
 ___ ___| |_ ___ _| |_ _| |___ ___ 
|_ -|  _|   | -_| . | | | | -_|  _|
|___|___|_|_|___|___|___|_|___|_|  
                                      
                           
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
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gopher_sched.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
