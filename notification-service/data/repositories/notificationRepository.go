package repositories

import (
	"context"
	"time"

	"github.com/ahmetkoprulu/go-playground/notification-service/data"
	"github.com/ahmetkoprulu/go-playground/notification-service/models"
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

func (repo *NotificationRepository) GetDeliveryById(id string) (*models.NotificationDelivery, error) {
	return repo.DbContext.NotificationDeliveries().FirstOrDefault(bson.M{"_id": id})
}

func (repo *NotificationRepository) GetRecipientsByDeliveryId(id string) ([]*models.NotificationRecipient, error) {
	return repo.DbContext.NotificationRecipients().Where(bson.M{"deliveryId": id})
}

func (r *NotificationRepository) GetScheduledNotifications() ([]models.Notification, error) {
	var notifications []models.Notification
	filter := bson.M{"schedule": bson.M{"$exists": true, "$ne": ""}} // Only get scheduled notifications
	cursor, err := r.DbContext.Notifications().Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var notification models.Notification
		if err := cursor.Decode(&notification); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	return notifications, nil
}
