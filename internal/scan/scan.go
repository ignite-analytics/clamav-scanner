package scan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"cloud.google.com/go/storage"
	"github.com/lyimmi/go-clamd"
)

type requestBody struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Bucket string `json:"bucket"`
}

var performScanFunc = performScan

// Handle is the HTTP handler used to scan files using ClamAV.
func Handle(quarantineBucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		ctx := r.Context()
		var reqBody requestBody
		if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
			log.Printf("Error decoding request body: %v", err)
			http.Error(w, "Bad request", http.StatusBadRequest)
			return
		}
		defer r.Body.Close()

		log.Printf("Received scan request for bucket: %s, file: %s\n", reqBody.Bucket, reqBody.Name)

		safe, err := performScanFunc(ctx, reqBody.Bucket, reqBody.Name, quarantineBucket)
		if err != nil {
			log.Printf("Error scanning file: %v", err)
			http.Error(w, "Failed to scan file", http.StatusInternalServerError)
			return
		}

		resultMsg := "No threats found"
		if !safe {
			resultMsg = "Threat found"
		}

		log.Printf("Scan completed for %s: %s\n", reqBody.Name, resultMsg)

		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Scan completed for " + reqBody.Name,
			"result":  resultMsg,
		}); err != nil {
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		}

	}
}

func performScan(ctx context.Context, bucketName, fileName, quarantineBucket string) (bool, error) {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to create storage client: %w", err)
	}
	defer storageClient.Close()

	fileData, err := fetchFileFromBucket(ctx, storageClient, bucketName, fileName)
	if err != nil {
		return false, err
	}

	scanner := clamd.NewClamd(clamd.WithUnix("/tmp/clamd.sock"))
	dataReader := bytes.NewReader(fileData)

	safe, err := scanner.ScanStream(ctx, dataReader)
	if !safe && err == nil {
		err = moveFileToQuarantine(ctx, storageClient, bucketName, fileName, quarantineBucket)
		if err != nil {
			log.Printf("Failed to quarantine file: %v", err)
		}
	}

	return safe, err
}

func fetchFileFromBucket(ctx context.Context, client *storage.Client, bucketName, fileName string) ([]byte, error) {
	rc, err := client.Bucket(bucketName).Object(fileName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create object reader: %w", err)
	}
	defer rc.Close()

	return io.ReadAll(rc)
}

func moveFileToQuarantine(ctx context.Context, client *storage.Client, sourceBucket, sourceFile, quarantineBucket string) error {
	src := client.Bucket(sourceBucket).Object(sourceFile)
	dst := client.Bucket(quarantineBucket).Object(sourceBucket + "/" + sourceFile)

	_, err := dst.CopierFrom(src).Run(ctx)
	if err != nil {
		return fmt.Errorf("failed to copy object to quarantine: %w", err)
	}

	if err := src.Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete original object after quarantine: %w", err)
	}

	return nil
}
