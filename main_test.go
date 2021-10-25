package main

import (
	"io"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssist(t *testing.T) {
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
			expectedBody:  `{"id":0,"body":"Hello, I am a virtual asistant. How can I help you?","options":[{"id":0,"body":"I need help with my passowrd","nextMessageId":1},{"id":1,"body":"I need help with my account","nextMessageId":2}]}`,
		},
	}

	// Setup the app as it is done in the main function
	app := Setup()

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

		// Verify, that the reponse body equals the expected body
		assert.Equalf(t, tt.expectedBody, string(body), tt.description)
	}
}

func AssistantBDTest(t *testing.T) {

}
