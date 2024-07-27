package forms_test

import (
	"os"
	"testing"

	"github.com/Slug-Boi/aion-cli/forms"
)

func cleanup() {
	os.Remove("../config.json")
}

func TestGetConfig(t *testing.T) {

	defer cleanup()

	t.Log("Testing GetConfig")

	err := os.WriteFile("../config.json", []byte(`{"DefaultSolver": "min_cost"}`), 0644)
	if err != nil {
		t.Error("Failed to create or write to config.json file",err)
	}

	conf, err := forms.GetConfigFile("../config.json")
	if err != nil {
		t.Error(err)
	}

	if conf.DefaultSolver != "min_cost" {
		t.Error("Failed to read config.json file \n Expected: min_cost \n Got:", conf.DefaultSolver)
	}

	// Expand test as conf struct grows
}
