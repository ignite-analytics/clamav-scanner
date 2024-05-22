package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var isProcessRunningMock func(name string) (bool, error)

func TestHandle(t *testing.T) {
	originalIsProcessRunning := isProcessRunning
	defer func() { isProcessRunning = originalIsProcessRunning }()

	isProcessRunning = func(name string) (bool, error) {
		return isProcessRunningMock(name)
	}

	tests := []struct {
		name           string
		processRunning bool
		processError   error
		expectedStatus int
	}{
		{
			name:           "Process running",
			processRunning: true,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Process not running",
			processRunning: false,
			expectedStatus: http.StatusServiceUnavailable,
		},
		{
			name:           "Process error",
			processError:   assert.AnError,
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			isProcessRunningMock = func(name string) (bool, error) {
				return tt.processRunning, tt.processError
			}

			req, err := http.NewRequest(http.MethodGet, "/health", nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			rr := httptest.NewRecorder()
			handler := Handle()
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code, "unexpected status code")
		})
	}

	t.Run("Method not allowed", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodPost, "/health", nil)
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}

		rr := httptest.NewRecorder()
		handler := Handle()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusMethodNotAllowed, rr.Code, "unexpected status code")
		assert.Equal(t, "Method Not Allowed\n", rr.Body.String(), "unexpected body")
	})
}
