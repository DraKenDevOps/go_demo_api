package middleware

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"

	"go_demo_api/config"
	"go_demo_api/logger"
	"go_demo_api/models"
	"go_demo_api/utils"
)

type LoginResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func AuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		accessToken := c.GetHeader("X-Access-Token")
		authorization := c.GetHeader("Authorization")

		if authorization != "" {
			accessToken = strings.Replace(authorization, "Bearer ", "", 1)
		}
		if accessToken == "" {
			c.JSON(200, LoginResponse{Status: "error", Message: "Authorization header required"})
			c.Abort()
			return
		}

		claims, err := utils.Verify(accessToken, cfg)
		if err != nil {
			fmt.Print(err)
			logger.Log.Error().Err(err)
			c.JSON(200, LoginResponse{Status: "error", Message: "Invalid or expired token"})
			c.Abort()
			return
		}

		c.Set("jwt", accessToken)
		c.Set("user_claims", claims)
		c.Next()
	}
}

func GetUserClaim(c *gin.Context) models.JwtUserClaim {
	claims := c.MustGet("user_claims").(jwt.MapClaims)

	fmt.Printf("%+v\n", claims)
	fmt.Printf("%#v\n", claims)

	userID := uint(claims["user_id"].(float64))
	userName := claims["username"].(string)
	email := claims["email"].(string)
	roleAction := claims["role_action"].(string)
	level := claims["level"].(string)
	telephone := claims["telephone"].(string)
	opId := claims["op_id"]

	user := models.JwtUserClaim{
		UserId:     userID,
		Telephone:  telephone,
		Username:   userName,
		Email:      email,
		Level:      level,
		RoleAction: roleAction,
	}

	if opId != nil {
		user.OpId = opId.(*uint)
	}

	return user
}
