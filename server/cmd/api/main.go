package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"manumental-effort/server/internal/auth"
	"manumental-effort/server/internal/platform/config"
	"manumental-effort/server/internal/platform/mongodb"
	"manumental-effort/server/internal/users"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg, err := config.Load("configs/app-local.yaml")
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoClient, err := mongodb.Connect(ctx, cfg.MongoDB.URI, cfg.MongoDB.Database)
	if err != nil {
		log.Fatalf("connect mongodb: %v", err)
	}

	log.Printf("connected to mongodb database=%s", mongoClient.Database.Name())

	userRepository := users.NewRepository(mongoClient.Database)

	indexCtx, indexCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer indexCancel()

	if err := userRepository.EnsureIndexes(indexCtx); err != nil {
		log.Fatalf("ensure user indexes: %v", err)
	}

	authRepository := auth.NewRepository(mongoClient.Database)
	tokenManager := auth.NewTokenManager(cfg.Auth.JWTSigningKey, cfg.Auth.TokenExpiryMinutes)
	authService := auth.NewService(authRepository, tokenManager)
	authHandler := auth.NewHandler(authService)

	userService := users.NewService(userRepository, authRepository)
	userHandler := users.NewHandler(userService)

	r := gin.Default()

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":   "ok",
			"database": mongoClient.Database.Name(),
		})
	})

	r.POST("/users", userHandler.CreateUser)
	r.GET("/users/:id", userHandler.GetUserByID)
	r.POST("/auth/login", authHandler.Login)

	authGroup := r.Group("/auth")
	authGroup.Use(auth.AuthMiddleware(tokenManager))
	authGroup.GET("/me", userHandler.GetCurrentUser)

	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("run server: %v", err)
	}
}
