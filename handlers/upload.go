package handlers

import (
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

type FileUploadResponse struct {
	Status   string `json:"status"`
	Message  string `json:"message,omitempty"`
	FileName string `json:"fileName,omitempty"`
	FileSize int64  `json:"fileSize,omitempty"`
	FileURL  string `json:"fileUrl,omitempty"`
	Error    string `json:"error,omitempty"`
}

// type UploadHandler struct {
// 	db  *gorm.DB
// 	cfg *config.Config
// }

// func NewUploadHandler(db *gorm.DB, cfg *config.Config) *UploadHandler {
// 	return &UploadHandler{db: db, cfg: cfg}
// }

func (h *ApiHandler) UploadFile(c *gin.Context) {
	// Get file
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(200, FileUploadResponse{
			Status:  "error",
			Message: "No file provided",
		})
		return
	}

	// Validate file size
	if file.Size > 10*1024*1024 { // 5MB
		c.JSON(200, FileUploadResponse{
			Status:  "error",
			Message: "File size exceeds 5MB limit",
		})
		return
	}

	// Create directory
	uploadDir := filepath.Join(h.cfg.Cwd, "/uploads")
	os.MkdirAll(uploadDir, 0755)

	// Save file
	filepath := filepath.Join(uploadDir, file.Filename)
	if err := c.SaveUploadedFile(file, filepath); err != nil {
		c.JSON(200, FileUploadResponse{
			Status:  "error",
			Message: "Failed to save file",
		})
		return
	}

	c.JSON(200, FileUploadResponse{
		Status:   "success",
		Message:  "File uploaded successfully",
		FileName: file.Filename,
		FileSize: file.Size,
		FileURL:  "/uploads/" + file.Filename,
	})
}
