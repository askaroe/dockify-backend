package document

import (
	"context"

	"github.com/askaroe/dockify-backend/internal/models"
	"github.com/askaroe/dockify-backend/pkg/psql"
)

type Document interface {
	Create(ctx context.Context, doc models.Document) (string, error)
	GetByUserID(ctx context.Context, userID int) ([]models.Document, error)
	GetByID(ctx context.Context, id string) (*models.Document, error)
	Delete(ctx context.Context, id string) error
}

type document struct {
	db *psql.Client
}

func NewDocumentRepository(db *psql.Client) Document {
	return &document{db: db}
}

func (d *document) Create(ctx context.Context, doc models.Document) (string, error) {
	query := `INSERT INTO documents (user_id, file_name, file_path, file_size, content_type)
	           VALUES ($1, $2, $3, $4, $5) RETURNING id`
	var id string
	err := d.db.QueryRow(ctx, query, doc.UserId, doc.FileName, doc.FilePath, doc.FileSize, doc.ContentType).Scan(&id)
	if err != nil {
		return "", err
	}
	return id, nil
}

func (d *document) GetByUserID(ctx context.Context, userID int) ([]models.Document, error) {
	query := `SELECT id, user_id, file_name, file_path, file_size, content_type, uploaded_at
	           FROM documents WHERE user_id = $1 ORDER BY uploaded_at DESC`
	rows, err := d.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var docs []models.Document
	for rows.Next() {
		var doc models.Document
		err := rows.Scan(&doc.ID, &doc.UserId, &doc.FileName, &doc.FilePath, &doc.FileSize, &doc.ContentType, &doc.UploadedAt)
		if err != nil {
			return nil, err
		}
		docs = append(docs, doc)
	}
	return docs, rows.Err()
}

func (d *document) GetByID(ctx context.Context, id string) (*models.Document, error) {
	query := `SELECT id, user_id, file_name, file_path, file_size, content_type, uploaded_at
	           FROM documents WHERE id = $1`
	var doc models.Document
	err := d.db.QueryRow(ctx, query, id).Scan(&doc.ID, &doc.UserId, &doc.FileName, &doc.FilePath, &doc.FileSize, &doc.ContentType, &doc.UploadedAt)
	if err != nil {
		return nil, err
	}
	return &doc, nil
}

func (d *document) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM documents WHERE id = $1`
	_, err := d.db.Exec(ctx, query, id)
	return err
}
