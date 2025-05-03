package controller

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/service"
)

type EventController interface {
	GetAllEvents(c *gin.Context)
	GetEventByID(c *gin.Context)
	CreateEvent(c *gin.Context)
	UpdateEvent(c *gin.Context)
	DeleteEvent(c *gin.Context)
}

type eventController struct {
	eventService service.EventService
}

func NewEventController(eventService service.EventService) EventController {
	return &eventController{
		eventService: eventService,
	}
}

// GetAllEvents godoc
// @Summary Get all events
// @Description Get a list of all events with pagination
// @Tags events
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Success 200 {object} map[string]interface{}
// @Router /events [get]
func (ctrl *eventController) GetAllEvents(c *gin.Context) {
	page, limit := GetPaginationParams(c)

	events, count, err := ctrl.eventService.GetAllEvents(page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GeneratePaginationResponse(events, page, limit, count))
}

// GetEventByID godoc
// @Summary Get event by ID
// @Description Get details of a specific event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Success 200 {object} entity.Event
// @Failure 404 {object} map[string]interface{}
// @Router /events/{id} [get]
func (ctrl *eventController) GetEventByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	event, err := ctrl.eventService.GetEventByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Event not found"})
		return
	}

	c.JSON(http.StatusOK, event)
}

// CreateEvent godoc
// @Summary Create a new event
// @Description Create a new event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Param event body entity.Event true "Event Data"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /events [post]
func (ctrl *eventController) CreateEvent(c *gin.Context) {
	var event entity.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Parse date strings if they come in string format
	if event.StartDate.IsZero() || event.EndDate.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date and end_date are required and must be valid dates"})
		return
	}

	// Validate event data
	if event.Name == "" || event.Location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and location are required"})
		return
	}

	if event.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be in the future"})
		return
	}

	if event.EndDate.Before(event.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	if event.Capacity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "capacity must be greater than 0"})
		return
	}

	if event.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price cannot be negative"})
		return
	}

	err := ctrl.eventService.CreateEvent(&event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Event created successfully", "event_id": event.ID})
}

// UpdateEvent godoc
// @Summary Update an event
// @Description Update an existing event with the provided details
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Param event body entity.Event true "Event Data"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Router /events/{id} [put]
func (ctrl *eventController) UpdateEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	var event entity.Event
	if err := c.ShouldBindJSON(&event); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate event data
	if event.Name == "" || event.Location == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "name and location are required"})
		return
	}

	if event.StartDate.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "start_date must be in the future"})
		return
	}

	if event.EndDate.Before(event.StartDate) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "end_date must be after start_date"})
		return
	}

	if event.Capacity <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "capacity must be greater than 0"})
		return
	}

	if event.Price < 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "price cannot be negative"})
		return
	}

	err = ctrl.eventService.UpdateEvent(uint(id), &event)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event updated successfully"})
}

// DeleteEvent godoc
// @Summary Delete an event
// @Description Delete an existing event by its ID
// @Tags events
// @Accept json
// @Produce json
// @Param id path int true "Event ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Router /events/{id} [delete]
func (ctrl *eventController) DeleteEvent(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid event ID"})
		return
	}

	err = ctrl.eventService.DeleteEvent(uint(id))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Event deleted successfully"})
} 