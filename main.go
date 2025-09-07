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
	_ = godotenv.Load()

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET is required")
	}

	middleware.SetJWTSecret(jwtSecret)

	var wg sync.WaitGroup
	wg.Add(1)

	myApp := app.NewApp("Banking App", &wg, jwtSecret)

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	myApp.Server.Handler = c.Handler(myApp.Router)

	log.Printf("Server running on port %s...", os.Getenv("PORT"))
	if err := myApp.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	wg.Wait()
}
