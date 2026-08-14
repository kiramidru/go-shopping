package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"kiramidru/go-shopping/internal/catalog"
	"kiramidru/go-shopping/internal/config"
	"kiramidru/go-shopping/internal/db"
	"kiramidru/go-shopping/internal/users"
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

	tokenManager := pkgs.NewTokenManager(cfg.JWTSecret, 15*time.Minute, 7*24*time.Hour)
	userRepo := users.NewRepository(db)
	userService := users.NewService(userRepo, tokenManager)
	userHandler := users.NewHandler(userService)

	catalogService := catalog.NewService(http.DefaultClient, cfg.CatalogURL, cfg.CatalogAPIKey)
	catalogHandler := catalog.NewHandler(catalogService)

	if cfg.Env == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	port := fmt.Sprintf(":%s", cfg.Port)

	r := gin.Default()
	api := r.Group("/api/v1")
	userHandler.RegisterRoutes(api)
	catalogHandler.RegisterRoutes(api)
	r.Run(port)
}
