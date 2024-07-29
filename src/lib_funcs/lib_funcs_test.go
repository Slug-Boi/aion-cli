package libfuncs_test

import (
	"os"
	"testing"

	libfuncs "github.com/Slug-Boi/aion-cli/lib_funcs"
)

func cleanup() {
	os.Remove("config.json")
}

func TestSetupConfig(t *testing.T) {
	defer cleanup()
	
	err := os.WriteFile("config.json", []byte(`
	{
		"DefaultSolver": "min_cost",
		"FormID": "",
		"ical_save": false
	}`), 0644)
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

}

func TestSetupConfigWithFormID(t *testing.T) {
	defer cleanup()

	err := os.WriteFile("config.json", []byte(`
	{
		"DefaultSolver": "min_cost",
		"FormID": "",
		"ical_save": false
	}`), 0644)
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
