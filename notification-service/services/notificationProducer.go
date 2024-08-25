package services

import (
	"encoding/json"
	"log"

	"github.com/ahmetkoprulu/go-playground/notification-service/models"
	"github.com/streadway/amqp"
)

type NotificationProducer struct {
	Connection *amqp.Connection
}

func NewNotificationProducer(conn *amqp.Connection) *NotificationProducer {
	return &NotificationProducer{
		Connection: conn,
	}
}

func (p *NotificationProducer) Publish(notification models.Notification) error {
	channel, err := p.Connection.Channel()
	if err != nil {
		return err
	}

	defer channel.Close()
	queue, err := p.declareQueue(channel)
	if err != nil {
		return err
	}

	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	if err != nil {
		return err
	}

	log.Printf("Published notification: %s", notification.Id)
	return nil
}

func (p *NotificationProducer) declareQueue(channel *amqp.Channel) (*amqp.Queue, error) {
	queue, err := channel.QueueDeclare(
		"notification_queue", // name
		true,                 // durable
		false,                // delete when unused
		false,                // exclusive
		false,                // no-wait
		nil,                  // arguments
	)

	if err != nil {
		return nil, err
	}

	return &queue, nil
}
