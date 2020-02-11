package job

import (
	"fmt"
	"github.com/rknizzle/livetest/datastore"
	"github.com/rknizzle/livetest/notification"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func Test(t *testing.T) {
	// generate a test server to mock the requests
	testServer := httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(200)
		res.Write([]byte("body"))
	}))
	defer func() { testServer.Close() }()

	// example job
	j := Job{
		Title:      "example",
		URL:        testServer.URL,
		HTTPMethod: "GET",
		Headers:    make(map[string]interface{}),
		Frequency:  500,
		Success:    true,
		ExpectedResponse: Response{
			StatusCode: 200,
			Body:       "",
		},
		RequestBody: make(map[string]interface{}),
	}

	// run the tests
	t.Run("Execute", testExecute(j))
	t.Run("Execute with missing URL", testExecuteWithError(j))
	t.Run("HandleResponse", testHandleResponse(j))
	t.Run("Handle Response appropriately triggrs a notification", testNotificationFromHandleResponse(j))
	t.Run("Schedule", testSchedule(j))
}

// run execute and expect a 200 response
func testExecute(j Job) func(*testing.T) {
	return func(t *testing.T) {
		res := j.execute()
		fmt.Println("response:")
		fmt.Println(res.res.StatusCode)
		if res.res.StatusCode != 200 {
			t.Error("Response status code was not 200")
		}
	}
}

// run execute with a missing URL and it should return an error and no response
func testExecuteWithError(j Job) func(*testing.T) {
	return func(t *testing.T) {
		// test execute with a job with a missing URL
		j.URL = ""
		res := j.execute()
		if res.err == nil || res.res != nil {
			t.Error("Execute did not return an error with a bad job")
		}
	}
}

// HandleResponse returns the appropriate data record
func testHandleResponse(j Job) func(*testing.T) {
	return func(t *testing.T) {
		response := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
		}
		result := result{response, nil, 0 * time.Second}
		expectedRecord := datastore.Record{true, "example", 200, 0 * time.Second}
		record, _ := j.HandleResponse(result)
		if record != expectedRecord {
			t.Error("Incorrect record created")
		}
	}
}

// test that the notification boolean is set when the response isnt expected
func testNotificationFromHandleResponse(j Job) func(*testing.T) {
	return func(t *testing.T) {
		j.ExpectedResponse.StatusCode = 0
		response := &http.Response{
			Status:     "200 OK",
			StatusCode: 200,
		}
		result := result{response, nil, 0 * time.Second}
		_, shouldNotify := j.HandleResponse(result)
		if shouldNotify == false {
			t.Error("shouldNotify boolean not set to true")
		}
	}
}

// test that Schedule will return the correct data record into the record channel
func testSchedule(j Job) func(*testing.T) {
	return func(t *testing.T) {
		blocker := make(chan struct{}, 1)
		recChan := make(chan datastore.Record)
		notification := notification.HTTPRequest{}
		j.Schedule(recChan, blocker, notification)

		expectedRecord := datastore.Record{true, "example", 200, 0 * time.Second}

		// fail if no record is recieved in the channel in 10 seconds
		go func() {
			time.Sleep(10 * time.Second)
			t.Error("TIMED OUT")
		}()

		for rec := range recChan {
			rec.Duration = 0 * time.Second

			if rec != expectedRecord {
				fmt.Println(rec)
				t.Error("Database record did not equal expected")
			} else {
				return
			}
		}
	}
}
