package routes

import (
	"github.com/gin-gonic/gin"

	"go_demo_api/config"
	"go_demo_api/handlers"
	"go_demo_api/middleware"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, apiHandler *handlers.ApiHandler) {
	api := r.Group(cfg.BasePath)
	{
		api.POST("/login", apiHandler.Login)
		api.Use(middleware.AuthChecker(cfg)).GET("/refresh", apiHandler.Refresh)
		api.Use(middleware.AuthChecker(cfg)).GET("/users", apiHandler.GetUsersRaw)
		api.Use(middleware.AuthChecker(cfg)).GET("/user/:id", apiHandler.GetUser)
		api.Use(middleware.AuthChecker(cfg)).POST("/create_user", apiHandler.CreateUser)
		api.Use(middleware.AuthChecker(cfg)).PUT("/update_user/:id", apiHandler.UpdateUser)
		api.Use(middleware.AuthChecker(cfg)).PATCH("/delete_user/:id", apiHandler.SoftDeleteUser)
		api.Use(middleware.AuthChecker(cfg)).DELETE("/delete_user/:id", apiHandler.HardDeleteUser)
		api.Use(middleware.AuthChecker(cfg)).POST("/upload", apiHandler.UploadFile)
	}
}
