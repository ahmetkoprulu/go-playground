package main

import (
	"github.com/ahmetkoprulu/go-playground/web-api/internal/data"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/middlewares"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	initEnv()
	initDb()

	r := gin.Default()
	r.Use(middlewares.ErrorMiddleware())
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	r.Run()
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
}
