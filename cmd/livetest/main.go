package main

import (
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/rknizzle/livetest/pkg/config"
	"github.com/rknizzle/livetest/pkg/datastore"
	"github.com/rknizzle/livetest/pkg/datastore/postgres"
	"github.com/rknizzle/livetest/pkg/scheduler"
	"os"
	"time"
)

func main() {
	var file string
	// check for a filename supplied as an argument
	if len(os.Args) > 1 {
		if !fileExists(os.Args[1]) {
			fmt.Printf("%s not found", os.Args[1])
			return
		}
		file = os.Args[1]
	} else {
		// if there are no arguments provided look for config.json
		if !fileExists("config.json") {
			fmt.Println("config.json not found")
			return
		}
		file = "config.json"
	}
	// parse the config file
	config := config.Parse(file)

	hasDatastore := false
	var store datastore.Datastore
	// create the connection to the datastore from the info in the config
	if config.Datastore != nil {
		// only connect to postgres for now
		store = &postgres.Postgres{}

		connect := func() error {
			fmt.Println("Attempting database connection...")
			err := store.Connect(config.Datastore)
			if err != nil {
				return err
			}
			return nil
		}

		fmt.Println("Starting connection to database")
		// attempt to connect to the database with exponential backoff
		bo := backoff.NewExponentialBackOff()
		bo.MaxElapsedTime = 30 * time.Second
		err := backoff.Retry(connect, bo)
		if err != nil {
			// exit if not able to make connection with database
			fmt.Println("Database connection failed. Exiting now...")
			return
		}

		hasDatastore = true
	} else {
		fmt.Println("No datastore connection info in config. Running without database.")
	}

	// blocking channel
	// concurrency limit is specified in the config
	blocker := make(chan struct{}, config.Concurrency)
	// channel of request results
	resChan := make(chan scheduler.Result)
	// loop through each job
	for _, j := range config.Jobs {
		// and schedule the job to run at the specified interval
		scheduler.Schedule(j, time.Duration(j.Frequency)*time.Millisecond, resChan, blocker)
	}
	for res := range resChan {
		// check if the request has the expected response
		record := scheduler.HandleResponse(res, config.Jobs, config.Notification)

		// and write the data to the datastore
		if hasDatastore {
			store.Write(record)
		}
	}
}

// Function found at:
// https://golangcode.com/check-if-a-file-exists/ (MIT LICENSE)
// fileExists checks if a file exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
