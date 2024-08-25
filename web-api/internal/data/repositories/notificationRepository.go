package repositories

import (
	"time"

	"github.com/ahmetkoprulu/go-playground/web-api/internal/data"
	data_models "github.com/ahmetkoprulu/go-playground/web-api/internal/data/abstract"
	"github.com/ahmetkoprulu/go-playground/web-api/internal/models"
	"go.mongodb.org/mongo-driver/bson"
)

type NotificationRepository struct {
	DbContext *data.MongoDbContext
}

func (repo *NotificationRepository) Save(model *models.Notification) (*models.Notification, error) {
	if model.CreatedDate.IsZero() {
		model.CreatedDate = time.Now().UTC()
	}

	var _, err = repo.DbContext.Notifications().Upsert(model)
	return model, err
}

func (repo *NotificationRepository) GetById(id string) (*models.Notification, error) {
	return repo.DbContext.Notifications().FirstOrDefault(bson.M{"_id": id})
}

func (repo *NotificationRepository) GetAll() ([]*models.Notification, error) {
	return repo.DbContext.Notifications().Where(bson.M{})
}

func (repo *NotificationRepository) Paginate(filter string, size, page int) (data_models.PagingModel[*models.Notification], error) {
	return repo.DbContext.Notifications().Paginate(filter, size, page)
}
