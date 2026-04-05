package document

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/askaroe/dockify-backend/internal/entity"
	"github.com/askaroe/dockify-backend/internal/services"
	"github.com/askaroe/dockify-backend/pkg/utils"
	"github.com/gin-gonic/gin"
)

type Document interface {
	UploadDocument(c *gin.Context)
	ListDocuments(c *gin.Context)
	DownloadDocument(c *gin.Context)
	DeleteDocument(c *gin.Context)
}

type document struct {
	s      *services.Service
	logger *utils.Logger
}

func NewDocumentHandler(s *services.Service, logger *utils.Logger) Document {
	return &document{s: s, logger: logger}
}

// UploadDocument godoc
// @Summary Upload a medical document
// @Description Upload a file (PDF, Word, image) for a user. Max 10 MB.
// @Tags Documents
// @Accept multipart/form-data
// @Produce json
// @Param user_id formData integer true "User ID"
// @Param file formData file true "Document file"
// @Success 201 {object} entity.DocumentResponse
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/documents/upload [post]
func (d *document) UploadDocument(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.PostForm("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "user_id is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid user_id"})
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "file is required"})
		return
	}

	// 10 MB limit
	if fileHeader.Size > 10*1024*1024 {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "file size exceeds 10 MB limit"})
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		d.logger.Errorf("UploadDocument open file error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to read uploaded file"})
		return
	}
	defer file.Close()

	doc, err := d.s.Document.Upload(ctx, userID, fileHeader.Filename, fileHeader.Header.Get("Content-Type"), fileHeader.Size, file)
	if err != nil {
		d.logger.Errorf("UploadDocument error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to upload document"})
		return
	}

	c.JSON(http.StatusCreated, entity.DocumentResponse{
		ID:          doc.ID,
		UserID:      doc.UserId,
		FileName:    doc.FileName,
		FileSize:    doc.FileSize,
		ContentType: doc.ContentType,
		Summary:     doc.Summary,
		UploadedAt:  "just now",
	})
}

// ListDocuments godoc
// @Summary List user documents
// @Description Returns all documents uploaded by a user
// @Tags Documents
// @Accept json
// @Produce json
// @Param user_id query integer true "User ID"
// @Success 200 {array} entity.DocumentResponse
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/documents [get]
func (d *document) ListDocuments(c *gin.Context) {
	ctx := c.Request.Context()

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "user_id is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid user_id"})
		return
	}

	docs, err := d.s.Document.ListByUser(ctx, userID)
	if err != nil {
		d.logger.Errorf("ListDocuments error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to list documents"})
		return
	}

	var resp []entity.DocumentResponse
	for _, doc := range docs {
		uploadedAt := ""
		if doc.UploadedAt != nil {
			uploadedAt = doc.UploadedAt.Format("2006-01-02 15:04:05")
		}
		resp = append(resp, entity.DocumentResponse{
			ID:          doc.ID,
			UserID:      doc.UserId,
			FileName:    doc.FileName,
			FileSize:    doc.FileSize,
			ContentType: doc.ContentType,
			Summary:     doc.Summary,
			UploadedAt:  uploadedAt,
		})
	}

	if resp == nil {
		resp = []entity.DocumentResponse{}
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteDocument godoc
// @Summary Delete a document
// @Description Delete a document by ID. Validates user ownership.
// @Tags Documents
// @Produce json
// @Param id path string true "Document UUID"
// @Param user_id query integer true "User ID"
// @Success 200 {object} entity.DocumentDeleteResponse
// @Failure 400 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/documents/{id} [delete]
func (d *document) DeleteDocument(c *gin.Context) {
	ctx := c.Request.Context()

	docID := c.Param("id")
	if docID == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "document id is required"})
		return
	}

	userIDStr := c.Query("user_id")
	if userIDStr == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "user_id is required"})
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "invalid user_id"})
		return
	}

	err = d.s.Document.Delete(ctx, docID, userID)
	if err != nil {
		d.logger.Errorf("DeleteDocument error: %v", err)
		c.JSON(http.StatusInternalServerError, entity.ErrorMessage{Message: "failed to delete document"})
		return
	}

	c.JSON(http.StatusOK, entity.DocumentDeleteResponse{Message: "document deleted"})
}

// DownloadDocument godoc
// @Summary Download a document
// @Description Download a file by document UUID
// @Tags Documents
// @Produce octet-stream
// @Param id path string true "Document UUID"
// @Success 200 {file} binary "file content"
// @Failure 400 {object} entity.ErrorMessage
// @Failure 404 {object} entity.ErrorMessage
// @Failure 500 {object} entity.ErrorMessage
// @Router /api/v1/documents/{id}/download [get]
func (d *document) DownloadDocument(c *gin.Context) {
	ctx := c.Request.Context()

	docID := c.Param("id")
	if docID == "" {
		c.JSON(http.StatusBadRequest, entity.ErrorMessage{Message: "document id is required"})
		return
	}

	doc, err := d.s.Document.GetByID(ctx, docID)
	if err != nil {
		d.logger.Errorf("DownloadDocument error: %v", err)
		c.JSON(http.StatusNotFound, entity.ErrorMessage{Message: "document not found"})
		return
	}

	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, doc.FileName))
	c.Header("Content-Type", doc.ContentType)
	c.File(doc.FilePath)
}
