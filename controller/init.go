package controller

import "github.com/taufikmulyawan/ticketing-system/service"

// Controllers holds all controller instances
type Controllers struct {
	UserController   UserController
	EventController  EventController
	TicketController TicketController
	ReportController ReportController
	AuditController  AuditController
}

// InitControllers initializes all controllers with their required services
func InitControllers(services *service.Services) *Controllers {
	return &Controllers{
		UserController:   NewUserController(services.UserService, services.AuditService),
		EventController:  NewEventController(services.EventService, services.AuditService),
		TicketController: NewTicketController(services.TicketService, services.AuditService),
		ReportController: NewReportController(services.ReportService),
		AuditController:  NewAuditController(services.AuditService),
	}
} 