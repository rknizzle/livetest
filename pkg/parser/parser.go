// Package parser parses the config file and creates jobs, connects to database and sets up notifications

package parser

import (
	"encoding/json"
	"github.com/rknizzle/livetest/pkg/job"
	"os"
)

// Stores the configuration
type Config struct {
	Jobs        []job.Job `json:"jobs"`
	Concurrency int       `json:"concurrency"`
}

// Parses the specified config file and loads the data into a config object
func ParseFile(filepath string) Config {
	// read in the config file
	configFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}

	// convert the JSON into a config object
	c := Config{}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&c)
	if err != nil {
		panic(err)
	}
	return c
}
