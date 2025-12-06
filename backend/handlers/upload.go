package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

const (
	uploadsDir   = "/app/data/uploads"
	maxFileSize  = 10 << 20 // 10MB
)

// UploadResponse represents the response from file upload
type UploadResponse struct {
	Success  bool   `json:"success"`
	FileURL  string `json:"file_url"`
	FileName string `json:"file_name"`
	FileType string `json:"file_type"`
	FileSize int64  `json:"file_size"`
}

// UploadFile handles file upload
func UploadFile(c echo.Context) error {
	// Get tenant ID from context
	tenantID := getTenantIDFromContext(c)
	if tenantID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Tenant not found",
		})
	}

	// Get uploaded file
	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File is required",
		})
	}

	// Check file size
	if file.Size > maxFileSize {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "File size exceeds maximum allowed (10MB)",
		})
	}

	// Open uploaded file
	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to open uploaded file",
		})
	}
	defer src.Close()

	// Determine file type
	ext := strings.ToLower(filepath.Ext(file.Filename))
	fileType := determineFileType(ext)

	// Create unique filename
	uniqueID := uuid.New().String()
	timestamp := time.Now().Format("20060102")
	newFileName := fmt.Sprintf("%s_%s_%s%s", tenantID, timestamp, uniqueID[:8], ext)

	// Create tenant upload directory
	tenantDir := filepath.Join(uploadsDir, tenantID)
	if err := os.MkdirAll(tenantDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create upload directory",
		})
	}

	// Create destination file
	dstPath := filepath.Join(tenantDir, newFileName)
	dst, err := os.Create(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to create destination file",
		})
	}
	defer dst.Close()

	// Copy file contents
	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Failed to save file",
		})
	}

	// Generate public URL
	fileURL := fmt.Sprintf("/uploads/%s/%s", tenantID, newFileName)

	return c.JSON(http.StatusOK, UploadResponse{
		Success:  true,
		FileURL:  fileURL,
		FileName: file.Filename,
		FileType: fileType,
		FileSize: file.Size,
	})
}

// determineFileType determines the type of file based on extension
func determineFileType(ext string) string {
	switch ext {
	case ".jpg", ".jpeg", ".png", ".gif", ".webp":
		return "image"
	case ".pdf":
		return "document"
	case ".doc", ".docx":
		return "document"
	case ".xls", ".xlsx":
		return "document"
	case ".mp4", ".mov", ".avi":
		return "video"
	case ".mp3", ".wav", ".ogg":
		return "audio"
	default:
		return "file"
	}
}
