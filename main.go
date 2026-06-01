package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"

	"go_demo_api/config"
	"go_demo_api/database"
	"go_demo_api/handlers"
	"go_demo_api/jobs"
	"go_demo_api/logger"
	"go_demo_api/middleware"
	"go_demo_api/routes"
	"go_demo_api/scheduler"
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

	schedule := scheduler.New()

	jobList := []scheduler.Job{
		jobs.CleanupJob{},
		jobs.EmailJob{},
	}

	for _, job := range jobList {
		if err := schedule.Register(job); err != nil {
			fmt.Printf("Run cron job error: %s", err.Error())
		}
	}

	schedule.Start()

	router := gin.Default()
	router.MaxMultipartMemory = 50 << 20

	router.Static("/static", filepath.Join(cfg.Cwd, "uploads"))

	router.Use(middleware.Logger("/api/upload"))
	router.GET("/health", func(c *gin.Context) {
		res := gin.H{
			"status":    "ok",
			"timestamp": time.Now().Unix(),
			"uptime":    time.Since(time.Unix(0, 0)) / time.Second,
			"version":   cfg.Version,
		}
		c.JSON(200, res)
	})
	routes.SetupRoutes(router, cfg, handlers.NewApiHandler(db, cfg))

	addr := cfg.GetServerAddress()
	server := &http.Server{
		Addr: addr, Handler: router,
	}

	logger.Log.Info().Msg(fmt.Sprintf("Server starting on %s", addr))
	errr := router.Run(addr)
	if errr != nil {
		msg := fmt.Sprintf("Failed to start server: %v", err)
		logger.Log.Error().Msg(msg)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	<-quit
	log.Println("shutdown signal received")

	schedule.Stop()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("graceful shutdown complete")
}
