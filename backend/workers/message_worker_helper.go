package workers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

// fetchMediaData fetches media data from URL or local file
func (w *MessageWorker) fetchMediaData(urlOrPath string) ([]byte, error) {
	// Check if it's a local file
	if !strings.HasPrefix(urlOrPath, "http") {
		// Assume local file
		return os.ReadFile(urlOrPath)
	}

	// Fetch from URL
	resp, err := http.Get(urlOrPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch media: status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}
