package tests

import (
	"bytes"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taufikmulyawan/ticketing-system/controller"
)

// MockFileService implements service.FileService
type MockFileService struct {
	mock.Mock
}

func (m *MockFileService) UploadFile(file *multipart.FileHeader, fileType string) (string, error) {
	args := m.Called(file, fileType)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) GetFilePath(filename string, fileType string) (string, error) {
	args := m.Called(filename, fileType)
	return args.String(0), args.Error(1)
}

func (m *MockFileService) DeleteFile(filename string, fileType string) error {
	args := m.Called(filename, fileType)
	return args.Error(0)
}

func TestUploadFile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockFileService)
	fileController := controller.NewFileController(mockService)

	// Create a test router
	router := gin.Default()
	router.POST("/files/upload", fileController.UploadFile)

	// Create a test request
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", "test.txt")
	part.Write([]byte("test content"))
	writer.WriteField("type", "events")
	writer.Close()

	// Set expectations
	mockService.On("UploadFile", mock.Anything, "events").Return("test-12345.txt", nil)

	// Create test request
	req, _ := http.NewRequest("POST", "/files/upload", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
}

func TestDownloadFile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockFileService)
	fileController := controller.NewFileController(mockService)

	// Create a test router
	router := gin.Default()
	router.GET("/files/:type/:filename", fileController.DownloadFile)

	// Create a temporary test file
	tempDir, err := os.MkdirTemp("", "test-files")
	assert.NoError(t, err)
	defer os.RemoveAll(tempDir)
	
	testFilePath := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFilePath, []byte("test content"), 0644)
	assert.NoError(t, err)

	// Set expectations - return the path to our test file
	mockService.On("GetFilePath", "test.txt", "events").Return(testFilePath, nil)

	// Create test request
	req, _ := http.NewRequest("GET", "/files/events/test.txt", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)
	assert.Equal(t, "test content", resp.Body.String())
	mockService.AssertExpectations(t)
}

func TestDeleteFile(t *testing.T) {
	// Setup
	gin.SetMode(gin.TestMode)
	mockService := new(MockFileService)
	fileController := controller.NewFileController(mockService)

	// Create a test router
	router := gin.Default()
	router.DELETE("/files/:type/:filename", fileController.DeleteFile)

	// Set expectations
	mockService.On("DeleteFile", "test.txt", "events").Return(nil)

	// Create test request
	req, _ := http.NewRequest("DELETE", "/files/events/test.txt", nil)
	resp := httptest.NewRecorder()

	// Perform the request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)
	mockService.AssertExpectations(t)
} 