package tests

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/service"
)

// Mock AuditRepository
type MockAuditRepository struct {
	mock.Mock
}

func (m *MockAuditRepository) CreateAuditLog(auditLog *entity.AuditLog) error {
	args := m.Called(auditLog)
	return args.Error(0)
}

func (m *MockAuditRepository) FindAuditLogs(page, limit int, userID uint, entityType string, startDate, endDate time.Time) ([]entity.AuditLog, int64, error) {
	args := m.Called(page, limit, userID, entityType, startDate, endDate)
	return args.Get(0).([]entity.AuditLog), args.Get(1).(int64), args.Error(2)
}

func (m *MockAuditRepository) FindAuditLogsByEntityID(entityType string, entityID uint) ([]entity.AuditLog, error) {
	args := m.Called(entityType, entityID)
	return args.Get(0).([]entity.AuditLog), args.Error(1)
}

func TestLogActivity_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockAuditRepository)
	auditService := service.NewAuditService(mockRepo)
	
	// Set expectations
	mockRepo.On("CreateAuditLog", mock.AnythingOfType("*entity.AuditLog")).Return(nil)
	
	// Test data
	userID := uint(1)
	action := entity.ActionCreate
	entityType := "event"
	entityID := uint(10)
	var oldValue interface{} = nil
	newValue := map[string]interface{}{"name": "New Event", "date": "2023-12-31"}
	ipAddress := "127.0.0.1"
	userAgent := "Mozilla/5.0"
	
	// Execute test
	err := auditService.LogActivity(userID, action, entityType, entityID, oldValue, newValue, ipAddress, userAgent)
	
	// Assertions
	assert.NoError(t, err)
	mockRepo.AssertExpectations(t)
}

func TestGetAuditLogs_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockAuditRepository)
	auditService := service.NewAuditService(mockRepo)
	
	// Test data
	page := 1
	limit := 10
	userID := uint(1)
	entityType := "event"
	startDate := time.Now().AddDate(0, 0, -7) // Last week
	endDate := time.Now()
	
	expectedLogs := []entity.AuditLog{
		{
			ID:         1,
			UserID:     1,
			Action:     entity.ActionCreate,
			EntityType: "event",
			EntityID:   10,
			CreatedAt:  time.Now(),
		},
	}
	expectedCount := int64(1)
	
	// Set expectations
	mockRepo.On("FindAuditLogs", page, limit, userID, entityType, startDate, endDate).
		Return(expectedLogs, expectedCount, nil)
	
	// Execute test
	logs, count, err := auditService.GetAuditLogs(page, limit, userID, entityType, startDate, endDate)
	
	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedCount, count)
	assert.Equal(t, expectedLogs, logs)
	mockRepo.AssertExpectations(t)
}

func TestGetAuditLogsByEntity_Success(t *testing.T) {
	// Setup
	mockRepo := new(MockAuditRepository)
	auditService := service.NewAuditService(mockRepo)
	
	// Test data
	entityType := "event"
	entityID := uint(10)
	
	expectedLogs := []entity.AuditLog{
		{
			ID:         1,
			UserID:     1,
			Action:     entity.ActionCreate,
			EntityType: "event",
			EntityID:   10,
			CreatedAt:  time.Now(),
		},
		{
			ID:         2,
			UserID:     2,
			Action:     entity.ActionUpdate,
			EntityType: "event",
			EntityID:   10,
			CreatedAt:  time.Now(),
		},
	}
	
	// Set expectations
	mockRepo.On("FindAuditLogsByEntityID", entityType, entityID).Return(expectedLogs, nil)
	
	// Execute test
	logs, err := auditService.GetAuditLogsByEntity(entityType, entityID)
	
	// Assertions
	assert.NoError(t, err)
	assert.Equal(t, expectedLogs, logs)
	mockRepo.AssertExpectations(t)
} 