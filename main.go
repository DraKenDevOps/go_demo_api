package main

import (
	"fmt"
	"path/filepath"

	"github.com/gin-gonic/gin"

	"go_demo_api/config"
	"go_demo_api/database"
	"go_demo_api/handlers"
	"go_demo_api/logger"
	"go_demo_api/middleware"
	"go_demo_api/routes"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		msg := fmt.Sprintf("Failed to load config: %v", err)
		logger.Log.Error().Msg(msg)
	}

	logger.Init(cfg)

	db, err := database.Connect(cfg)
	if err != nil {
		msg := fmt.Sprintf("Failed to connect to database: %v", err)
		logger.Log.Error().Msg(msg)
	}

	if cfg.EnvMode == "production" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.MaxMultipartMemory = 50 << 20

	router.Static("/public", filepath.Join(cfg.Cwd, "uploads"))

	router.Use(middleware.Logger())
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})
	routes.SetupRoutes(router, cfg, handlers.NewUserHandler(db), handlers.NewAuthHandler(db, cfg), handlers.NewUploadHandler(db, cfg))

	addr := cfg.GetServerAddress()
	logger.Log.Info().Msg(fmt.Sprintf("Server starting on %s", addr))
	errr := router.Run(addr)
	if errr != nil {
		msg := fmt.Sprintf("Failed to start server: %v", err)
		logger.Log.Error().Msg(msg)
	}
}
