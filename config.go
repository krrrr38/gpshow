package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/krrrr38/gpshow/utils"
)

// Configuration for picture-show
type Configuration struct {
	Title    string   `json:"title"`
	Sections []string `json:"sections"`
}

// Config parse json file to generate picture-show configuration
func Config(filepath string) (config Configuration) {
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		utils.Log("error", fmt.Sprintf("conf.js not found under @ `%s`. Please try `gpshow init <project_name>`", filepath))
		os.Exit(1)
	}

	err = json.Unmarshal(bytes, &config)
	utils.DieIf(err)
	return config
}
