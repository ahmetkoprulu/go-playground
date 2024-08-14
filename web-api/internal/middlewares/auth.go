package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/ahmetkoprulu/go-playground/web-api/internal/helpers"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := getToken(c)
		if err != nil {
			unauthorized(c, err.Error())
			return
		}

		claims, err := helpers.ValidateJwtToken(token)
		if err != nil {
			unauthorized(c, err.Error())
			return
		}

		c.Set("username", claims.Username)
		c.Next()
	}
}

func getToken(c *gin.Context) (string, error) {
	token := c.GetHeader("Authorization")
	if token == "" {
		return "", errors.New("authorization token is required")
	}

	tokenParts := strings.Split(token, " ")
	if len(tokenParts) != 2 {
		return "", errors.New("invalid authorization token format")
	}

	return tokenParts[1], nil
}

func unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, gin.H{"error": message})
	c.Abort()
}
