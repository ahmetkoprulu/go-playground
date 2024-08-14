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

	context.UserRepository = &UserRepository{DbContext: dbContext}
}

type RepositoryContext struct {
	UserRepository *UserRepository
}
