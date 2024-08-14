package repositories

import (
	"github.com/ahmetkoprulu/go-playground/web-api/internal/data"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/helpers"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	DbContext *data.MongoDbContext
}

func (repo *UserRepository) Register(username, email, password string) (*models.User, error) {
	var id = primitive.NewObjectID().String()
	user := &models.User{
		Id:       id,
		Username: username,
		Email:    email,
		Password: helpers.HashPassword(password, id),
	}

	var _, err = repo.DbContext.Users().Upsert(user)
	return user, err
}

func (repo *UserRepository) GetByEmail(email string) *models.User {
	// return repo.DbContext.Users().FirstOrDefault(bson.M{"email": email})
	return nil
}
