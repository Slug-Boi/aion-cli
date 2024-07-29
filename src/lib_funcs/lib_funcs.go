package libfuncs

import (
	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/Slug-Boi/aion-cli/logger"
)

var Sugar = logger.SetupLogger()

// This function will get the config file and setup the config struct
func SetupConfig(args []string, testing ...string) forms.Config {
	var conf forms.Config
	var err error

	if len(testing) > 0 {
		// override formID from config file if formID is provided as an argument
		if len(args) == 1 {
			conf, err = forms.GetConfigFile(testing[0])
			if err != nil {
				Sugar.Panicf("Error getting the config file using provided formID:\n", err.Error())
			}
			conf.FormID = args[0]
			return conf
		}
		// get config file
		conf, err = forms.GetConfigFile(testing[0])
		if err != nil {
			Sugar.Panicf("Error getting the config file:\n", err.Error())
		}
		return conf
	}

	// if formID is provided as an argument
	if len(args) == 1 {
		// override formID from config file if formID is provided as an argument
		conf, err = forms.GetConfigFile()
		if err != nil {
			Sugar.Panicf("Error getting the config file using provided formID:\n", err.Error())
		}
		conf.FormID = args[0]
		return conf
	}

	// get config file
	conf, err = forms.GetConfigFile()
	if err != nil {
		Sugar.Panicf("Error getting the config file:\n", err.Error())
	}
	return conf

}
