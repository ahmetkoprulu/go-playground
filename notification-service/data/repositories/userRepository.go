package repositories

import (
	"github.com/ahmetkoprulu/go-playground/notification-service/data"
	"github.com/ahmetkoprulu/go-playground/notification-service/models"
	"go.mongodb.org/mongo-driver/bson"
)

type UserRepository struct {
	DbContext *data.MongoDbContext
}

func (repo *UserRepository) GetById(id string) (*models.User, error) {
	return repo.DbContext.Users().FirstOrDefault(bson.M{"_id": id})
}
