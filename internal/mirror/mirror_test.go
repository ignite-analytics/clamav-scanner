package mirror

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fsouza/fake-gcs-server/fakestorage"
)

const TestBucketName = "some-bucket"

func Test_Blob_Exists(t *testing.T) {
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: TestBucketName,
				Name:       "some/object/file",
			},
		},
	})
	defer server.Stop()
	client := server.Client()
	storageProxy := NewStorageProxy(client.Bucket(TestBucketName))

	req, err := http.NewRequest(http.MethodHead, "/mirror/some/object/file", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storageProxy.checkBlobExists(w, r, "some/object/file")
	})
	handler.ServeHTTP(response, req)

	if response.Code == http.StatusOK {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d'", response.Code)
	}
}

func Test_Blob_Download(t *testing.T) {
	expectedBlobContent := "my content"
	server := fakestorage.NewServer([]fakestorage.Object{
		{
			ObjectAttrs: fakestorage.ObjectAttrs{
				BucketName: TestBucketName,
				Name:       "some/file",
			},
			Content: []byte(expectedBlobContent),
		},
	})
	defer server.Stop()
	client := server.Client()
	storageProxy := NewStorageProxy(client.Bucket(TestBucketName))

	req, err := http.NewRequest(http.MethodGet, "/mirror/some/file", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}

	response := httptest.NewRecorder()
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		storageProxy.downloadBlob(w, r, "some/file")
	})
	handler.ServeHTTP(response, req)

	if response.Code == http.StatusOK {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong status: '%d'", response.Code)
	}

	downloadedBlobContent := response.Body.String()
	if downloadedBlobContent == expectedBlobContent {
		t.Log("Passed")
	} else {
		t.Errorf("Wrong content: '%s'", downloadedBlobContent)
	}
}
