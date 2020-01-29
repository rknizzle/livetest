// Package scheduler creates tickers to run each job periodically at the correct interval

package scheduler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rknizzle/livetest/pkg/datastore"
	"github.com/rknizzle/livetest/pkg/job"
	"io"
	"net/http"
	"time"
)

type Result struct {
	Title string
	Res   *http.Response
	Err   error
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

	r := Result{job.Title, resp, err}
	resChan <- r
}

func HandleResponse(res Result, jobs []*job.Job) *datastore.Record {
	for _, j := range jobs {
		if j.Title == res.Title {
			// check for error
			if res.Err != nil {
				j.Status = "failing"
				fmt.Println("Failing with error:")
				fmt.Println(res.Err)
				// check if response status code matches the expected status code
			} else if res.Res.StatusCode == j.ExpectedResponse.StatusCode {
				j.Status = "passing"
				fmt.Println("passing")
			} else {
				j.Status = "failing"
				fmt.Println("failing")
			}
			// turn the result into a data record
			return &datastore.Record{j.Status, j.Title, res.Res.StatusCode}
		}
	}
	return nil
}
