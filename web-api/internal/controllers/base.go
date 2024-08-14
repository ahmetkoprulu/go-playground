package controllers

import (
	"github.com/ahmetkoprulu/go-playground/web-api/internal/data"
	repository "github.com/ahmetkoprulu/go-playground/web-api/internal/data/repositories"

	"github.com/gin-gonic/gin"
)

type HttpContext struct {
	*gin.Context
}

func BindModel[T any](ctx *gin.Context) *T {
	var model T
	if err := ctx.BindJSON(&model); err != nil {
		BadRequest(ctx, err.Error())
		return nil
	}

	return &model
}

func GetRepositoryContext() *repository.RepositoryContext {
	var context = repository.RepoContext()
	return context
}

func GetDb() *data.MongoDbContext {
	var db, err = data.Context()
	if err != nil {
		panic(err)
	}

	return db
}

func Ok(ctx *gin.Context, data any) {
	ctx.JSON(200, data)
}

func NotFound(ctx *gin.Context, message string) {
	ctx.JSON(404, gin.H{"error": message})
	ctx.Abort()
}

func BadRequest(ctx *gin.Context, message string) {
	ctx.JSON(400, gin.H{"error": message})
	ctx.Abort()
}

func InternalServerError(ctx *gin.Context, message string) {
	ctx.JSON(500, gin.H{"error": message})
	ctx.Abort()
}
