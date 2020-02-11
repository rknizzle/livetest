package config

import (
	"encoding/json"
	"github.com/rknizzle/livetest/datastore"
	"github.com/rknizzle/livetest/job"
	"github.com/rknizzle/livetest/notification"
	"reflect"
	"testing"
)

func getExampleConfig() *Config {
	j := &job.Job{
		Title:      "example GET",
		URL:        "http://postman-echo.com/get?foo1=bar1&foo2=bar2",
		HTTPMethod: "GET",
		ExpectedResponse: job.Response{
			StatusCode: 200,
		},
		RequestBody: map[string]interface{}{},
		Headers:     map[string]interface{}{},
		Frequency:   5000,
	}
	n := notification.HTTPRequest{
		URL:         "http://postman-echo.com/post",
		HTTPMethod:  "POST",
		RequestBody: map[string]interface{}{"notification": "true"},
		Headers:     map[string]interface{}{"Content-Type": "application/json"},
	}

	d := &datastore.Connection{
		Host:     "localhost",
		Port:     5432,
		User:     "postgres",
		Password: "password",
		DBname:   "postgres",
	}

	var jobArray []*job.Job
	jobArray = append(jobArray, j)

	c := &Config{
		Jobs:         jobArray,
		Concurrency:  1,
		Datastore:    d,
		Notification: n,
	}
	return c
}

// Test that parse returns the expected Config object
func TestParse(t *testing.T) {
	inputFile := "testdata/config.json"
	expectedConfig := getExampleConfig()
	c, err := Parse(inputFile)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(c, expectedConfig) {
		t.Error("Config object doesnt match expected")
	}
}

// Test that initializeNotification with correct input will not return an error
func TestInitializeNotification(t *testing.T) {
	data := []byte(`{
            "headers": {
                "Content-Type": "application/json"
            },
            "httpMethod": "POST",
            "requestBody": {
                "notification": "true"
            },
            "url": "http://postman-echo.com/post"
        }`)
	e := Envelope{
		Type: "http",
		Msg:  (*json.RawMessage)(&data),
	}

	_, err := initializeNotification(e)
	if err != nil {
		t.Error(err)
	}

}

// test that initializeNotification returns an error if an unsupported notification type is specified
func TestInitializeNotificationWithInvalidType(t *testing.T) {
	data := []byte(`{
            "headers": {
                "Content-Type": "application/json"
            },
            "httpMethod": "POST",
            "requestBody": {
                "notification": "true"
            },
            "url": "http://postman-echo.com/post"
        }`)
	e := Envelope{
		Type: "unsupported",
		Msg:  (*json.RawMessage)(&data),
	}
	_, err := initializeNotification(e)
	if err == nil {
		t.Error("An error should have been returned but it was not")
	}
}

// Test that Initialize function returns nils when an empty envelope is passed in
func TestInitializeNotificationWithNoValue(t *testing.T) {
	e := Envelope{}
	n, err := initializeNotification(e)
	if err != nil {
		t.Error(err)
	}
	if n != nil {
		t.Error("Expected notification to be nil")
	}
}
