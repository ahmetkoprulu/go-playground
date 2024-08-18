package main

import (
	"os"

	"github.com/ahmetkoprulu/go-playground/web-api/internal/controllers"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/data"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/data/repositories"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	initEnv()
	var port = os.Getenv("PORT")
	if port == "" {
		panic("PORT is not set")
	}

	r := setupServer()
	initDb()
	r.Run(":" + port)
}

func initEnv() {
	err := godotenv.Load("../.env")
	if err != nil {
		panic(err)
	}
}

func initDb() {
	err := data.InitializeMongoDb()
	if err != nil {
		panic(err)
	}

	repositories.InitializeRepositoryContext()
}

func setupServer() *gin.Engine {
	r := gin.Default()
	setupMiddlewares(r)
	setupControllers(r)

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r
}

func setupControllers(r *gin.Engine) {
	controllers.SetupAccountRouter(r)
	controllers.SetupNotificationRouter(r)
}

func setupMiddlewares(r *gin.Engine) {
	r.Use(middlewares.ErrorMiddleware())
}
