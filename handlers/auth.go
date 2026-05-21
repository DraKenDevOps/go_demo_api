package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go_demo_api/config"
	"go_demo_api/logger"
	"go_demo_api/middleware"
	"go_demo_api/models"
	"go_demo_api/utils"
)

type AuthHandler struct {
	db  *gorm.DB
	cfg *config.Config
}

func NewAuthHandler(db *gorm.DB, cfg *config.Config) *AuthHandler {
	return &AuthHandler{db: db, cfg: cfg}
}

type LoginRequest struct {
	Telephone string `json:"telephone" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

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

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(200, gin.H{"status": "error", "message": "Failed to read request body"})
		return
	}

	var sql string = "SELECT id, name, password, reportGroupId, role, status FROM users WHERE telephone = ?"
	var user models.AuthUser
	err := h.db.Exec(sql, req.Telephone).First(&user).Error
	// err := h.db.Select("id, password, reportGroupId, role, status").Where("telephone = ?", req.Telephone).First(&user).Error
	if err != nil {
		msg := fmt.Sprintf("Internal server error: User phone %s not found %v", req.Telephone, err)
		logger.Log.Error().Msg(msg)
		c.JSON(200, gin.H{"status": "error", "message": "Failed to log in"})
		return
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		msg := fmt.Sprintf("User phone %s wrong password", req.Telephone)
		logger.Log.Info().Msg(msg)
		c.JSON(200, gin.H{"status": "error", "message": "Invalid credentials"})
		return
	}

	if user.Status != "active" {
		msg := fmt.Sprintf("User phone %s is not active", req.Telephone)
		logger.Log.Info().Msg(msg)
		c.JSON(200, gin.H{"status": "error", "message": "Invalid credentials"})
		return
	}

	token, err := utils.CreateToken(user.ID, user.ReportGroupID, user.Name, user.Telephone, user.Role, h.cfg)
	if err != nil {
		msg := fmt.Sprintf("User phone %s failed to generate token", req.Telephone)
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

func (h *AuthHandler) Refresh(c *gin.Context) {
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
