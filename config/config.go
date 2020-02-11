// Package parser parses the config file and creates jobs, connects to database and sets up notifications

package config

import (
	"encoding/json"
	"errors"
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

// take in notification data from config and return either a Notification object,
// an error if something goes wrong, or nil if there was no notification specified in config
func initializeNotification(e Envelope) (notification.Notification, error) {
	// if envelope is empty return nil notification
	if e == (Envelope{}) {
		return nil, nil
	}
	switch e.Type {
	case "http":
		var h notification.HTTPRequest
		if err := json.Unmarshal(*e.Msg, &h); err != nil {
			return nil, errors.New("Couldn't parse HTTP notification data from config")
		}
		return h, nil
	default:
		// return error if the notification type is not supported
		return nil, errors.New(e.Type + " is not a supported notification type")
	}
}

// Parses the specified config file and loads the data into a config object
func Parse(filepath string) (*Config, error) {
	var pre Pre
	// read in the config file
	configFile, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	jsonParser := json.NewDecoder(configFile)
	err = jsonParser.Decode(&pre)
	if err != nil {
		return nil, err
	}

	// initialize notification from config value
	n, err := initializeNotification(pre.Notification)
	if err != nil {
		return nil, err
	}

	// create the config object from JSON data
	c := &Config{pre.Jobs, pre.Concurrency, pre.Datastore, n}

	return c, nil
}
