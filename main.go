package main

import (
	"log"
	"sync"

	"banking-app/app"
	"banking-app/config"
	"banking-app/middleware"
)

func main() {
	cfg := config.LoadConfig()

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET is required in .env or environment variables")
	}

	middleware.SetJWTSecret(cfg.JWTSecret)

	var wg sync.WaitGroup
	wg.Add(1)

	myApp := app.NewApp("Banking App", &wg, cfg.JWTSecret)

	log.Printf("Server running on port %s...", cfg.Port)

	if err := myApp.Start(); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	wg.Wait()
}
