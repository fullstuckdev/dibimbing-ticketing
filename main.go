package main

import (
	"fmt"
	"log"

	"github.com/taufikmulyawan/ticketing-system/config"
	"github.com/taufikmulyawan/ticketing-system/controller"
	_ "github.com/taufikmulyawan/ticketing-system/docs"
	"github.com/taufikmulyawan/ticketing-system/repository"
	"github.com/taufikmulyawan/ticketing-system/router"
	"github.com/taufikmulyawan/ticketing-system/service"
)

func main() {
	// Load configuration
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	config.ConnectDatabase()

	// Initialize all application components
	repositories := repository.InitRepositories()
	services := service.InitServices(repositories)
	controllers := controller.InitControllers(services)

	// Setup router
	r := router.InitRouter(controllers, services.AuditService)

	// Start the server
	port := config.AppConfig.Port
	if port == "" {
		port = "8080"
	}

	fmt.Printf("Server running on port %s\n", port)
	r.Run(":" + port)
} 