package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go_demo_api/models"
)

type UserHandler struct {
	db *gorm.DB
}

func NewUserHandler(db *gorm.DB) *UserHandler {
	return &UserHandler{db: db}
}

func (h *UserHandler) GetUsers(c *gin.Context) {
	// Get query parameters
	page := c.DefaultQuery("page", "1")
	text := c.Query("text")

	// Convert page to integer
	var pageInt int
	fmt.Sscanf(page, "%d", &pageInt)
	if pageInt < 1 {
		pageInt = 1
	}

	// Pagination settings
	pageSize := 10
	offset := (pageInt - 1) * pageSize

	// Build query
	var users []models.User
	var total int64
	var col string = `id, name, telephone, reportGroupId, role, status, 
	DATE_FORMAT(createdAt, '%d/%m/%Y %r') as createdAt, 
	DATE_FORMAT(updatedAt, '%d/%m/%Y %r') as updatedAt,
	DATE_FORMAT(deletedAt, '%d/%m/%Y %r') as deletedAt`

	// col += "DATE_FORMAT(createdAt, '%d/%m/%Y %r') as createdAt"
	// col += "DATE_FORMAT(updatedAt, '%d/%m/%Y %r') as updatedAt"
	// col += "DATE_FORMAT(deletedAt, '%d/%m/%Y %r') as deletedAt"

	query := h.db.Model(&models.User{}).Select(col)

	// Apply text filter if provided
	if text != "" {
		text = "%" + text + "%"
		query = query.Where("name LIKE ? OR telephone LIKE ?", text, text)
	}

	// Get total count for pagination info
	query.Count(&total)

	// Apply pagination and get results
	query.Offset(offset).Limit(pageSize).Find(&users)

	// Return response with pagination metadata
	c.JSON(200, gin.H{
		"data":        users,
		"page":        pageInt,
		"page_size":   pageSize,
		"total":       total,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	err := h.db.First(&user, id).Error
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "message": "User not found", "data": user})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": user})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.SaveUser
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	h.db.Create(&user)
	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := h.db.First(&user, id).Error; err != nil {
		c.JSON(200, gin.H{"error": "User not found"})
		return
	}
	var input models.SaveUser
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(200, gin.H{"error": err.Error()})
		return
	}
	h.db.Model(&user).Updates(input)
	c.JSON(200, gin.H{"data": user})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	if err := h.db.Delete(&user, id).Error; err != nil {
		c.JSON(200, gin.H{"error": "User not found"})
		return
	}
	c.JSON(200, gin.H{"message": "User deleted"})
}
