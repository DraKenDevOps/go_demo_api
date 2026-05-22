package routes

import (
	"github.com/gin-gonic/gin"

	"go_demo_api/config"
	"go_demo_api/handlers"
	"go_demo_api/middleware"
)

func SetupRoutes(r *gin.Engine, cfg *config.Config, user *handlers.UserHandler, auth *handlers.AuthHandler, upload *handlers.UploadHandler) {
	api := r.Group(cfg.BasePath)
	{
		api.POST("/login", auth.Login)
		api.Use(middleware.AuthMiddleware(cfg)).GET("/refresh", auth.Refresh)
		api.Use(middleware.AuthMiddleware(cfg)).GET("/users", user.GetUsersRaw)
		api.Use(middleware.AuthMiddleware(cfg)).GET("/user/:id", user.GetUser)
		api.Use(middleware.AuthMiddleware(cfg)).POST("/create_user", user.CreateUser)
		api.Use(middleware.AuthMiddleware(cfg)).PUT("/update_user/:id", user.UpdateUser)
		api.Use(middleware.AuthMiddleware(cfg)).PATCH("/delete_user/:id", user.SoftDeleteUser)
		api.Use(middleware.AuthMiddleware(cfg)).DELETE("/delete_user/:id", user.HardDeleteUser)
		api.Use(middleware.AuthMiddleware(cfg)).POST("/upload", upload.UploadFile)
	}
}
