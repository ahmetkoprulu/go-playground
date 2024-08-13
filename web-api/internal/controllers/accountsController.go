package account_router

import (
	middlewares "github.com/ahmetkoprulu/go-playground/web-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupRouter(router *gin.Engine) {
	router.POST("/sign-in", signin)
	router.POST("/sign-up", signup)

	protected := router.Group("/", middlewares.AuthMiddleware())
	{
		protected.GET("/me", getMe)
		protected.GET("/all", getUsers)
	}
}

func signin(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "signin",
	})
}

func signup(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "signup",
	})
}

func getMe(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "getMe",
	})
}

func getUsers(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "getUsers",
	})
}
