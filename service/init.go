package service

import "github.com/taufikmulyawan/ticketing-system/repository"

// Services holds all service instances
type Services struct {
	UserService   UserService
	EventService  EventService
	TicketService TicketService
	ReportService ReportService
	AuditService  AuditService
}

// InitServices initializes all services with their required repositories
func InitServices(repos *repository.Repositories) *Services {
	return &Services{
		UserService:   NewUserService(repos.UserRepository),
		EventService:  NewEventService(repos.EventRepository),
		TicketService: NewTicketService(repos.TicketRepository, repos.EventRepository),
		ReportService: NewReportService(repos.TicketRepository, repos.EventRepository),
		AuditService:  NewAuditService(repos.AuditRepository),
	}
} 