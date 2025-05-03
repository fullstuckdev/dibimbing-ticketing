package router

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/taufikmulyawan/ticketing-system/controller"
	"github.com/taufikmulyawan/ticketing-system/middleware"
	"github.com/taufikmulyawan/ticketing-system/service"
)

// SetupRouter configures all the routes for the application
func SetupRouter(
	userController controller.UserController,
	eventController controller.EventController,
	ticketController controller.TicketController,
	reportController controller.ReportController,
	auditController controller.AuditController,
	auditService service.AuditService,
) *gin.Engine {
	// Initialize router
	router := gin.Default()

	// Add middleware for CORS
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Add audit middleware to all routes
	router.Use(middleware.AuditMiddleware(auditService))

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Public routes
	router.POST("/register", userController.Register)
	router.POST("/login", userController.Login)
	router.GET("/events", eventController.GetAllEvents)
	router.GET("/events/:id", eventController.GetEventByID)

	// Protected routes
	authRoutes := router.Group("/")
	authRoutes.Use(middleware.AuthMiddleware())
	{
		// User routes
		authRoutes.GET("/profile", userController.Profile)
		authRoutes.POST("/logout", userController.Logout)
		authRoutes.GET("/my-audit-logs", userController.GetMyAuditLogs)

		// Ticket routes
		authRoutes.GET("/tickets", ticketController.GetAllTickets)
		authRoutes.GET("/tickets/:id", ticketController.GetTicketByID)
		authRoutes.POST("/tickets", ticketController.PurchaseTicket)
		authRoutes.PATCH("/tickets/:id", ticketController.CancelTicket)
	}

	// Admin routes
	adminRoutes := router.Group("/")
	adminRoutes.Use(middleware.AdminMiddleware())
	{
		// Event management
		adminRoutes.POST("/events", eventController.CreateEvent)
		adminRoutes.PUT("/events/:id", eventController.UpdateEvent)
		adminRoutes.DELETE("/events/:id", eventController.DeleteEvent)

		// Reports
		adminRoutes.GET("/reports/summary", reportController.GetSalesReport)
		adminRoutes.GET("/reports/event/:id", reportController.GetEventSalesReport)
		
		// Audit logs (admin only)
		adminRoutes.GET("/audit/logs", auditController.GetAuditLogs)
		adminRoutes.GET("/audit/:entity_type/:entity_id", auditController.GetEntityAuditLogs)
	}

	return router
} 