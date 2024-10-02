package scan

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TestBucketName       = "test-bucket"
	TestQuarantineBucket = "quarantine-bucket"
	TestFileName         = "test-file.txt"
	TestFileContent      = "test file content"
)

func TestHandle(t *testing.T) {
	originalPerformScan := performScanFunc
	defer func() { performScanFunc = originalPerformScan }()

	tests := []struct {
		name             string
		method           string
		body             interface{}
		mockScanResult   bool
		mockScanError    error
		expectedStatus   int
		expectedResponse map[string]string
	}{
		{
			name:           "Successful scan with no threats",
			method:         "POST",
			body:           requestBody{Kind: "test", Name: TestFileName, Bucket: TestBucketName},
			mockScanResult: true,
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]string{
				"message": "Scan completed for " + TestFileName,
				"result":  "NO_THREAT",
			},
		},
		{
			name:           "Successful scan with threats",
			method:         "POST",
			body:           requestBody{Kind: "test", Name: TestFileName, Bucket: TestBucketName},
			mockScanResult: false,
			expectedStatus: http.StatusOK,
			expectedResponse: map[string]string{
				"message": "Scan completed for " + TestFileName,
				"result":  "THREAT",
			},
		},
		{
			name:           "Method not allowed",
			method:         "GET",
			body:           nil,
			expectedStatus: http.StatusMethodNotAllowed,
		},
		{
			name:           "Bad request",
			method:         "POST",
			body:           "invalid json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			performScanFunc = func(ctx context.Context, bucketName, fileName, quarantineBucket string) (bool, error) {
				return tt.mockScanResult, tt.mockScanError
			}

			var reqBody []byte
			if tt.body != nil {
				reqBody, _ = json.Marshal(tt.body)
			}

			req, err := http.NewRequest(tt.method, "/scan", bytes.NewReader(reqBody))
			assert.NoError(t, err)

			rr := httptest.NewRecorder()
			handler := Handle(TestQuarantineBucket)
			handler.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedStatus == http.StatusOK {
				var respBody map[string]string
				err = json.NewDecoder(rr.Body).Decode(&respBody)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedResponse, respBody)
			}
		})
	}
}
