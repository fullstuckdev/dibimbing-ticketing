package repository

import (
	"time"

	"github.com/taufikmulyawan/ticketing-system/config"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"gorm.io/gorm"
)

type AuditRepository interface {
	CreateAuditLog(auditLog *entity.AuditLog) error
	FindAuditLogs(page, limit int, userID uint, entityType string, startDate, endDate time.Time) ([]entity.AuditLog, int64, error)
	FindAuditLogsByEntityID(entityType string, entityID uint) ([]entity.AuditLog, error)
}

type auditRepository struct {
	db *gorm.DB
}

func NewAuditRepository() AuditRepository {
	return &auditRepository{
		db: config.DB,
	}
}

func (r *auditRepository) CreateAuditLog(auditLog *entity.AuditLog) error {
	return r.db.Create(auditLog).Error
}

func (r *auditRepository) FindAuditLogs(page, limit int, userID uint, entityType string, startDate, endDate time.Time) ([]entity.AuditLog, int64, error) {
	var auditLogs []entity.AuditLog
	var count int64
	
	offset := (page - 1) * limit
	query := r.db
	
	// Add filters
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	
	if entityType != "" {
		query = query.Where("entity_type = ?", entityType)
	}
	
	// Add date range filter if provided
	if !startDate.IsZero() && !endDate.IsZero() {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	} else if !startDate.IsZero() {
		query = query.Where("created_at >= ?", startDate)
	} else if !endDate.IsZero() {
		query = query.Where("created_at <= ?", endDate)
	}
	
	// Get total count
	if err := query.Model(&entity.AuditLog{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated audit logs with user info
	if err := query.Preload("User").Order("created_at DESC").Offset(offset).Limit(limit).Find(&auditLogs).Error; err != nil {
		return nil, 0, err
	}
	
	return auditLogs, count, nil
}

func (r *auditRepository) FindAuditLogsByEntityID(entityType string, entityID uint) ([]entity.AuditLog, error) {
	var auditLogs []entity.AuditLog
	
	if err := r.db.Where("entity_type = ? AND entity_id = ?", entityType, entityID).
		Preload("User").
		Order("created_at DESC").
		Find(&auditLogs).Error; err != nil {
		return nil, err
	}
	
	return auditLogs, nil
} 