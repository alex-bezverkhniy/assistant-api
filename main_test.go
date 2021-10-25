package main

import (
	"io"
	"io/ioutil"
	"net/http"
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
			expectedBody:  `{"id":0,"body":"Hello, I am a virtual assistant. How can I help you?","options":[{"id":0,"body":"I need help with my password","nextMessageId":1},{"id":1,"body":"I need help with my account","nextMessageId":2}]}`,
		},
		{
			description:   "get message with id = 1",
			route:         "/api/v1/assist/1",
			method:        "PUT",
			body:          nil,
			expectedError: false,
			expectedCode:  200,
			expectedBody:  `{"id":1,"body":"Let me clarify what exactly you need?","options":[{"id":0,"body":"restore password","nextMessageId":3},{"id":1,"body":"change password","nextMessageId":4}]}`,
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
			expectedBody:  `[{"id":0,"body":"Hello, I am a virtual assistant. How can I help you?","options":[{"id":0,"body":"I need help with my password","nextMessageId":1},{"id":1,"body":"I need help with my account","nextMessageId":2}]},{"id":1,"body":"Let me clarify what exactly you need?","options":[{"id":0,"body":"restore password","nextMessageId":3},{"id":1,"body":"change password","nextMessageId":4}]},{"id":2,"body":"Let me clarify what exactly you need?","options":[{"id":0,"body":"unclock my account","nextMessageId":5},{"id":1,"body":"block my account","nextMessageId":6}]},{"id":3,"body":"If you need restore your password please fallow next [link](http://help.com/pwd)","options":[{"id":0,"body":"I have other questions","nextMessageId":8},{"id":1,"body":"Thank you I got what I need","nextMessageId":7}]},{"id":7,"body":"Thank you for using our service! Have good day!","options":[]}]`,
		},
		{
			description:   "Get record by ID",
			route:         "/api/v1/assistant/db/1",
			method:        "GET",
			body:          nil,
			expectedError: false,
			expectedCode:  200,
			expectedBody:  `{"id":1,"body":"Let me clarify what exactly you need?","options":[{"id":0,"body":"restore password","nextMessageId":3},{"id":1,"body":"change password","nextMessageId":4}]}`,
		},
		{
			description:   "Get record - 404",
			route:         "/api/v1/assistant/db/100",
			method:        "GET",
			body:          nil,
			expectedError: false,
			expectedCode:  404,
			expectedBody:  `{"message":"no messages found","status":"error"}`,
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
