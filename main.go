package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"syscall"
	"time"

	"cloud.google.com/go/storage"
	"github.com/ignite-analytics/clamav-scanner/internal/health"
	"github.com/ignite-analytics/clamav-scanner/internal/mirror"
	"github.com/ignite-analytics/clamav-scanner/internal/scan"
	"github.com/ignite-analytics/clamav-scanner/internal/update"
	"golang.org/x/net/context"
)

var (
	mirrorBucket     = os.Getenv("MIRROR_BUCKET")
	quarantineBucket = os.Getenv("QUARANTINE_BUCKET")
	listenAddress    = os.Getenv("LISTEN_ADDRESS")
	stop             = make(chan os.Signal, 1)
)

func main() {
	ctx := context.Background()
	client, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("Failed to create Google Cloud Storage client: %v", err)
	}

	// Bootstrap the ClamAV configurations
	bootstrap(ctx, client)

	// Reload the clamav and freshclam services
	reloadServices()

	// Create new HTTP server
	server := &http.Server{
		Addr: listenAddress,
	}

	setupHandlers()

	// Setup channel for listening to signals and server errors
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Scanner service listening on %s...\n", server.Addr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to listen on %s: %v", server.Addr, err)
		}
	}()

	// Wait for interrupt or external shutdown signal
	<-stop

	// Start graceful shutdown
	log.Println("Shutting down the server...")
	timeoutCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer shutdownCancel()
	if err := server.Shutdown(timeoutCtx); err != nil {
		log.Fatalf("Server shutdown failed: %v", err)
	}

	log.Println("Server gracefully stopped.")
}

func bootstrap(ctx context.Context, client *storage.Client) {
	if mirrorBucket == "" {
		log.Fatal("Please specify ClamAV mirror bucket")
	}

	// Check if the bucket contains any objects
	itr := client.Bucket(mirrorBucket).Objects(ctx, nil)
	if _, err := itr.Next(); err != nil {
		log.Printf("Mirror bucket is empty, initializing mirror in %s...\n", mirrorBucket)

		cvdCmd := exec.Command("cvdupdate", "update", "-V", "-c", "/clamav/config.json")
		if err := cvdCmd.Run(); err != nil {
			log.Printf("Failed to update CVDs: %v", err)
		}

		gsutilCmd := exec.Command("gsutil", "-m", "-q", "rsync", "-d", "-c", "-r", "/clamav/cvds",
			fmt.Sprintf("gs://%s/", mirrorBucket))
		if err := gsutilCmd.Run(); err != nil {
			log.Fatalf("Failed to upload inital CVDs to mirror: %v", err)
		}
	}

	log.Println("Downloading and updating CVDs from mirror...")

	cmd := exec.Command("gsutil", "-m", "rsync", "-c", "-r",
		fmt.Sprintf("gs://%s/", mirrorBucket), "/var/lib/clamav")
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to synchronize CVDs: %v", err)
	}
}

func reloadServices() {
	log.Println("Starting ClamAV services...")

	if err := exec.Command("service", "clamav-freshclam", "force-reload").Run(); err != nil {
		log.Printf("Failed to start clamav-freshclam: %v", err)
	}

	if err := exec.Command("service", "clamav-daemon", "force-reload").Run(); err != nil {
		log.Printf("Failed to start clamav-daemon: %v", err)
	}
}

func setupHandlers() {
	http.HandleFunc("/mirror/", mirror.Handle(mirrorBucket))
	http.HandleFunc("/update", update.Handle(mirrorBucket))
	http.HandleFunc("/scan", scan.Handle(quarantineBucket))
	http.HandleFunc("/health", health.Handle())
}
