package tests

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB() *gorm.DB {
	// Create a unique database name for each test to avoid conflicts
	dbName := fmt.Sprintf("file::memory:?cache=shared&key=%d", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	// Auto-migrate the schemas
	db.AutoMigrate(&entity.User{}, &entity.AuditLog{})

	return db
}

// auditRepositoryImpl is a test implementation of repository.AuditRepository
type auditRepositoryImpl struct {
	db *gorm.DB
}

func (r *auditRepositoryImpl) CreateAuditLog(auditLog *entity.AuditLog) error {
	return r.db.Create(auditLog).Error
}

func (r *auditRepositoryImpl) FindAuditLogs(page, limit int, userID uint, entityType string, startDate, endDate time.Time) ([]entity.AuditLog, int64, error) {
	var auditLogs []entity.AuditLog
	var count int64
	
	offset := (page - 1) * limit
	query := r.db
	
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}
	
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	} else if !endDate.IsZero() {
		query = query.Where("created_at <= ?", endDate)
	}
	
	if err := query.Model(&entity.AuditLog{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	
	if err := query.Offset(offset).Limit(limit).Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}
	
	return auditLogs, count, nil
}

func (r *auditRepositoryImpl) FindAuditLogsByEntityID(entityType string, entityID uint) ([]entity.AuditLog, error) {
	var auditLogs []entity.AuditLog
	
	err := r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).Find(&auditLogs).Error
	
	return auditLogs, err
}

func TestCreateAuditLog(t *testing.T) {
	// Setup
	db := setupTestDB()
	repo := &auditRepositoryImpl{
		db: db,
	}

	// Create test user
	user := &entity.User{
		Name:     "Test User",
		Email:    "test@example.com",
		Password: "password123",
		Role:     entity.RoleUser,
	}
	db.Create(user)

	// Create test audit log
	auditLog := &entity.AuditLog{
		UserID:     user.ID,
		Action:     entity.ActionCreate,
		EntityType: "event",
		EntityID:   1,
		CreatedAt:  time.Now(),
	}

	// Test
	err := repo.CreateAuditLog(auditLog)

	// Assertions
	assert.NoError(t, err)
	assert.NotZero(t, auditLog.ID)

	// Verify log was created
	var savedLog entity.AuditLog
	result := db.First(&savedLog, auditLog.ID)
	assert.NoError(t, result.Error)
	assert.Equal(t, auditLog.UserID, savedLog.UserID)
	assert.Equal(t, auditLog.Action, savedLog.Action)
	assert.Equal(t, auditLog.EntityType, savedLog.EntityType)
}

func TestFindAuditLogs(t *testing.T) {
	// Setup with a fresh database
	db := setupTestDB()
	
	// Clear any existing data to be safe
	db.Exec("DELETE FROM audit_logs")
	db.Exec("DELETE FROM users")
	
	repo := &auditRepositoryImpl{
		db: db,
	}

	// Create test user with a unique identifier in the email
	uniqueID := fmt.Sprintf("test_%d@example.com", time.Now().UnixNano())
	user := &entity.User{
		Name:     "Test User for Find Test",
		Email:    uniqueID,
		Password: "password123",
		Role:     entity.RoleUser,
	}
	result := db.Create(user)
	assert.NoError(t, result.Error)
	assert.NotZero(t, user.ID)

	// Create exactly 5 audit logs, ensuring they all have the test user's ID
	for i := 0; i < 5; i++ {
		auditLog := &entity.AuditLog{
			UserID:     user.ID,
			Action:     entity.ActionCreate,
			EntityType: "event",
			EntityID:   uint(i + 1),
			CreatedAt:  time.Now(),
		}
		createResult := db.Create(auditLog)
		assert.NoError(t, createResult.Error)
	}

	// Verify we have exactly 5 logs for our test user
	var countCheck int64
	db.Model(&entity.AuditLog{}).Where("user_id = ?", user.ID).Count(&countCheck)
	assert.Equal(t, int64(5), countCheck)

	// Test
	logs, count, err := repo.FindAuditLogs(1, 10, user.ID, "event", time.Time{}, time.Time{})

	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, int64(5), count)
	assert.Len(t, logs, 5)
} 