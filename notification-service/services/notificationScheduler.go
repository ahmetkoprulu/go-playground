package services

import (
	"context"
	"log"
	"time"

	"github.com/ahmetkoprulu/go-playground/notification-service/data/repositories"
	"github.com/ahmetkoprulu/go-playground/notification-service/models"
)

type NotificationScheduler struct {
	Producer   *NotificationProducer
	Repository repositories.RepositoryContext
	Interval   time.Duration
}

func NewNotificationScheduler(producer *NotificationProducer, repo repositories.RepositoryContext) *NotificationScheduler {
	return &NotificationScheduler{
		Producer:   producer,
		Repository: repo,
		Interval:   time.Duration(5) * time.Second,
	}
}

func (s *NotificationScheduler) Start(ctx context.Context) {

	ticker := time.NewTicker(s.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.processNotifications(ctx)
		}
	}
}

func (s *NotificationScheduler) processNotifications(ctx context.Context) {
	// Attempt to acquire the lock
	locked, err := s.acquireLock(ctx)
	if err != nil {
		log.Println("Error acquiring lock:", err)
		return
	}

	if !locked {
		log.Println("Another instance is already processing notifications.")
		return
	}

	notifications, err := s.Repository.NotificationRepository.GetScheduledNotifications()
	if err != nil {
		log.Printf("Error fetching scheduled notifications: %v", err)
		return
	}

	for _, notification := range notifications {
		s.scheduleNotification(notification)
	}

	// Release the lock after processing
	err = s.releaseLock(ctx)
	if err != nil {
		log.Println("Error releasing lock:", err)
	}
}

func (s *NotificationScheduler) acquireLock(ctx context.Context) (bool, error) {
	return false, nil
}

func (s *NotificationScheduler) releaseLock(ctx context.Context) error {
	return nil
}

func (s *NotificationScheduler) scheduleNotification(notification models.Notification) {
	// check if the notification crong job is passed next occurence
	err := s.Producer.Publish(notification)
	if err != nil {
		log.Printf("Error publishing notification: %v", err)
	}

	if err != nil {
		log.Printf("Error scheduling notification: %v", err)
	}
}
