package repositories

import "github.com/ahmetkoprulu/go-playground/notification-service/data"

var ctx *RepositoryContext = nil

func RepoContext() *RepositoryContext {
	if ctx == nil {
		ctx = &RepositoryContext{}
	}

	return ctx
}

func InitializeRepositoryContext() {
	var dbContext, err = data.Context()
	if err != nil {
		panic(err)
	}

	RepoContext()
	ctx.UserRepository = &UserRepository{DbContext: dbContext}
	ctx.NotificationRepository = &NotificationRepository{DbContext: dbContext}
}

type RepositoryContext struct {
	NotificationRepository *NotificationRepository
	UserRepository         *UserRepository
}
