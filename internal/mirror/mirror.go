package mirror

import (
	"bufio"
	"context"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
)

type StorageProxy struct {
	bucketHandler *storage.BucketHandle
}

func NewStorageProxy(bucketHandler *storage.BucketHandle) *StorageProxy {
	return &StorageProxy{
		bucketHandler: bucketHandler,
	}
}

// Handle is the HTTP handler function for the mirror route
func Handle(mirrorBucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		client, err := storage.NewClient(ctx)
		if err != nil {
			log.Printf("Failed to create storage client: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		defer func() {
			if err := client.Close(); err != nil {
				log.Printf("Failed to close storage client: %v", err)
			}
		}()

		bucketHandler := client.Bucket(mirrorBucket)
		storageProxy := NewStorageProxy(bucketHandler)

		key := r.URL.Path[len("/mirror/"):]

		switch r.Method {
		case "GET":
			storageProxy.downloadBlob(w, r, key)
		case "HEAD":
			storageProxy.checkBlobExists(w, r, key)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

func (proxy *StorageProxy) downloadBlob(w http.ResponseWriter, r *http.Request, name string) {
	object := proxy.bucketHandler.Object(name)
	if object == nil {
		http.NotFound(w, r)
		return
	}

	reader, err := object.NewReader(context.Background())
	if err != nil {
		http.NotFound(w, r)
		return
	}

	defer func() {
		if err := reader.Close(); err != nil {
			log.Printf("Failed to close reader: %v", err)
		}
	}()
	bufferedReader := bufio.NewReader(reader)
	_, err = bufferedReader.WriteTo(w)

	if err != nil {
		log.Printf("Failed to serve blob '%s': %v", name, err)
		http.Error(w, "Failed to serve blob", http.StatusInternalServerError)
	}
}

func (proxy *StorageProxy) checkBlobExists(w http.ResponseWriter, r *http.Request, name string) {
	object := proxy.bucketHandler.Object(name)
	if object == nil {
		http.NotFound(w, r)
		return
	}

	attrs, err := object.Attrs(context.Background())
	if err != nil || attrs == nil {
		http.NotFound(w, r)
		return
	}

	w.WriteHeader(http.StatusOK)
}
