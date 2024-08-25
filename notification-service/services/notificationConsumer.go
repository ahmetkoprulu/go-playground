package services

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/ahmetkoprulu/go-playground/notification-service/data/repositories"
	"github.com/ahmetkoprulu/go-playground/notification-service/models"
	"github.com/streadway/amqp"
)

func (s *NotificationService) Start(ctx context.Context) {
	// Channels to communicate between workers, consumer and the main goroutine
	messages := make(chan NotificationMessage, s.Workers) // consumer send notifications to message channel, workers process them and send errors to main channel to be handled by the main goroutine
	errors := make(chan error, s.Workers)
	repo := repositories.RepoContext()

	for w := 1; w <= s.Workers; w++ { // Workers that process consumed notifications and send them to the recipients
		go s.work(ctx, repo, messages, errors)
	}

	go s.ConsumeNotifications(ctx, messages) // Listen for notifications until the context is cancelled

	// Handle results or errors
	for res := range errors { // Retry mechanisms, alerting, logging to external systems
		if res == nil {
			continue
		}

		log.Println("Error processing notification:", res)
	}
}

func (s *NotificationService) ConsumeNotifications(ctx context.Context, messages chan NotificationMessage) {
	for {
		select {
		case <-ctx.Done(): // Cancel the context and close the channel to gradually stop the service
			close(messages)
			return
		default:
			if len(messages) >= cap(messages) { // Limit the number of messages in the channel to avoid memory issues and dead locks
				time.Sleep(1 * time.Second)
				continue
			}

			msg, err := s.Consumer.Consume() // Consume notifications from the queue
			if err != nil {
				log.Println("Failed to consume message:", err)
				continue
			}

			notification := parseMessage(msg) // Parse message into NotificationRecipient
			messages <- notification
		}
	}
}

func (s *NotificationService) work(ctx context.Context, repo *repositories.RepositoryContext, messages <-chan NotificationMessage, results chan<- error) {
	for message := range messages {
		notificationDelivery, err := repo.NotificationRepository.GetDeliveryById(message.DeliveryId)
		log.Println("Processing Delivery:", notificationDelivery.Id)
		deliveryStatus := models.NotificationStatusSent
		if err != nil {
			results <- err
			continue
		}

		for _, recipient := range notificationDelivery.Recipients {
			err := s.sendNotification(ctx, recipient, notificationDelivery)
			var recipientStatus = models.NotificationStatusSent
			var message = ""
			if err != nil {
				deliveryStatus = models.NotificationStatusPartial
				recipientStatus = models.NotificationStatusFailed
				message = err.Error()
			}

			notificationRecipient := &models.NotificationRecipient{
				NotificationId: notificationDelivery.NotificationId,
				DeliveryId:     notificationDelivery.Id,
				RecipientId:    recipient,
				Status:         recipientStatus,
				CreatedDate:    time.Now().UTC(),
				Message:        message,
			}

			_, err = repo.NotificationRepository.DbContext.NotificationRecipients().Upsert(notificationRecipient)
			if err != nil {
				results <- err
			}
		}

		notificationDelivery.Status = deliveryStatus
		notificationDelivery.UpdatedDate = time.Now().UTC()
		_, err = repo.NotificationRepository.DbContext.NotificationDeliveries().Upsert(notificationDelivery)
		if err != nil {
			results <- err
			continue
		}
	}
}

func (s *NotificationService) sendNotification(ctx context.Context, recipient string, notification *models.NotificationDelivery) error {
	// Send notification to the recipient
	time.Sleep(100 * time.Millisecond) // Simulate sending notification
	log.Println("Sending notification to recipient:", recipient, "with message:", notification.Message)

	return nil
}

func parseMessage(msg []byte) NotificationMessage {
	// Implement parsing logic here
	var notification NotificationMessage
	err := json.Unmarshal(msg, &notification)
	if err != nil {
		log.Println("Failed to parse message:", err)
		return NotificationMessage{}
	}

	return notification
}

type IConsumer interface {
	Consume() ([]byte, error)
}

type RabbitMqConsumer struct {
	Channel   *amqp.Channel
	QueueName string
}

func NewRabbitMQConsumer(channel *amqp.Channel, queue string) *RabbitMqConsumer {
	return &RabbitMqConsumer{
		Channel:   channel,
		QueueName: queue,
	}
}

func (c *RabbitMqConsumer) Consume() ([]byte, error) {
	msgs, err := c.Channel.Consume(c.QueueName, "", true, false, false, false, nil) // Queue name, consumer name, auto-ack, exclusive, no-local, no-wait, args
	if err != nil {
		log.Println("Failed to consume messages:", err)
		return nil, err
	}

	for msg := range msgs {
		return msg.Body, nil
	}

	return nil, nil
}

type NotificationService struct {
	Repo     repositories.NotificationRepository
	Consumer IConsumer
	Workers  int
}

func NewNotificationService(repo repositories.NotificationRepository, consumer IConsumer, workers int) *NotificationService {
	return &NotificationService{
		Repo:     repo,
		Consumer: consumer,
		Workers:  workers,
	}
}
