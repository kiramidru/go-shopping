package main

import (
	"fmt"
	"log"
	"time"

	"kiramidru/go-shopping/internal/auth"
	"kiramidru/go-shopping/internal/config"
	"kiramidru/go-shopping/internal/db"
	"kiramidru/go-shopping/pkgs"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	db, err := db.ConnectDatabase(cfg)
	if err != nil {
		log.Fatalf("Database connection failed: %v", err)
	}

	jwtManager := pkgs.NewJWTManager(cfg.JWTSecret, 15*time.Minute)
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo, jwtManager)
	authHandler := auth.NewHandler(authService)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	port := fmt.Sprintf(":%s", cfg.Port)

	r := gin.Default()
	api := r.Group("/api/v1")
	authHandler.RegisterRoutes(api)
	r.Run(port)
}
