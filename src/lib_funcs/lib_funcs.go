package libfuncs

import (
	"github.com/Slug-Boi/aion-cli/forms"
	"github.com/Slug-Boi/aion-cli/logger"
)

var Sugar = logger.SetupLogger()

// This function will get the config file and setup the config struct
func SetupConfig(args []string) forms.Config {
	var conf forms.Config
	var err error

	println(len(args))

	if len(args) == 1 {
		// override formID from config file if formID is provided as an argument
		conf, err = forms.GetConfigFile()
		if err != nil {
			Sugar.Panicf("Error getting the config file using provided formID:\n", err.Error())
		}
		conf.FormID = args[0]
	} else {
		// get config file
		conf, err = forms.GetConfigFile()
		if err != nil {
			Sugar.Panicf("Error getting the config file:\n", err.Error())
		}
	}

	return conf
}
