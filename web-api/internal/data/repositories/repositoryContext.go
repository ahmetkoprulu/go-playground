package repositories

import "github.com/ahmetkoprulu/go-playground/web-api/internal/data"

var context *RepositoryContext = nil

func RepoContext() *RepositoryContext {
	if context == nil {
		context = &RepositoryContext{}
	}

	return context
}

func InitializeRepositoryContext() {
	var dbContext, err = data.Context()
	if err != nil {
		panic(err)
	}

	RepoContext()
	context.UserRepository = &UserRepository{DbContext: dbContext}
	context.NotificationRepository = &NotificationRepository{DbContext: dbContext}
}

type RepositoryContext struct {
	NotificationRepository *NotificationRepository
	UserRepository         *UserRepository
}
