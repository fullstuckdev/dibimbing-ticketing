package service

import (
	"encoding/json"
	"time"

	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/repository"
)

type AuditService interface {
	LogActivity(userID uint, action entity.AuditAction, entityType string, entityID uint, oldValue, newValue interface{}, ipAddress, userAgent string) error
	GetAuditLogs(page, limit int, userID uint, entityType string, startDate, endDate time.Time) ([]entity.AuditLog, int64, error)
	GetAuditLogsByEntity(entityType string, entityID uint) ([]entity.AuditLog, error)
}

type auditService struct {
	auditRepo repository.AuditRepository
}

func NewAuditService(auditRepo repository.AuditRepository) AuditService {
	return &auditService{
		auditRepo: auditRepo,
	}
}

func (s *auditService) LogActivity(userID uint, action entity.AuditAction, entityType string, entityID uint, oldValue, newValue interface{}, ipAddress, userAgent string) error {
	// Convert old and new values to JSON strings
	var oldValueStr, newValueStr string
	
	if oldValue != nil {
		oldValueJson, err := json.Marshal(oldValue)
		if err == nil {
			oldValueStr = string(oldValueJson)
		}
	}
	
	if newValue != nil {
		newValueJson, err := json.Marshal(newValue)
		if err == nil {
			newValueStr = string(newValueJson)
		}
	}
	
	// Create audit log entry
	auditLog := &entity.AuditLog{
		UserID:     userID,
		Action:     action,
		EntityType: entityType,
		EntityID:   entityID,
		OldValue:   oldValueStr,
		NewValue:   newValueStr,
		IPAddress:  ipAddress,
		UserAgent:  userAgent,
		CreatedAt:  time.Now(),
	}
	
	return s.auditRepo.CreateAuditLog(auditLog)
}

func (s *auditService) GetAuditLogs(page, limit int, userID uint, entityType string, startDate, endDate time.Time) ([]entity.AuditLog, int64, error) {
	return s.auditRepo.FindAuditLogs(page, limit, userID, entityType, startDate, endDate)
}

func (s *auditService) GetAuditLogsByEntity(entityType string, entityID uint) ([]entity.AuditLog, error) {
	return s.auditRepo.FindAuditLogsByEntityID(entityType, entityID)
} 