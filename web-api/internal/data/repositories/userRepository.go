package repositories

import (
	"time"

	"github.com/ahmetkoprulu/go-playground/web-api/internal/data"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/helpers"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserRepository struct {
	DbContext *data.MongoDbContext
}

func (repo *UserRepository) Register(username, email, password string) (*models.User, error) {
	var id = primitive.NewObjectID().Hex()
	user := &models.User{
		Id:          id,
		Username:    username,
		Email:       email,
		Password:    helpers.HashPassword(password, id),
		CreatedDate: time.Now().UTC(),
	}

	var _, err = repo.DbContext.Users().Upsert(user)
	return user, err
}

func (repo *UserRepository) GetById(id string) (*models.User, error) {
	return repo.DbContext.Users().FirstOrDefault(bson.M{"_id": id})
}

func (repo *UserRepository) GetByEmail(email string) (*models.User, error) {
	return repo.DbContext.Users().FirstOrDefault(bson.M{"email": email})
}
