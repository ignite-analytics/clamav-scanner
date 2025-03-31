package scan

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"cloud.google.com/go/pubsub"
	"cloud.google.com/go/storage"
	"github.com/lyimmi/go-clamd"
)

type requestBody struct {
	Kind   string `json:"kind"`
	Name   string `json:"name"`
	Bucket string `json:"bucket"`
}

var (
	performScanFunc = performScan
	topicName       = os.Getenv("PUBSUB_TOPIC")
	projectId       = os.Getenv("PROJECT_ID")
)

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
		defer func() {
			if err := r.Body.Close(); err != nil {
				log.Printf("Failed to close request body: %v", err)
			}
		}()

		log.Printf("Received scan request for bucket: %s, file: %s\n", reqBody.Bucket, reqBody.Name)

		// Acknowledge the Pub/Sub message before it's proccessed
		if pubsubMsg, ok := r.Context().Value("pubsubMessage").(*pubsub.Message); ok {
			pubsubMsg.Ack()
		}

		safe, err := performScanFunc(ctx, reqBody.Bucket, reqBody.Name, quarantineBucket)
		if err != nil {
			log.Printf("Error scanning file: %v", err)
			http.Error(w, "Failed to scan file", http.StatusInternalServerError)
			return
		}

		resultMsg := "NO_THREAT"
		if !safe {
			resultMsg = "THREAT"
		}

		log.Printf("Scan completed for %s: %s\n", reqBody.Name, resultMsg)

		// Respond to the HTTP request with the scan result
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(map[string]string{
			"message": "Scan completed for " + reqBody.Name,
			"result":  resultMsg,
		}); err != nil {
			log.Printf("Error encoding JSON response: %v", err)
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}

		go publishScanResultToPubSub(context.Background(), reqBody.Bucket, reqBody.Name, resultMsg)
	}
}

func performScan(ctx context.Context, bucketName, fileName, quarantineBucket string) (bool, error) {
	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return false, fmt.Errorf("failed to create storage client: %w", err)
	}
	defer func() {
		if err := storageClient.Close(); err != nil {
			log.Printf("Failed to close storage client: %v", err)
		}
	}()

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

func publishScanResultToPubSub(ctx context.Context, bucketName, fileName, resultMsg string) {
	pubsubClient, err := pubsub.NewClient(ctx, projectId)
	if err != nil {
		log.Printf("Failed to create Pub/Sub client: %v", err)
		return
	}
	defer func() {
		if err := pubsubClient.Close(); err != nil {
			log.Printf("Failed to close Pub/Sub client: %v", err)
		}
	}()

	topic := pubsubClient.Topic(topicName)
	resultMessage := map[string]string{
		"file":   fileName,
		"bucket": bucketName,
		"result": resultMsg,
	}
	messageData, err := json.Marshal(resultMessage)
	if err != nil {
		log.Printf("Failed to marshal result message: %v", err)
		return
	}

	_, err = topic.Publish(ctx, &pubsub.Message{
		Data: messageData,
		Attributes: map[string]string{
			"bucketName": bucketName,
		},
	}).Get(ctx)
	if err != nil {
		log.Printf("Failed to publish result to Pub/Sub: %v", err)
	}
}

func fetchFileFromBucket(ctx context.Context, client *storage.Client, bucketName, fileName string) ([]byte, error) {
	rc, err := client.Bucket(bucketName).Object(fileName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create object reader: %w", err)
	}
	defer func() {
		if err := rc.Close(); err != nil {
			log.Printf("Failed to close object reader: %v", err)
		}
	}()

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
