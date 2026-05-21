package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"go_demo_api/models"
	"go_demo_api/utils"
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
	var col string = `user_id, username, telephone, email, op_id, level, role_action, status, 
	DATE_FORMAT(created_at, '%d/%m/%Y %r') as created_at, 
	DATE_FORMAT(updated_at, '%d/%m/%Y %r') as updated_at,
	DATE_FORMAT(deleted_at, '%d/%m/%Y %r') as deleted_at`

	// col += "DATE_FORMAT(created_at, '%d/%m/%Y %r') as created_at"
	// col += "DATE_FORMAT(updated_at, '%d/%m/%Y %r') as updated_at"
	// col += "DATE_FORMAT(deleted_at, '%d/%m/%Y %r') as deleted_at"

	query := h.db.Model(&models.User{}).Select(col)

	// Apply text filter if provided
	if text != "" {
		text = "%" + text + "%"
		query = query.Where("username LIKE ? OR telephone LIKE ?", text, text)
	}

	// Get total count for pagination info
	query.Count(&total)

	// Apply pagination and get results
	query.Offset(offset).Limit(pageSize).Find(&users)

	// Return response with pagination metadata
	c.JSON(200, gin.H{
		"data":        users,
		"page":        pageInt,
		"total":       total,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func (h *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")
	var col string = `user_id, username, telephone, email, op_id, level, role_action, status, 
	DATE_FORMAT(created_at, '%d/%m/%Y %r') as created_at, 
	DATE_FORMAT(updated_at, '%d/%m/%Y %r') as updated_at,
	DATE_FORMAT(deleted_at, '%d/%m/%Y %r') as deleted_at`
	var user models.User
	err := h.db.Select(col).First(&user, id).Error
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "message": "User not found", "data": user})
		return
	}
	c.JSON(200, gin.H{"status": "success", "data": user})
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.SaveUser
	if err := c.ShouldBindJSON(&user); err != nil {
		fmt.Printf("Failed to read request body: %v\n", err)
		c.JSON(200, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}

	if user.Password != nil || *user.Password != "" {
		hash, err := utils.HashPassword(*user.Password)
		if err != nil {
			fmt.Printf("Failed to hash password: %v\n", err)
			c.JSON(200, gin.H{"status": "error", "message": "Failed to create user"})
			return
		}
		user.Password = &hash
	}

	if user.Status == nil || *user.Status == "" {
		status := "ACTIVE"
		user.Status = &status
	}

	err := h.db.Create(&user).Error
	if err != nil {
		fmt.Printf("Failed to create user: %v\n", err)
		c.JSON(200, gin.H{"status": "error", "message": "Failed to create user"})
		return
	}

	c.JSON(200, gin.H{"status": "success", "message": "User created successfully"})
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	var user models.SaveUserPassword
	if err := h.db.First(&user, id).Error; err != nil {
		fmt.Printf("User not found: %v\n", err)
		c.JSON(200, gin.H{"status": "error", "message": "User not found"})
		return
	}
	var input models.SaveUser
	if err := c.ShouldBindJSON(&input); err != nil {
		fmt.Printf("Failed to read request body: %v\n", err)
		c.JSON(200, gin.H{"status": "error", "message": "Invalid request body"})
		return
	}

	// TODO: if in handler, check if password is not in request body, then keep existing password

	// TODO: if in handler, check if password is provided and different from existing one, then hash it before saving
	if input.Password != nil || *input.Password != "" {
		valid := utils.CheckPassword(*input.Password, user.Password)
		if !valid {
			hash, err := utils.HashPassword(*input.Password)
			if err != nil {
				fmt.Printf("Failed to hash password: %v\n", err)
				c.JSON(200, gin.H{"status": "error", "message": "Failed to update user"})
				return
			}
			input.Password = &hash
		}
	}

	if input.Status == nil || *input.Status == "" {
		status := "ACTIVE"
		input.Status = &status
	}

	err := h.db.Model(&user).Updates(input).Where("user_id = ?", id).Error
	if err != nil {
		fmt.Printf("Failed to update user: %v\n", err)
		c.JSON(200, gin.H{"status": "error", "message": "Failed to update user"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "User updated successfully"})
}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	var user models.User
	err := h.db.Delete(&user, id).Error
	if err != nil {
		c.JSON(200, gin.H{"status": "error", "message": "User not found"})
		return
	}
	c.JSON(200, gin.H{"status": "success", "message": "User deleted"})
}
