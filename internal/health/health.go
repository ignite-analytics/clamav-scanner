package health

import (
	"log"
	"net/http"

	"github.com/shirou/gopsutil/process"
)

var isProcessRunning = IsProcessRunning

// Handle is the HTTP handler used to check the status of ClamAV daemon.
func Handle() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		processName := "clamd"
		running, err := isProcessRunning(processName)
		if err != nil {
			log.Printf("Failed to get ClamAV daemon: %v", err)
			http.Error(w, "Failed to get ClamAV daemon", http.StatusInternalServerError)
			return
		}

		if running {
			w.WriteHeader(http.StatusOK)
		} else {
			http.Error(w, "ClamAV daemon is not running", http.StatusServiceUnavailable)
			log.Printf("ClamAV daemon is not running")
		}
	}
}

func IsProcessRunning(name string) (bool, error) {
	processes, err := process.Processes()
	if err != nil {
		return false, err
	}

	for _, p := range processes {
		pName, err := p.Name()
		if err == nil && pName == name {
			return true, nil
		}
	}
	return false, nil
}
