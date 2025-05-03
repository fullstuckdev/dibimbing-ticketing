package repository

import (
	"errors"

	"github.com/taufikmulyawan/ticketing-system/config"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"gorm.io/gorm"
)

type TicketRepository interface {
	FindAll(page, limit int, userID uint) ([]entity.Ticket, int64, error)
	FindByID(id uint) (*entity.Ticket, error)
	Save(ticket *entity.Ticket) error
	FindByEventID(eventID uint) ([]entity.Ticket, error)
	CountSoldTicketsByEventID(eventID uint) (int64, error)
}

type ticketRepository struct {
	db *gorm.DB
}

func NewTicketRepository() TicketRepository {
	return &ticketRepository{
		db: config.DB,
	}
}

func (r *ticketRepository) FindAll(page, limit int, userID uint) ([]entity.Ticket, int64, error) {
	var tickets []entity.Ticket
	var count int64

	offset := (page - 1) * limit
	query := r.db

	// Filter by user if userID is provided
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	// Get total count
	if err := query.Model(&entity.Ticket{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get tickets with pagination
	if err := query.Preload("Event").Preload("User").Offset(offset).Limit(limit).Find(&tickets).Error; err != nil {
		return nil, 0, err
	}

	return tickets, count, nil
}

func (r *ticketRepository) FindByID(id uint) (*entity.Ticket, error) {
	var ticket entity.Ticket
	result := r.db.Preload("Event").Preload("User").First(&ticket, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("ticket not found")
		}
		return nil, result.Error
	}
	return &ticket, nil
}

func (r *ticketRepository) Save(ticket *entity.Ticket) error {
	return r.db.Save(ticket).Error
}

func (r *ticketRepository) FindByEventID(eventID uint) ([]entity.Ticket, error) {
	var tickets []entity.Ticket
	if err := r.db.Where("event_id = ?", eventID).Find(&tickets).Error; err != nil {
		return nil, err
	}
	return tickets, nil
}

func (r *ticketRepository) CountSoldTicketsByEventID(eventID uint) (int64, error) {
	var count int64
	err := r.db.Model(&entity.Ticket{}).
		Where("event_id = ? AND status = ?", eventID, entity.TicketStatusPurchased).
		Count(&count).Error
	return count, err
} 