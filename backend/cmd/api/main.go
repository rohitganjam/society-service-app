package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/rohit/society-service-app/backend/internal/config"
	"github.com/rohit/society-service-app/backend/internal/database"
	"github.com/rohit/society-service-app/backend/internal/handlers"
	"github.com/rohit/society-service-app/backend/internal/middleware"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Set Gin mode
	if cfg.Environment == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize database (optional - will work without DB for initial testing)
	var db *database.DB
	if cfg.DatabaseURL != "" {
		db, err = database.New(cfg.DatabaseURL)
		if err != nil {
			log.Printf("Warning: Failed to connect to database: %v", err)
			log.Println("Server will start without database connection")
		} else {
			defer db.Close()
			log.Println("Connected to database")
		}
	} else {
		log.Println("DATABASE_URL not set, running without database")
	}

	// Initialize router
	router := gin.New()

	// Apply middleware
	router.Use(middleware.Recovery())
	router.Use(middleware.Logger())
	router.Use(middleware.CORS())

	// Initialize handlers
	healthHandler := handlers.NewHealthHandler(db)

	// Routes
	router.GET("/health", healthHandler.Health)
	router.GET("/ready", healthHandler.Ready)

	// API v1 routes
	v1 := router.Group("/api/v1")
	{
		v1.GET("/health", healthHandler.Health)

		// Auth routes (placeholder)
		// auth := v1.Group("/auth")
		// {
		// 	auth.POST("/send-otp", authHandler.SendOTP)
		// 	auth.POST("/verify-otp", authHandler.VerifyOTP)
		// }

		// Add more route groups here as you implement them
	}

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		log.Println("Shutting down server...")
		if db != nil {
			db.Close()
		}
		os.Exit(0)
	}()

	// Start server
	addr := ":" + cfg.Port
	log.Printf("Server starting on %s (environment: %s)", addr, cfg.Environment)
	if err := router.Run(addr); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
