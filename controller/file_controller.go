package controller

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/taufikmulyawan/ticketing-system/service"
)

type FileController interface {
	UploadFile(c *gin.Context)
	DownloadFile(c *gin.Context)
	DeleteFile(c *gin.Context)
}

type fileController struct {
	fileService service.FileService
}

func NewFileController(fileService service.FileService) FileController {
	return &fileController{
		fileService: fileService,
	}
}

// UploadFile godoc
// @Summary Upload a file
// @Description Upload a file to the server (event image, ticket file, profile picture)
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to upload"
// @Param type formData string true "File type (events, tickets, profiles)"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Router /files/upload [post]
func (ctrl *fileController) UploadFile(c *gin.Context) {
	// Get file from form
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file provided"})
		return
	}
	
	// Get file type from form
	fileType := c.PostForm("type")
	if fileType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File type is required"})
		return
	}
	
	// Upload file
	filename, err := ctrl.fileService.UploadFile(file, fileType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"filename": filename,
		"type": fileType,
	})
}

// DownloadFile godoc
// @Summary Download a file
// @Description Download a file from the server
// @Tags files
// @Produce octet-stream
// @Param filename path string true "File name"
// @Param type path string true "File type (events, tickets, profiles)"
// @Success 200 {file} binary
// @Failure 400,404 {object} map[string]string
// @Router /files/{type}/{filename} [get]
func (ctrl *fileController) DownloadFile(c *gin.Context) {
	// Get filename and file type from URL
	filename := c.Param("filename")
	fileType := c.Param("type")
	
	// Get file path
	filePath, err := ctrl.fileService.GetFilePath(filename, fileType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	// Set content disposition header to force download
	c.Header("Content-Disposition", "attachment; filename="+filename)
	c.Header("Content-Type", getContentType(filename))
	
	// Serve the file
	c.File(filePath)
}

// DeleteFile godoc
// @Summary Delete a file
// @Description Delete a file from the server
// @Tags files
// @Produce json
// @Param filename path string true "File name"
// @Param type path string true "File type (events, tickets, profiles)"
// @Security BearerAuth
// @Success 200 {object} map[string]string
// @Failure 400,404 {object} map[string]string
// @Router /files/{type}/{filename} [delete]
func (ctrl *fileController) DeleteFile(c *gin.Context) {
	// Get filename and file type from URL
	filename := c.Param("filename")
	fileType := c.Param("type")
	
	// Delete file
	err := ctrl.fileService.DeleteFile(filename, fileType)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}

// Helper function to determine content type based on file extension
func getContentType(filename string) string {
	ext := filepath.Ext(filename)
	switch ext {
	case ".pdf":
		return "application/pdf"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	default:
		return "application/octet-stream"
	}
} 