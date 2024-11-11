package utils

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// SyncLocalToBucket uploads files from a local directory to the specified GCS bucket.
func SyncLocalToBucket(ctx context.Context, client *storage.Client, localDir, bucketName string) error {
	dirEntries, err := os.ReadDir(localDir)
	if err != nil {
		return fmt.Errorf("failed to read local directory %s: %w", localDir, err)
	}

	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}

		localPath := filepath.Clean(fmt.Sprintf("%s/%s", localDir, entry.Name()))
		object := client.Bucket(bucketName).Object(entry.Name())
		writer := object.NewWriter(ctx)

		file, err := os.Open(localPath)
		if err != nil {
			return fmt.Errorf("failed to open local file %s: %w", localPath, err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				fmt.Printf("Failed to close file %s: %v\n", localPath, err)
			}
		}()

		defer func() {
			if err := writer.Close(); err != nil {
				fmt.Printf("Failed to close writer for object %s: %v\n", entry.Name(), err)
			}
		}()

		if _, err := io.Copy(writer, file); err != nil {
			return fmt.Errorf("failed to copy file %s to GCS object %s: %w", localPath, entry.Name(), err)
		}

		fmt.Printf("Successfully uploaded %s to %s\n", localPath, entry.Name())
	}

	return nil
}

// SyncBucketToLocal downloads files from a GCS bucket to a local directory.
func SyncBucketToLocal(ctx context.Context, client *storage.Client, bucketName, localDir string) error {
	it := client.Bucket(bucketName).Objects(ctx, nil)
	for {
		attr, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("failed to iterate objects in bucket: %w", err)
		}

		localPath := filepath.Clean(fmt.Sprintf("%s/%s", localDir, attr.Name))
		object := client.Bucket(bucketName).Object(attr.Name)
		reader, err := object.NewReader(ctx)
		if err != nil {
			return fmt.Errorf("failed to open reader for GCS object %s: %w", attr.Name, err)
		}

		defer func() {
			if err := reader.Close(); err != nil {
				fmt.Printf("Failed to close file %s: %v\n", localPath, err)
			}
		}()

		file, err := os.Create(localPath)
		if err != nil {
			return fmt.Errorf("failed to create local file %s: %w", localPath, err)
		}

		defer func() {
			if err := file.Close(); err != nil {
				fmt.Printf("Failed to close file %s: %v\n", localPath, err)
			}
		}()

		if _, err := io.Copy(file, reader); err != nil {
			return fmt.Errorf("failed to copy GCS object %s to local file: %w", attr.Name, err)
		}
	}

	return nil
}
