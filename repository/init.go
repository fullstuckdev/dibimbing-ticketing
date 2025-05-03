package repository

// Repositories holds all repository instances
type Repositories struct {
	UserRepository   UserRepository
	EventRepository  EventRepository
	TicketRepository TicketRepository
	AuditRepository  AuditRepository
}

// InitRepositories initializes all repositories
func InitRepositories() *Repositories {
	return &Repositories{
		UserRepository:   NewUserRepository(),
		EventRepository:  NewEventRepository(),
		TicketRepository: NewTicketRepository(),
		AuditRepository:  NewAuditRepository(),
	}
} 