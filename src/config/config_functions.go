package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type Config struct {
	ConfVersion   string `json:"version"`
	DefaultSolver string `json:"default_solver"`
	FormID        string `json:"formID"`
	Ical_save     bool   `json:"ical_save"`
	CsvSave       bool   `json:"csv_save"`
	DefaultSorter string `json:"default_sorter"`
}

var version = "0.1"

// UserConf returns the user config directory
func UserConf() string {
	userConfig, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("Error getting user config directory exiting")
		os.Exit(1)
	}
	// Append aion-cli to the user config directory
	return userConfig + "/aion-cli/"
}

// CheckConfig checks if the config file exists
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

// CreateConfig creates a new config file
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

		response := ""

		version := "0.1"

		fmt.Println("Enter the Default Solver you would like to use [min_cost or gurobi]: (min_cost)")
		fmt.Scanln(&response)

		defSolver := response

		// If no response is given, default to min_cost
		if defSolver == "" || (defSolver != "min_cost" && defSolver != "gurobi") {
			defSolver = "min_cost"
		}

		fmt.Println("Turn on auto ical calendar file saving for the generate command? [y/n]: (n)")
		fmt.Scanln(&response)

		var ical bool

		if response == "" || (response != "y" && response != "n") || response == "n" {
			ical = false
		} else {
			ical = true
		}

		var csvSave bool
		fmt.Println("Turn on auto csv file saving for the generate, solve and form commands?\nThis will cache the form.csv file and the program will use that on all future runs of the program\n[y/n]: (n)")
		fmt.Scanln(&response)

		if response == "" || (response != "y" && response != "n") || response == "n" {
			csvSave = false
		} else {
			csvSave = true
		}

		var defaultSorter string
		fmt.Println("Enter the Default Sorter you would like to use when generating HTML [timeslot or group_number]: (timeslot)")
		fmt.Scanln(&response)
		if response == "" || (response != "timeslot" && response != "group_number") {
			defaultSorter = "timeslot"
		} else {
			defaultSorter = response
		}

		// Call writer to write to config file
		WriteConfig(f, Config{ConfVersion: version, DefaultSolver: defSolver, Ical_save: ical, CsvSave: csvSave, DefaultSorter: defaultSorter})

		os.Exit(0)

	} else {
		fmt.Println("Exiting")
		os.Exit(0)
	}
}

func WriteConfig(f *os.File, conf Config) {
	// Write to config file
	writer := bufio.NewWriter(f)

	writer.WriteString("{\n")

	// Write version to config file
	fmt.Println("Writing version to config file...")
	writer.WriteString(fmt.Sprintf("\t\"version\": \"%s\",\n", conf.ConfVersion))

	// Write API key for strawpoll to config file
	fmt.Println("Writing Default Solver to config file...")
	writer.WriteString(fmt.Sprintf("\t\"default_solver\": \"%s\",\n", conf.DefaultSolver))

	// Write form ID field to config file (empty for now)
	fmt.Println("Writing form ID field to config file...")
	writer.WriteString(fmt.Sprintf("\t\"formID\": \"%s\",\n", conf.FormID))

	// Write ical_save field to config file
	fmt.Println("Writing ical_save field to config file...")
	writer.WriteString(fmt.Sprintf("\t\"ical_save\": %t,\n", conf.Ical_save))

	// Write csvCaching
	fmt.Println("Writing csvSave field to config file...")
	writer.WriteString(fmt.Sprintf("\t\"csv_save\": %t,\n", conf.CsvSave))

	// Write default sorter
	fmt.Println("Writing default_sorter field to config file...")
	writer.WriteString(fmt.Sprintf("\t\"default_sorter\": \"%s\"", conf.DefaultSorter))

	writer.WriteString("\n}")

	writer.Flush()
}

// parse json with golang https://tutorialedge.net/golang/parsing-json-with-golang/
func GetConfigFile(testing ...string) (Config, error) {
	var jsonFile *os.File
	if len(testing) > 0 {
		// Open test config file location
		var err error
		jsonFile, err = os.Open(testing[0])
		if err != nil {
			return Config{}, fmt.Errorf("error opening test config file: %v", err)
		}
	} else {

		userConf, err := os.UserConfigDir()
		if err != nil {
			return Config{}, fmt.Errorf("error getting user config directory: %v", err)
		}

		// Open config file
		jsonFile, err = os.Open(userConf + "/aion-cli/config.json")
		if err != nil {

			return Config{}, fmt.Errorf("error opening config file: %v", err)
		}
	}

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, fmt.Errorf("error reading config file: %v", err)
	}

	// initialize config var
	var conf Config

	json.Unmarshal(byteValue, &conf)

	if conf.ConfVersion != version {
		fmt.Println("Config file version is out of date, please update your config file:")
		fmt.Println("Current version:", version)
		fmt.Println("Would you like the program to remove the old file and create a new one? (y/n)")
		var response string
		fmt.Scanln(&response)
		if response == "y" || response == "Y" {
			RemoveConfig()
			CreateConfig()
			fmt.Println("Exiting please run the program again")
			os.Exit(0)
		} else {
			fmt.Println("Exiting")
			os.Exit(0)
		}
		return Config{}, fmt.Errorf("Error with the config version")

	}

	return conf, nil
}

func RemoveConfig() {
	_, err := os.Stat(UserConf())
	if os.IsNotExist(err) {
		fmt.Println("Configuration file does not exist")
		os.Exit(0)
	}

	err = os.Remove(UserConf() + "config.json")
	if err != nil {
		fmt.Println("Failed to remove the configuration file")
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Remove(UserConf())
	if err != nil {
		fmt.Println("Failed to remove the configuration folder")
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Configuration file removed successfully")
}
