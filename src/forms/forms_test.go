package forms_test

import (
	"os"
	"testing"

	"github.com/Slug-Boi/aion-cli/config"
	"github.com/Slug-Boi/aion-cli/forms"
)

var data = []byte(`Timestamp,Group Number,Lottery String,12-04-24 Monday [8:00-10:00],12-04-24 Monday [10:00-12:00],15-04-24 Thursday [10:00-12:00],15-04-24 Thursday [14:30-16:30]
26/07/2024 10:50:58,Group 1,JavaBois,Want,Can do,,
26/07/2024 10:51:21,Group 2,PartyInTheSewers,Want,Want,,`)

func cleanup() {
	os.Remove("form.csv")
	os.Remove("config.json")
}

var configData = []byte(`
{
	"version": "0.1",
	"default_solver": "min_cost",
	"FormID": "",
	"ical_save": false,
	"csv_save": false,
	"default_sorter": "group_number"
}`)

func TestGetConfig(t *testing.T) {

	defer cleanup()

	err := os.WriteFile("config.json", configData, 0644)
	if err != nil {
		t.Error("Failed to create or write to config.json file", err)
	}

	conf, err := config.GetConfigFile("config.json")
	if err != nil {
		t.Error(err)
	}

	if conf.DefaultSolver != "min_cost" {
		t.Error("Failed to read config.json file \n Expected: min_cost \n Got:", conf.DefaultSolver)
	}

	if conf.FormID != "" {
		t.Error("Failed to read config.json file \n Expected: '' \n Got:", conf.FormID)
	}

	if conf.Ical_save != false {
		t.Error("Failed to read config.json file \n Expected: false \n Got:", conf.Ical_save)
	}

	// Expand test as conf struct grows
}

func TestGetForm(t *testing.T) {
	err := os.WriteFile("form.csv", data, 0644)
	if err != nil {
		t.Error("Failed to create or write to form.csv file", err)
	}

	defer cleanup()

	var conf config.Config

	form := forms.GetForm(conf, true)

	if len(form) == 0 {
		t.Error("Failed to read form")
	}
}
