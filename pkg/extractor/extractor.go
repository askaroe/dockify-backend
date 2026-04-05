package extractor

import (
	"bytes"
	"fmt"
	"os"
	"strings"

	"github.com/ledongthuc/pdf"
	"github.com/nguyenthenguyen/docx"
)

// ExtractText reads a file and returns its text content based on the content type.
func ExtractText(filePath, contentType string) (string, error) {
	switch {
	case contentType == "application/pdf":
		return extractPDF(filePath)
	case contentType == "application/vnd.openxmlformats-officedocument.wordprocessingml.document":
		return extractDOCX(filePath)
	case strings.HasPrefix(contentType, "text/"):
		return extractPlainText(filePath)
	default:
		// Fallback: try reading as plain text
		return extractPlainText(filePath)
	}
}

func extractPDF(filePath string) (string, error) {
	f, r, err := pdf.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("open pdf: %w", err)
	}
	defer f.Close()

	var buf bytes.Buffer
	totalPages := r.NumPage()
	for i := 1; i <= totalPages; i++ {
		p := r.Page(i)
		if p.V.IsNull() {
			continue
		}
		text, err := p.GetPlainText(nil)
		if err != nil {
			continue
		}
		buf.WriteString(text)
		buf.WriteString("\n")
	}

	return buf.String(), nil
}

func extractDOCX(filePath string) (string, error) {
	r, err := docx.ReadDocxFile(filePath)
	if err != nil {
		return "", fmt.Errorf("open docx: %w", err)
	}
	defer r.Close()

	doc := r.Editable()
	return doc.GetContent(), nil
}

func extractPlainText(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("read file: %w", err)
	}
	return string(data), nil
}
