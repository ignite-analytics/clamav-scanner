package update

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/exec"

	"cloud.google.com/go/storage"
	"github.com/ignite-analytics/clamav-scanner/internal/utils"
)

// Handle is the HTTP handler used to trigger ClamAV updates.
func Handle(mirrorBucket string, client *storage.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		_, _ = fmt.Fprintln(w, "Update request accepted. Processing in the background.")

		go func() {
			if err := update(); err != nil {
				log.Printf("Error updating CVDs: %v", err)
				return
			}

			if err := upload(context.Background(), client, mirrorBucket); err != nil {
				log.Printf("Error uploading CVDs to GCS: %v", err)
				return
			}

			log.Println("Update and upload completed successfully.")
		}()
	}
}

func update() error {
	log.Println("Updating CVDs using cvdtool...")
	cvdCmd := exec.Command("cvdupdate", "update", "-V", "-c", "/clamav/config.json")
	if err := cvdCmd.Run(); err != nil {
		log.Printf("Failed to update CVDs: %v", err)
		return err
	}
	return nil
}

func upload(ctx context.Context, client *storage.Client, mirrorBucket string) error {
	log.Println("Uploading CVDs to GCS...")

	if err := utils.SyncLocalToBucket(ctx, client, "/clamav/cvds", mirrorBucket); err != nil {
		log.Printf("Failed to upload CVDs to mirror: %v", err)
		return err
	}
	return nil
}
