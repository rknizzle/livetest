// Package scheduler creates tickers to run each job periodically at the correct interval

package scheduler

import (
	"bytes"
	"encoding/json"
	"github.com/rknizzle/livetest/pkg/job"
	"io"
	"net/http"
	"time"
)

type Result struct {
	ID  int
	Res *http.Response
	Err error
}

// Create a ticker for a job to trigger it periodically
func Schedule(j *job.Job, interval time.Duration, resChan chan<- Result) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				// run the job
				execute(j, resChan)
			}
		}
	}()
	return ticker
}

// Sends an HTTP request for a job
func execute(job *job.Job, resChan chan<- Result) {
	client := &http.Client{}

	// set body to nil unless there is a RequestBody to add
	var body io.Reader = nil
	// if the RequestBody map is not empty
	if len(job.RequestBody) > 0 {
		// convert the body into an io.Reader to pass to the http request
		requestByte, err := json.Marshal(job.RequestBody)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(requestByte)
	}

	// configure http request
	req, err := http.NewRequest(job.HTTPMethod, job.URL, body)
	if err != nil {
		panic(err)
	}
	// if the request contains a body, set the content type to json
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// execute the HTTP request
	resp, err := client.Do(req)

	r := Result{job.ID, resp, err}
	resChan <- r
}
