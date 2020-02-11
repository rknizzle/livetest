package main

import (
	"fmt"
	"github.com/cenkalti/backoff"
	"github.com/rknizzle/livetest/config"
	"github.com/rknizzle/livetest/datastore"
	"github.com/rknizzle/livetest/datastore/postgres"
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
	config, err := config.Parse(file)
	if err != nil {
		fmt.Println(err)
		return
	}

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

		fmt.Println("Successfully connected to database!")
		hasDatastore = true
	} else {
		fmt.Println("No datastore connection info in config. Running without database.")
	}

	// blocking channel
	// concurrency limit is specified in the config
	blocker := make(chan struct{}, config.Concurrency)
	// channel of job response as a record
	recChan := make(chan datastore.Record)
	// loop through each job
	for _, j := range config.Jobs {
		// and schedule the job to run at the specified interval
		j.Schedule(recChan, blocker, config.Notification)
	}
	for rec := range recChan {
		// and write the data record to the datastore
		if hasDatastore {
			store.Write(&rec)
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
