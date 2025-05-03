package router

import (
	"github.com/gin-gonic/gin"
	"github.com/taufikmulyawan/ticketing-system/controller"
	"github.com/taufikmulyawan/ticketing-system/service"
)

// InitRouter initializes the router with all controllers and services
func InitRouter(controllers *controller.Controllers, auditService service.AuditService) *gin.Engine {
	return SetupRouter(
		controllers.UserController,
		controllers.EventController,
		controllers.TicketController,
		controllers.ReportController,
		controllers.AuditController,
		auditService,
	)
} 