package main

import (
	"fmt"
	"github.com/rknizzle/livetest/pkg/parser"
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
	config := parser.ParseFile(file)
	// buffered response channel
	// concurrency limit is specified in the config
	resChan := make(chan scheduler.Result, config.Concurrency)
	// loop through each job
	for _, j := range config.Jobs {
		// and schedule the job to run at the specified interval
		scheduler.Schedule(j, time.Duration(j.Frequency)*time.Millisecond, resChan)
	}
	for res := range resChan {
		for _, j := range config.Jobs {
			if j.Title == res.Title {
				// check for error
				if res.Err != nil {
					fmt.Println("Failing with error:")
					fmt.Println(res.Err)
					// check if response status code matches the expected status code
				} else if res.Res.StatusCode == j.ExpectedResponse.StatusCode {
					fmt.Println("passing")
				} else {
					fmt.Println("failing")
				}
			}
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
