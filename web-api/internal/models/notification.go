package models

import "time"

const (
	NotificationTypeMessage = "message"
	NotificationTypeAlert   = "alert"
)

const (
	NotificationChannelEmail = "email"
	NotificationChannelSms   = "sms"
	NotificationChannelPush  = "push"
)

type Notification struct {
	Id          string    `bson:"_id" json:"id"`
	Label       string    `bson:"label" json:"label" binding:"required"`
	Subject     string    `bson:"subject" json:"subject" binding:"required"`
	Body        string    `bson:"body" json:"body" binding:"required"`
	Channels    []string  `bson:"channels" json:"channels" binding:"required"`
	Type        string    `bson:"type" json:"type" binding:"required"`
	CreatedDate time.Time `bson:"createdDate" json:"createdDate"`
}

func (e *Notification) GetId() string {
	return e.Id
}

func (e *Notification) SetId(id string) {
	e.Id = id
}
