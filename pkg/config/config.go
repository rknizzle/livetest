// Package parser parses the config file and creates jobs, connects to database and sets up notifications

package config

import (
	"encoding/json"
	"fmt"
	"github.com/rknizzle/livetest/pkg/datastore"
	"github.com/rknizzle/livetest/pkg/job"
	"github.com/rknizzle/livetest/pkg/notification"
	"os"
)

// Stores the testing configuration
type Config struct {
	Jobs         []*job.Job                `json:"jobs"`
	Concurrency  int                       `json:"concurrency"`
	Datastore    *datastore.Connection     `json:"datastore"`
	Notification notification.Notification `json:"notification"`
}

// Config before unpacking the envelopes
type Pre struct {
	Jobs         []*job.Job            `json:"jobs"`
	Concurrency  int                   `json:"concurrency"`
	Datastore    *datastore.Connection `json:"datastore"`
	Notification Envelope              `json:"notification"`
}

// intermediate data structure for reading in dynamic json
type Envelope struct {
	Type string           `json:"type"`
	Msg  *json.RawMessage `json:"msg"`
}

func initializeNotification(e Envelope) notification.Notification {
	switch e.Type {
	case "http":
		var h notification.HTTPRequest
		if err := json.Unmarshal(*e.Msg, &h); err != nil {
			panic(err)
		}
		return h
	default:
		fmt.Println("No valid notification type in config. Running without notification.")
		return nil
	}
}

// Parses the specified config file and loads the data into a config object
func Parse(filepath string) *Config {
	var pre Pre
	// read in the config file
	configFile, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&pre)
	if err != nil {
		panic(err)
	}

	n := initializeNotification(pre.Notification)

	// create the config object from JSON data
	c := &Config{pre.Jobs, pre.Concurrency, pre.Datastore, n}

	return c
}
