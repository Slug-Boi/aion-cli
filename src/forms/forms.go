package forms

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

type Form struct {
	Participant_count int `json:"participant_count"`
	PollOptions       []struct {
		Id         string `json:"id"`
		Start_time int64  `json:"start_time"`
		End_time   int64  `json:"end_time"`
	} `json:"poll_options"`
	PollResults []struct {
		Name  string `json:"name"`
		Id    string `json:"id"`
		Votes []int  `json:"poll_votes"`
	} `json:"poll_participants"`
}

type Config struct {
	Apikey string `json:"spAPI"`
	FormID string `json:"formID"`
}

// Generated by curl-to-Go: https://mholt.github.io/curl-to-go

//	curl --request GET \
//	  --url https://api.strawpoll.com/v3/polls/XmZRQjmaPgd/results \
//	  --header 'Accept: application/json' \
//	  --header 'X-API-Key: a55c3866-3ac0-11ef-9ad2-345670717338'
//
// TODO: Add optional flag for api key currently uses config file
func GetForm(conf Config) Form {
	// create http get request to strawpoll api with form id
	url := fmt.Sprintf("https://api.strawpoll.com/v3/polls/%s/results", conf.FormID)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}

	// set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Api-Key", conf.Apikey)

	//TODO: Potentially change this to a return error when changing to zap logger
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// turn resp header into a byte array
	byteValue, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// turn byte array into struct via json unmarshal
	var form Form
	json.Unmarshal(byteValue, &form)

	return form
}

// parse json with golang https://tutorialedge.net/golang/parsing-json-with-golang/
func GetConfigFile() (Config, error) {

	// Open config file
	jsonFile, err := os.Open("../config.json")
	if err != nil {
		return Config{}, fmt.Errorf("Error opening config file: %v", err)
	}

	// read our opened jsonFile as a byte array.
	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, fmt.Errorf("Error reading config file: %v", err)
	}

	// initialize config var
	var conf Config

	json.Unmarshal(byteValue, &conf)

	return conf, nil
}
