package main

import (
	"log"
	"os"
	"sync"

	"banking-app/app"
	"banking-app/middleware"

	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	// Load .env
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	// Set JWT secret in middleware
	middleware.SetJWTSecret(jwtSecret)

	var wg sync.WaitGroup
	wg.Add(1)

	// Create the app instance
	myApp := app.NewApp("Banking App", &wg, jwtSecret)

	// Wrap the router with CORS middleware
	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"}, // Angular URL
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	myApp.Server.Handler = c.Handler(myApp.Router) // Wrap the router with CORS

	log.Printf("Server running on port %s...", os.Getenv("PORT"))
	// Start the server
	if err := myApp.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	wg.Wait()
}
