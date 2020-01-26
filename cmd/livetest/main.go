package main

import (
	"fmt"
	"github.com/rknizzle/livetest/pkg/parser"
	"github.com/rknizzle/livetest/pkg/scheduler"
	"net/http"
	"time"
)

func main() {
	config := parser.ParseFile("examples/config.json")
	// buffered response channel
	// concurrency limit is specified in the config
	resChan := make(chan *http.Response, config.Concurrency)
	// loop through each job
	for _, j := range config.Jobs {
		// and schedule the job to run at the specified interval
		scheduler.Schedule(j, time.Duration(j.Frequency)*time.Millisecond, resChan)
	}
	for res := range resChan {
		// just print out the HTTP response for now
		fmt.Println("got a result")
		fmt.Println(res)
	}
}
