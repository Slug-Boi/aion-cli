package libfuncs_test

import (
	"os"
	"testing"

	libfuncs "github.com/Slug-Boi/aion-cli/lib_funcs"
)

func cleanup() {
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

func TestSetupConfig(t *testing.T) {
	defer cleanup()
	
	err := os.WriteFile("config.json", configData, 0644)
	if err != nil {
		t.Error("Failed to create or write to config.json file", err)
	}

	form := libfuncs.SetupConfig([]string{""}, "config.json")

	if form.DefaultSolver != "min_cost" {
		t.Error("Failed to read config.json file \n Expected: min_cost \n Got:", form.DefaultSolver)
	}

	if form.FormID != "" {
		t.Error("Failed to read config.json file \n Expected: '' \n Got:", form.FormID)
	}

	if form.Ical_save != false {
		t.Error("Failed to read config.json file \n Expected: false \n Got:", form.Ical_save)
	}

	if form.CsvSave != false {
		t.Error("Failed to read config.json file \n Expected: false \n Got:", form.CsvSave)
	}

}



func TestSetupConfigWithFormID(t *testing.T) {
	defer cleanup()

	err := os.WriteFile("config.json", configData, 0644)
	if err != nil {
		t.Error("Failed to create or write to config.json file", err)
	}

	form := libfuncs.SetupConfig([]string{"test"}, "config.json")

	if form.DefaultSolver != "min_cost" {
		t.Error("Failed to read config.json file \n Expected: min_cost \n Got:", form.DefaultSolver)
	}

	if form.FormID != "test" {
		t.Error("Failed to read config.json file \n Expected: '' \n Got:", form.FormID)
	}

	if form.Ical_save != false {
		t.Error("Failed to read config.json file \n Expected: false \n Got:", form.Ical_save)
	}

}
