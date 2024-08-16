package controllers

import (
	"github.com/ahmetkoprulu/go-playground/web-api/internal/helpers"
	middlewares "github.com/ahmetkoprulu/go-playground/web-api/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func SetupAccountRouter(router *gin.Engine) {
	router.POST("/sign-in", signin)
	router.POST("/sign-up", signup)

	protected := router.Group("/", middlewares.AuthMiddleware())
	{
		protected.GET("/me", getMe)
		// protected.GET("/all", getUsers)
	}
}

func signin(c *gin.Context) {
	var model = BindModel[SignInModel](c)
	if model == nil {
		return
	}

	var repoContext = GetRepositoryContext()
	var exist, existError = repoContext.UserRepository.GetByEmail(model.Email)
	if existError == nil {
		BadRequest(c, "Invalid credentials")
		return
	}

	var pHash = helpers.HashPassword(model.Password, exist.Id)
	if pHash != exist.Password {
		BadRequest(c, "Invalid credentials")
		return
	}

	var token, err = helpers.GenerateJwtToken(exist.Username)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Ok(c, &AuthenticatedUser{exist.Username, exist.Email, token})
}

func signup(c *gin.Context) {
	var model = BindModel[SignUpModel](c)
	if model == nil {
		return
	}

	var repoContext = GetRepositoryContext()
	var exist, _ = repoContext.UserRepository.GetByEmail(model.Email)
	if exist != nil {
		BadRequest(c, "The email is taken")
		return
	}

	var _, err = repoContext.UserRepository.Register(model.Username, model.Email, model.Password)
	if err != nil {
		InternalServerError(c, err.Error())
		return
	}

	Ok(c, nil)
}

func getMe(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "getMe",
	})
}

type SignInModel struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type SignUpModel struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthenticatedUser struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	AccessToken string `json:"accessToken"`
}
