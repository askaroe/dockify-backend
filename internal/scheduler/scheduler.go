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
	// Remove contents inside the documents directory (not the directory itself)
	dir := "documents"
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // nothing to clean
		}
		return fmt.Errorf("read documents directory: %w", err)
	}
	for _, entry := range entries {
		path := fmt.Sprintf("%s/%s", dir, entry.Name())
		if err := os.RemoveAll(path); err != nil {
			return fmt.Errorf("remove %s: %w", path, err)
		}
	}

	// Truncate the documents table
	_, err = db.Exec(context.Background(), "DELETE FROM documents")
	if err != nil {
		return fmt.Errorf("truncate documents table: %w", err)
	}

	return nil
}
