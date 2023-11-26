package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostUsersHandler(t *testing.T) {
	// Initialize your mock database
	mockDB := &MockDB{}

	// Create an instance of your apiConfig with the mockDB
	cfg := &apiConfig{
		db: mockDB,
		// Add other configuration properties as needed
	}

	// Create a test server with the postUsersHandler as the handler
	server := httptest.NewServer(http.HandlerFunc(cfg.postUsersHandler))
	defer server.Close()

	// Your test cases
	testCases := []struct {
		name       string
		request    map[string]interface{}
		statusCode int
	}{
		{
			name: "Valid Request",
			request: map[string]interface{}{
				"email":        "test@example.com",
				"password":     "testpassword",
				"display_name": "testuser",
				"phone":        "123456789",
			},
			statusCode: http.StatusCreated,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Convert the request map to JSON
			reqBody, err := json.Marshal(tc.request)
			assert.NoError(t, err)

			// Make a POST request to the test server
			resp, err := http.Post(server.URL, "application/json", bytes.NewBuffer(reqBody))
			assert.NoError(t, err)
			defer resp.Body.Close()

			// Check the response status code
			assert.Equal(t, tc.statusCode, resp.StatusCode)

			// If the response code is StatusCreated, you might want to parse the response body as JSON
			if tc.statusCode == http.StatusCreated {
				var resBody map[string]interface{}
				err := json.NewDecoder(resp.Body).Decode(&resBody)
				assert.NoError(t, err)

				// Add more assertions on the response body if needed
			}
		})
	}
}
