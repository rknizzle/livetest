//Package notification hadnles the notifications to notify the user if a job is not passing

package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Notification interface {
	Notify()
}

type HTTPRequest struct {
	URL         string                 `json:"url"`
	HTTPMethod  string                 `json:"httpMethod"`
	RequestBody map[string]interface{} `json:"requestBody"`
	Headers     map[string]interface{} `json:"headers"`
}

func (h HTTPRequest) Notify() {
	client := &http.Client{}

	// set body to nil unless there is a RequestBody to add
	var body io.Reader = nil
	// if the RequestBody map is not empty
	if len(h.RequestBody) > 0 {
		// convert the body into an io.Reader to pass to the http request
		requestByte, err := json.Marshal(h.RequestBody)
		if err != nil {
			panic(err)
		}
		body = bytes.NewBuffer(requestByte)
	}

	// configure http request
	req, err := http.NewRequest(h.HTTPMethod, h.URL, body)
	if err != nil {
		panic(err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// execute the HTTP request
	_, err = client.Do(req)
	if err != nil {
		panic(err)
	}
}
