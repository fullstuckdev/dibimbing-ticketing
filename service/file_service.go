package service

import (
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileService handles file upload and download operations
type FileService interface {
	UploadFile(file *multipart.FileHeader, fileType string) (string, error)
	GetFilePath(filename string, fileType string) (string, error)
	DeleteFile(filename string, fileType string) error
}

type fileService struct {
	baseStoragePath string
}

func NewFileService() FileService {
	// Create uploads directory if it doesn't exist
	uploadDir := "uploads"
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, 0755)
	}

	// Create subdirectories for different file types
	for _, dir := range []string{"events", "tickets", "profiles"} {
		fullPath := filepath.Join(uploadDir, dir)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			os.MkdirAll(fullPath, 0755)
		}
	}

	return &fileService{
		baseStoragePath: uploadDir,
	}
}

func (s *fileService) UploadFile(file *multipart.FileHeader, fileType string) (string, error) {
	// Validate file type
	if !isValidFileType(fileType) {
		return "", errors.New("invalid file type")
	}

	// Generate unique file name to prevent collisions
	extension := filepath.Ext(file.Filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), extension)
	
	// Create the full path for storing the file
	storagePath := filepath.Join(s.baseStoragePath, fileType)
	filePath := filepath.Join(storagePath, newFilename)
	
	// Open the uploaded file
	src, err := file.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()
	
	// Create the destination file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer dst.Close()
	
	// Copy the uploaded file to the destination file
	if _, err = io.Copy(dst, src); err != nil {
		return "", err
	}
	
	return newFilename, nil
}

func (s *fileService) GetFilePath(filename string, fileType string) (string, error) {
	// Validate file type
	if !isValidFileType(fileType) {
		return "", errors.New("invalid file type")
	}
	
	// Prevent directory traversal attacks
	filename = filepath.Base(filename)
	
	// Create the full path
	fullPath := filepath.Join(s.baseStoragePath, fileType, filename)
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return "", errors.New("file not found")
	}
	
	return fullPath, nil
}

func (s *fileService) DeleteFile(filename string, fileType string) error {
	// Validate file type
	if !isValidFileType(fileType) {
		return errors.New("invalid file type")
	}
	
	// Prevent directory traversal attacks
	filename = filepath.Base(filename)
	
	// Create the full path
	fullPath := filepath.Join(s.baseStoragePath, fileType, filename)
	
	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return errors.New("file not found")
	}
	
	// Delete the file
	return os.Remove(fullPath)
}

// isValidFileType checks if the file type is allowed
func isValidFileType(fileType string) bool {
	validTypes := []string{"events", "tickets", "profiles"}
	for _, t := range validTypes {
		if t == fileType {
			return true
		}
	}
	return false
}

// Helper function to check if a file is an image
func isImageFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

// Helper function to check if a file is a PDF
func isPdfFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".pdf"
} 