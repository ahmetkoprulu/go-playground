package controllers

import (
	"strconv"

	middlewares "github.com/ahmetkoprulu/go-playground/web-api/internal/middlewares"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/models"
	"github.com/gin-gonic/gin"
)

func SetupNotificationRouter(router *gin.Engine) {
	protected := router.Group("/notifications", middlewares.AuthMiddleware())
	{
		protected.POST("/save", saveNotification)
		protected.GET("/all", getNotifications)
		protected.GET("/page", paginateNotifications)
		protected.GET("/:id", getNotification)
	}
}

func saveNotification(c *gin.Context) {
	var model = BindModel[models.Notification](c)
	if model == nil {
		return
	}

	var repoContext = GetRepositoryContext()
	var notification, err = repoContext.NotificationRepository.Save(model)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Ok(c, notification)
}

func getNotifications(c *gin.Context) {
	var repoContext = GetRepositoryContext()
	notifications, err := repoContext.NotificationRepository.GetAll()
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Ok(c, notifications)
}

func getNotification(c *gin.Context) {
	var id = c.Param("id")
	var repoContext = GetRepositoryContext()
	notification, err := repoContext.NotificationRepository.GetById(id)
	if err != nil {
		NotFound(c, err.Error())
		return
	}

	Ok(c, notification)
}

func paginateNotifications(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	takeStr := c.DefaultQuery("take", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		BadRequest(c, "Invalid 'page' parameter")
		return
	}

	take, err := strconv.Atoi(takeStr)
	if err != nil || take < 1 {
		BadRequest(c, "Invalid 'take' parameter")
		return
	}

	var repoContext = GetRepositoryContext()
	notifications, err := repoContext.NotificationRepository.DbContext.Notifications().Paginate(nil, page, take)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Ok(c, notifications)
}
