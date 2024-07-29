package html_test

import (
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/Slug-Boi/aion-cli/html"
)

func cleanup() {
	os.Remove("calendar.ics")
}

// These tests should be run sequentially
// https://stackoverflow.com/questions/31201858/how-to-run-golang-tests-sequentially

// Test requires a public google sheet if this link dies the test will fail and need to be updated
// Sheets ID: 15_RVajfepi7MxTZ_ZWFgz8PCE1axb3NcAyamF8mGl3E
func TestGenerateHTML(t *testing.T) {

	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}
	
	resChan := make(chan *http.Response)

	go html.GenerateHTML([]string{"15_RVajfepi7MxTZ_ZWFgz8PCE1axb3NcAyamF8mGl3E"}, "min_cost")

	go func(resChan chan *http.Response) {
		time.Sleep(1 * time.Second)
		res, err := http.Get("http://localhost:80/")
		if err != nil {
			t.Errorf("Error getting the response: %v", err)
		}
		resChan <- res
	}(resChan)
	select {
	case <-time.After(5 * time.Second):
		t.Errorf("Timed out waiting for response")
		return
	case res := <-resChan:
		if res.StatusCode != 200 {
			t.Errorf("Expected status code 200, got %v", res.StatusCode)
			return
		}
	}

}

func TestICal(t *testing.T) {

	if os.Getenv("CI") == "true" {
		t.Skip("Skipping test in CI environment")
	}

	defer cleanup()

	html.CreateICal()

	// Check if the file exists
	_, err := os.Stat("calendar.ics")
	if err != nil {
		t.Errorf("Error getting the file: %v", err)
	}
}
