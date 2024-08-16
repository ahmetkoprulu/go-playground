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
	initDb()

	var port = os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r := gin.Default()
	r.Use(middlewares.ErrorMiddleware())
	controllers.SetupAccountRouter(r)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

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
