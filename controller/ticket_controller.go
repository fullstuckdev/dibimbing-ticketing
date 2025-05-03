package controller

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/taufikmulyawan/ticketing-system/entity"
	"github.com/taufikmulyawan/ticketing-system/service"
)

type TicketController interface {
	GetAllTickets(c *gin.Context)
	GetTicketByID(c *gin.Context)
	PurchaseTicket(c *gin.Context)
	CancelTicket(c *gin.Context)
}

type ticketController struct {
	ticketService service.TicketService
	auditService  service.AuditService
}

func NewTicketController(ticketService service.TicketService, auditService service.AuditService) TicketController {
	return &ticketController{
		ticketService: ticketService,
		auditService:  auditService,
	}
}

// GetAllTickets godoc
// @Summary Get all tickets
// @Description Get a list of all tickets with pagination (admin sees all, user sees only their own)
// @Tags tickets
// @Accept json
// @Produce json
// @Param page query int false "Page number"
// @Param limit query int false "Items per page"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Router /tickets [get]
func (ctrl *ticketController) GetAllTickets(c *gin.Context) {
	page, limit := GetPaginationParams(c)

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Check if user is admin
	userRole, _ := c.Get("user_role")
	role, _ := userRole.(string)

	var userIDFilter uint
	if role == string(entity.RoleAdmin) {
		// Admin can see all tickets
		userIDFilter = 0
	} else {
		// Regular users can only see their own tickets
		userIDFilter = uint(id)
	}

	tickets, count, err := ctrl.ticketService.GetAllTickets(page, limit, userIDFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, GeneratePaginationResponse(tickets, page, limit, count))
}

// GetTicketByID godoc
// @Summary Get ticket by ID
// @Description Get details of a specific ticket by its ID
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Security BearerAuth
// @Success 200 {object} entity.Ticket
// @Failure 404 {object} map[string]interface{}
// @Router /tickets/{id} [get]
func (ctrl *ticketController) GetTicketByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	ticket, err := ctrl.ticketService.GetTicketByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}

	// Check if user is authorized to view this ticket
	userID, _ := c.Get("user_id")
	userRole, _ := c.Get("user_role")
	
	if userRole != string(entity.RoleAdmin) && ticket.UserID != uint(userID.(float64)) {
		c.JSON(http.StatusForbidden, gin.H{"error": "Not authorized to view this ticket"})
		return
	}

	c.JSON(http.StatusOK, ticket)
}

// PurchaseTicket godoc
// @Summary Purchase a ticket
// @Description Purchase a ticket for an event
// @Tags tickets
// @Accept json
// @Produce json
// @Param ticket body entity.Ticket true "Ticket Data (eventID is required)"
// @Security BearerAuth
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /tickets [post]
func (ctrl *ticketController) PurchaseTicket(c *gin.Context) {
	var ticket entity.Ticket
	if err := c.ShouldBindJSON(&ticket); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Set user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}

	// Validate ticket data
	if ticket.EventID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "event_id is required"})
		return
	}

	ticket.UserID = uint(id)

	// Store the pre-purchase ticket data for audit
	oldTicket, _ := json.Marshal(nil) // No old ticket exists
	
	err := ctrl.ticketService.PurchaseTicket(&ticket)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Explicitly log the ticket purchase in the audit trail
	newTicket, _ := json.Marshal(ticket)
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	
	go ctrl.auditService.LogActivity(
		uint(id),
		entity.ActionCreate,
		"ticket",
		ticket.ID,
		string(oldTicket),
		string(newTicket),
		ipAddress,
		userAgent,
	)

	c.JSON(http.StatusCreated, gin.H{"message": "Ticket purchased successfully", "ticket_id": ticket.ID})
}

// CancelTicket godoc
// @Summary Cancel a ticket
// @Description Cancel a purchased ticket
// @Tags tickets
// @Accept json
// @Produce json
// @Param id path int true "Ticket ID"
// @Security BearerAuth
// @Success 200 {object} map[string]interface{}
// @Failure 400,404 {object} map[string]interface{}
// @Router /tickets/{id} [patch]
func (ctrl *ticketController) CancelTicket(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ticket ID"})
		return
	}

	// Get user ID from token
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userIDUint, ok := userID.(float64)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid user ID"})
		return
	}
	
	// Get the ticket before cancellation for audit purposes
	oldTicket, _ := ctrl.ticketService.GetTicketByID(uint(id))
	if oldTicket == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Ticket not found"})
		return
	}
	
	oldTicketJSON, _ := json.Marshal(oldTicket)

	err = ctrl.ticketService.CancelTicket(uint(id), uint(userIDUint))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	// Get updated ticket after cancellation
	updatedTicket, _ := ctrl.ticketService.GetTicketByID(uint(id))
	updatedTicketJSON, _ := json.Marshal(updatedTicket)
	
	// Explicitly log the ticket cancellation in the audit trail
	ipAddress := c.ClientIP()
	userAgent := c.Request.UserAgent()
	
	go ctrl.auditService.LogActivity(
		uint(userIDUint),
		entity.ActionUpdate, // Cancellation is an update to the ticket status
		"ticket",
		uint(id),
		string(oldTicketJSON),
		string(updatedTicketJSON),
		ipAddress,
		userAgent,
	)

	c.JSON(http.StatusOK, gin.H{"message": "Ticket cancelled successfully"})
} 