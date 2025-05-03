package repository

import (
	"errors"

	"github.com/taufikmulyawan/ticketing-system/config"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"gorm.io/gorm"
)

type EventRepository interface {
	FindAll(page, limit int) ([]entity.Event, int64, error)
	FindByID(id uint) (*entity.Event, error)
	Save(event *entity.Event) error
	Delete(id uint) error
}

type eventRepository struct {
	db *gorm.DB
}

func NewEventRepository() EventRepository {
	return &eventRepository{
		db: config.DB,
	}
}

func (r *eventRepository) FindAll(page, limit int) ([]entity.Event, int64, error) {
	var events []entity.Event
	var count int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&entity.Event{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated events
	if err := r.db.Offset(offset).Limit(limit).Find(&events).Error; err != nil {
		return nil, 0, err
	}

	return events, count, nil
}

func (r *eventRepository) FindByID(id uint) (*entity.Event, error) {
	var event entity.Event
	result := r.db.First(&event, id)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errors.New("event not found")
		}
		return nil, result.Error
	}
	return &event, nil
}

func (r *eventRepository) Save(event *entity.Event) error {
	return r.db.Save(event).Error
}

func (r *eventRepository) Delete(id uint) error {
	// First check if event exists
	event, err := r.FindByID(id)
	if err != nil {
		return err
	}

	// Then check if there are tickets sold for this event
	var ticketCount int64
	if err := r.db.Model(&entity.Ticket{}).Where("event_id = ? AND status = ?", id, entity.TicketStatusPurchased).Count(&ticketCount).Error; err != nil {
		return err
	}

	if ticketCount > 0 {
		return errors.New("cannot delete event with sold tickets")
	}

	return r.db.Delete(event).Error
} 