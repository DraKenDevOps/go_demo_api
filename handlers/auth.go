package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"

	"go_demo_api/logger"
	"go_demo_api/middleware"
	"go_demo_api/models"
	"go_demo_api/utils"
)

// type AuthHandler struct {
// 	db  *gorm.DB
// 	cfg *config.Config
// }

// func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
// 	return &AuthHandler{db: db, cfg: cfg}
// }

type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// type LoginRequest struct {
// 	Telephone string `json:"telephone" binding:"required"`
// 	Password string `json:"password" binding:"required"`
// }

type LoginResponse struct {
	AccessToken string          `json:"access_token"`
	User        models.AuthUser `json:"user"`
	Status      string          `json:"status"`
	Message     string          `json:"message"`
}

type RefreshResponse struct {
	AccessToken string              `json:"access_token"`
	User        models.JwtUserClaim `json:"user"`
	Status      string              `json:"status"`
	Message     string              `json:"message"`
}

func (h *ApiHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"status": "error", "message": "Failed to read request body"})
		return
	}

	var sql string = "SELECT user_id, telephone, password, email, role_action, op_id, level, status FROM users WHERE username = ?"
	var user models.AuthUser
	err := h.db.Exec(sql, req.Username).First(&user).Error
	if err != nil {
		msg := fmt.Sprintf("Internal server error: User %s not found %v", req.Username, err)
		logger.Log.Error().Msg(msg)
		c.JSON(200, gin.H{"status": "error", "message": "Failed to log in"})
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		msg := fmt.Sprintf("User %s wrong password", req.Username)
		logger.Log.Info().Msg(msg)
		c.JSON(200, gin.H{"status": "error", "message": "Invalid credentials"})
		return
	}

	if user.Status != "ACTIVE" {
		msg := fmt.Sprintf("User %s is not active", req.Username)
		logger.Log.Info().Msg(msg)
		c.JSON(200, gin.H{"status": "error", "message": "Invalid credentials"})
		return
	}

	token, err := utils.CreateToken(user.UserId, user.OpId, user.Username, user.Email, user.Telephone, user.Level, user.RoleAction, h.cfg)
	if err != nil {
		msg := fmt.Sprintf("User %s failed to generate token", req.Username)
		logger.Log.Info().Msg(msg)
		c.JSON(200, gin.H{"status": "error", "message": "Failed to log in"})
		return
	}

	c.JSON(200, LoginResponse{
		AccessToken: token,
		User:        user,
		Status:      "success",
		Message:     "Welcome Back",
	})
}

func (h *ApiHandler) Refresh(c *gin.Context) {
	accessToken := c.GetString("jwt")
	if accessToken == "" {
		c.JSON(200, gin.H{"status": "error", "message": "Failed to authentication"})
		return
	}

	user := middleware.GetUserClaim(c)

	c.JSON(200, RefreshResponse{
		AccessToken: accessToken,
		Status:      "success",
		User:        user,
		Message:     "Welcome Back",
	})
}
