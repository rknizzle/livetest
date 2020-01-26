// Package scheduler creates tickers to run each job periodically at the correct interval

package scheduler

import (
	"github.com/rknizzle/livetest/pkg/job"
	"net/http"
	"time"
)

// Create a ticker for a job to trigger it periodically
func Schedule(j job.Job, interval time.Duration, resChan chan<- *http.Response) *time.Ticker {
	ticker := time.NewTicker(interval)
	go func() {
		for {
			select {
			case <-ticker.C:
				// run the job
				execute(&j, resChan)
			}
		}
	}()
	return ticker
}

// Sends an HTTP request for a job
func execute(job *job.Job, resChan chan<- *http.Response) {
	client := &http.Client{}

	// configure the HTTP request
	r, err := http.NewRequest(job.HTTPMethod, job.URL, nil)
	if err != nil {
		panic(err)
	}

	// execute the HTTP request
	resp, _ := client.Do(r)
	// TODO: pass in and handle errors as well
	resChan <- resp
}
