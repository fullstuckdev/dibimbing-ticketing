package tests

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/taufikmulyawan/ticketing-system/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB creates an in-memory SQLite database for testing
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the schemas
	db.AutoMigrate(&entity.User{}, &entity.AuditLog{})

	return db
}

// GetProjectRoot returns the absolute path to the project root directory
func GetProjectRoot() string {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(filepath.Dir(b))
}

// ClearTestDirs removes temporary test directories
func ClearTestDirs(paths ...string) {
	for _, path := range paths {
		os.RemoveAll(path)
	}
} 