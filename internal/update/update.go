package update

import (
	"fmt"
	"log"
	"net/http"
	"os/exec"
)

// Handle is the HTTP handler used to trigger ClamAV updates.
func Handle(mirrorBucket string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		if err := update(); err != nil {
			http.Error(w, "Failed to update CVDs", http.StatusInternalServerError)
			return
		}

		if err := upload(mirrorBucket); err != nil {
			http.Error(w, "Failed to upload CVDs", http.StatusInternalServerError)
			return
		}

		fmt.Fprintln(w, "Update completed.")
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

func upload(mirrorBucket string) error {
	log.Println("Uploading CVDs to GCS using gsutil...")
	gsutilCmd := exec.Command("gsutil", "-m", "-q", "rsync", "-d", "-c", "-r", "/clamav/cvds",
		fmt.Sprintf("gs://%s/", mirrorBucket))
	if err := gsutilCmd.Run(); err != nil {
		log.Printf("Failed to upload CVDs to mirror: %v", err)
		return err
	}
	return nil
}
