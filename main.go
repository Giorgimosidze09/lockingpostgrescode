package main

import (
	"log"
	"net/http"

	db "lockingpostgrescode/database"
	_ "lockingpostgrescode/docs"
	"lockingpostgrescode/handlers"

	"github.com/rs/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Account Management API
// @version 1.0
// @description API for managing accounts with withdrawal, deposit, and exchange operations.
// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com
// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization

func main() {
	// Connect to the database
	db, err := db.ConnectDB()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Setup routes
	r := handlers.SetupRoutes(db)

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Allow frontend URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

	// Apply CORS middleware
	r.Use(corsHandler.Handler)

	// Swagger documentation
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Start the server
	log.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", r))
}
