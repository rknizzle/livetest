// Package job represents a feature or endpoint to be tested.
// The result of the job is checked against an expected result
// to express a passing or failing feature

package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/rknizzle/livetest/pkg/datastore"
	"github.com/rknizzle/livetest/pkg/notification"
	"io"
	"net/http"
	"time"
)

type Job struct {
	Title            string                 `json:"title"`
	URL              string                 `json:"url"`
	HTTPMethod       string                 `json:"httpMethod"`
	Headers          map[string]interface{} `json:"headers"`
	Frequency        int                    `json:"frequency"`
	Success          bool                   `json:"passing"`
	ExpectedResponse Response               `json:"expectedResponse"`
	RequestBody      map[string]interface{} `json:"requestBody"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

// Create a ticker for a job to trigger it periodically
func (j Job) Schedule(recChan chan<- datastore.Record, blocker chan struct{}, n notification.Notification) {
	ticker := time.NewTicker(time.Duration(j.Frequency) * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				// run the job

				// send an empty struct to the blocker when the concurrency limit
				// has been reached it will block until there is room in the channel
				blocker <- struct{}{}
				// execute the job
				res := j.Execute()
				// handle the request response and format it into a database record
				record, shouldNotify := j.HandleResponse(res)
				// send the notification if the job failed
				if shouldNotify {
					if n != nil {
						n.Notify()
					}
				}

				// send the database record to the record channel
				recChan <- record
				// remove struct from the blocker channel to open up space for another request to run
				<-blocker
			}
		}
	}()
}

// store the result of an http request
// either a response or an error
type result struct {
	res *http.Response
	err error
}

// Sends an HTTP request for a job
func (j Job) Execute() result {
	//blocker <- struct{}{}
	client := &http.Client{}

	// set body to nil unless there is a RequestBody to add
	var body io.Reader = nil
	// if the RequestBody map is not empty
	if len(j.RequestBody) > 0 {
		// convert the body into an io.Reader to pass to the http request
		requestByte, err := json.Marshal(j.RequestBody)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(requestByte)
	}

	// configure http request
	req, err := http.NewRequest(j.HTTPMethod, j.URL, body)
	if err != nil {
		panic(err)
	}

	// loop through each header and add it to request
	for k, v := range j.Headers {
		req.Header.Set(k, v.(string))
	}

	// execute the HTTP request
	resp, err := client.Do(req)

	return result{resp, err}
}

// verify if the request got the expected result and format the data
// for writing to the database
func (j Job) HandleResponse(res result) (datastore.Record, bool) {
	shouldNotify := false
	log := j.Title + ": "
	// check for error
	if res.err != nil {
		j.Success = false
		log += "Fail with error: "
		log += res.err.Error()
		// check if response status code matches the expected status code
	} else if res.res.StatusCode == j.ExpectedResponse.StatusCode {
		j.Success = true
		log += "Success"
	} else {
		j.Success = false
		log += "Fail"
	}
	fmt.Println(log)
	if j.Success == false {
		// job has failed so trigger a notification if there is one
		shouldNotify = true
	}

	// turn the result into a data record
	if res.res == nil {
		return datastore.Record{j.Success, j.Title, 0}, shouldNotify
	}
	return datastore.Record{j.Success, j.Title, res.res.StatusCode}, shouldNotify
}
