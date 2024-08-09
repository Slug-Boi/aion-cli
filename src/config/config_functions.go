package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/Slug-Boi/aion-cli/tui"
	"github.com/inancgumus/screen"
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

		// Default solver
		screen.Clear()
		choices := []string{"min_cost", "gurobi"}
		terminal_msg := "Which default solver would you like to use?\n(Gurobi requires third-party setup)"
		ans := tui.RunConfigTea(choices, terminal_msg)

		defSolver := ans

		// If no response is given, default to min_cost
		if defSolver == "" || (defSolver != "min_cost" && defSolver != "gurobi") {
			defSolver = "min_cost"
		}
		// formID
		screen.Clear()
		terminal_msg = "Enter a google sheets ID you would like to use by default\nThis will be used when no argument is given with commands\nIt can be changed later using the config command or flags\nExit to leave blank"
		ans = tui.StartTextTUI(terminal_msg, "formID")

		if ans == "reset" {
			ans = ""
		}

		formID := ans

		// ICal Save
		screen.Clear()
		choices = []string{"Enable", "Disable"}
		terminal_msg = "Would you like automatic ICal ics file saving when running the generate command?"
		ans = tui.RunConfigTea(choices, terminal_msg)

		var ical bool

		if ans == "Enable" {
			ical = true
		} else {
			ical = false
		}

		// CSV Save
		screen.Clear()
		choices = []string{"Enable", "Disable"}
		terminal_msg = "Would you like to enable saving and use of local CSV files?\n(This requires manual deletion of CSV files use only for testing or repeat runs on same form data)"
		ans = tui.RunConfigTea(choices, terminal_msg)

		var csvSave bool

		if ans == "Enable" {
			csvSave = true
		} else {
			csvSave = false
		}

		// Default Sorter
		screen.Clear()
		choices = []string{"timeslot", "group_number"}
		terminal_msg = "Which default sorting method would you like to use in the generated HTML?\n(Time slots - Earliest to Latest)\n(Group Number = Natural sorting - group 1, group 2, ..., group 10, group 11, ...)"
		ans = tui.RunConfigTea(choices,terminal_msg)

		var defaultSorter = ans

		// Call writer to write to config file
		WriteConfig(f, Config{ConfVersion: version, DefaultSolver: defSolver, FormID: formID, Ical_save: ical, CsvSave: csvSave, DefaultSorter: defaultSorter})

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

func StartConfigEdit(conf Config) {
	CheckConfig()

	screen.Clear()

	fmt.Println("Reading current config file")

	err := os.Truncate(UserConf()+"config.json", 0)
	if err != nil {
		fmt.Println(err)
	}

	f, err := os.OpenFile(UserConf()+"config.json", os.O_RDWR, 0644)
	if err != nil {
		fmt.Println(err)
	}

	WriteConfig(f, conf)
	time.Sleep(1 * time.Second)
	screen.Clear()
}
