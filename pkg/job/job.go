// Package job represents a feature or endpoint to be tested.
// The result of the job is checked against an expected result
// to express a passing or failing feature

package job

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
