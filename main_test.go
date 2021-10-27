package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEndpoints(t *testing.T) {
	tests := []struct {
		description string
		route       string

		// Request
		method string
		body   io.Reader

		// Expected output
		expectedError bool
		expectedCode  int
		expectedBody  string
	}{
		{
			description:   "get default message",
			route:         "/api/v1/assist/",
			method:        "PUT",
			body:          nil,
			expectedError: false,
			expectedCode:  200,
			expectedBody:  `{"ID":1,"CreatedAt":"000","UpdatedAt":"000","DeletedAt":"000","body":"Hello, I am a virtual assistant. How can I help you?","options":null,"FlowID":1}`,
		},
		{
			description:   "get message with id = 1",
			route:         "/api/v1/assist/1",
			method:        "PUT",
			body:          nil,
			expectedError: false,
			expectedCode:  200,
			expectedBody:  `{"ID":1,"CreatedAt":"000","UpdatedAt":"000","DeletedAt":"000","body":"Hello, I am a virtual assistant. How can I help you?","options":null,"FlowID":1}`,
		},
		{
			description:   "get 404",
			route:         "/api/v1/assist/100",
			method:        "PUT",
			body:          nil,
			expectedError: false,
			expectedCode:  404,
			expectedBody:  `{"message":"no messages found","status":"error"}`,
		},
		{
			description:   "Get all records from assistant db",
			route:         "/api/v1/assistant/db",
			method:        "GET",
			body:          nil,
			expectedError: false,
			expectedCode:  200,
			expectedBody:  `[{"ID":1,"CreatedAt":"000","UpdatedAt":"000","DeletedAt":"000","body":"Hello, I am a virtual assistant. How can I help you?","options":null,"FlowID":1},{"ID":2,"CreatedAt":"000","UpdatedAt":"000","DeletedAt":"000","body":"Let me clarify what exactly you need?","options":null,"FlowID":1},{"ID":3,"CreatedAt":"000","UpdatedAt":"000","DeletedAt":"000","body":"If you need restore your password please fallow next [link](http://help.com/pwd)","options":null,"FlowID":1},{"ID":7,"CreatedAt":"000","UpdatedAt":"000","DeletedAt":"000","body":"Thank you for using our service! Have good day!","options":null,"FlowID":1}]`,
		},
		{
			description:   "Get record by ID",
			route:         "/api/v1/assistant/db/1",
			method:        "GET",
			body:          nil,
			expectedError: false,
			expectedCode:  200,
			expectedBody:  `{"ID":1,"CreatedAt":"000","UpdatedAt":"000","DeletedAt":"000","body":"Hello, I am a virtual assistant. How can I help you?","options":null,"FlowID":1}`,
		},
		{
			description:   "Get record - 404",
			route:         "/api/v1/assistant/db/100",
			method:        "GET",
			body:          nil,
			expectedError: false,
			expectedCode:  404,
			expectedBody:  `{"message":"message not found","status":"error"}`,
		},
	}

	// Setup the app as it is done in the main function
	testDb, err := SetupDB("test.db")
	if err != nil {
		t.Error(err.Error())
	}

	store := NewFlowStorage(testDb)
	app := Setup(store)

	for _, tt := range tests {
		// Create request
		req, _ := http.NewRequest(
			tt.method,
			tt.route,
			tt.body,
		)
		req.Header.Set("Content-Type", "application/json; charset=UTF-8")

		res, err := app.Test(req, -1)

		assert.Equalf(t, tt.expectedError, err != nil, tt.description)

		if tt.expectedError {
			continue
		}

		assert.Equalf(t, tt.expectedCode, res.StatusCode, tt.description)

		// Read the response body
		body, err := ioutil.ReadAll(res.Body)
		assert.Nilf(t, err, tt.description)

		gotBody := string(body)
		gotBody, err = cleanupDates(gotBody)
		if err != nil {
			t.Fatal(err)
		}

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, tt.expectedBody, string(gotBody), tt.description)
	}
	os.Remove("test.db")
}

// (.*["CreatedAt"||"UpdatedAt"||"DeletedAt"]:)(["\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d{9}-\d{2}:\d{2}"])
func cleanupDates(str string) (string, error) {
	regexp, err := regexp.Compile(`("CreatedAt"|"UpdatedAt"|"DeletedAt"):("\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}.\d*-\d{2}:\d{2}"|null)`)
	if err != nil {
		return "", err
	}
	return regexp.ReplaceAllString(str, `$1:"000"`), nil
}
