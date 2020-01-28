// Package scheduler creates tickers to run each job periodically at the correct interval

package scheduler

import (
	"github.com/rknizzle/livetest/pkg/job"
	"net/http"
	"time"
)

type Result struct {
	ID  int
	Res http.Response
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

	// configure the HTTP request
	req, err := http.NewRequest(job.HTTPMethod, job.URL, nil)
	if err != nil {
		panic(err)
	}

	// execute the HTTP request
	resp, err := client.Do(req)

	r := Result{job.ID, *resp, err}
	resChan <- r
}
