package scheduler

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/askaroe/dockify-backend/pkg/psql"
	"github.com/askaroe/dockify-backend/pkg/utils"
)

// StartDocumentCleanup runs a background goroutine that clears the documents
// directory and truncates the documents table every 3 hours.
func StartDocumentCleanup(db *psql.Client, logger *utils.Logger) {
	ticker := time.NewTicker(3 * time.Minute)

	go func() {
		for range ticker.C {
			logger.Infof("Scheduler: starting document cleanup...")
			if err := cleanup(db); err != nil {
				logger.Errorf("Scheduler: cleanup failed: %v", err)
			} else {
				logger.Infof("Scheduler: document cleanup complete")
			}
		}
	}()
}

func cleanup(db *psql.Client) error {
	// Remove files from disk
	dir := "documents"
	if err := os.RemoveAll(dir); err != nil {
		return fmt.Errorf("remove documents directory: %w", err)
	}
	// Recreate empty directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("recreate documents directory: %w", err)
	}

	// Truncate the documents table
	_, err := db.Exec(context.Background(), "DELETE FROM documents")
	if err != nil {
		return fmt.Errorf("truncate documents table: %w", err)
	}

	return nil
}
