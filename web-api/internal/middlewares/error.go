package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ErrorMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer restorePanic(c)
		// Continue to the next middleware/handler
		c.Next()
		// Check if any errors occurred during the request handling
		checkError(c)
	}
}

func checkError(ctx *gin.Context) {
	if len(ctx.Errors) < 1 {
		return
	}

	err := ctx.Errors.Last().Err
	log.Printf("Error: %v", err)
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
	ctx.Abort()
}

func restorePanic(ctx *gin.Context) {
	r := recover()
	if r == nil {
		return
	}

	// Log the panic if necessary
	log.Printf("Panic recovered: %v", r)
	// Respond with a 500 Internal Server Error
	ctx.JSON(http.StatusInternalServerError, gin.H{
		"error": "Internal Server Error",
	})
	ctx.Abort()
}
