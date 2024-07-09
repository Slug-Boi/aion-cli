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

	err := os.WriteFile("../config.json", []byte(`{"spAPI": "test_api_key"}`), 0644)
	if err != nil {
		t.Error("Failed to create or write to config.json file",err)
	}

	conf, err := forms.GetConfigFile("../config.json")
	if err != nil {
		t.Error(err)
	}

	if conf.Apikey != "test_api_key" {
		t.Error("Failed to read config.json file \n Expected: test_api_key \n Got:", conf.Apikey)
	}

	// Expand test as conf struct grows
}
