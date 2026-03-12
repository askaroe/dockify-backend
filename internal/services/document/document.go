package document

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/google/uuid"

	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/internal/repository"
)

type Document interface {
	Upload(ctx context.Context, userID int, fileName string, contentType string, fileSize int64, fileReader io.Reader) (*models.Document, error)
	ListByUser(ctx context.Context, userID int) ([]models.Document, error)
	Delete(ctx context.Context, docID string, userID int) error
}

type document struct {
	repo    *repository.Repository
	baseDir string
}

func NewDocumentService(repo *repository.Repository) Document {
	return &document{repo: repo, baseDir: "documents"}
}

func (d *document) Upload(ctx context.Context, userID int, fileName string, contentType string, fileSize int64, fileReader io.Reader) (*models.Document, error) {
	// Ensure directory exists
	userDir := filepath.Join(d.baseDir, fmt.Sprintf("%d", userID))
	if err := os.MkdirAll(userDir, 0755); err != nil {
		return nil, fmt.Errorf("create upload directory: %w", err)
	}

	// Generate unique file name
	uniqueName := fmt.Sprintf("%s_%s", uuid.New().String(), fileName)
	filePath := filepath.Join(userDir, uniqueName)

	// Write file to disk
	dst, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("create file: %w", err)
	}
	defer dst.Close()

	if _, err := io.Copy(dst, fileReader); err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("write file: %w", err)
	}

	// Insert into database
	doc := models.Document{
		UserId:      userID,
		FileName:    fileName,
		FilePath:    filePath,
		FileSize:    fileSize,
		ContentType: contentType,
	}

	id, err := d.repo.Document.Create(ctx, doc)
	if err != nil {
		os.Remove(filePath)
		return nil, fmt.Errorf("save document metadata: %w", err)
	}

	doc.ID = id
	return &doc, nil
}

func (d *document) ListByUser(ctx context.Context, userID int) ([]models.Document, error) {
	return d.repo.Document.GetByUserID(ctx, userID)
}

func (d *document) Delete(ctx context.Context, docID string, userID int) error {
	// Verify the document belongs to the user
	doc, err := d.repo.Document.GetByID(ctx, docID)
	if err != nil {
		return fmt.Errorf("document not found: %w", err)
	}
	if doc.UserId != userID {
		return fmt.Errorf("unauthorized: document does not belong to user")
	}

	// Delete from disk
	if err := os.Remove(doc.FilePath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("remove file: %w", err)
	}

	// Delete from database
	return d.repo.Document.Delete(ctx, docID)
}
