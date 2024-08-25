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

const (
	NotificationStatusPending = "pending"
	NotificationStatusSent    = "sent"
	NotificationStatusFailed  = "failed"
)

type Notification struct {
	Id          string        `bson:"_id" json:"id"`
	Label       string        `bson:"label" json:"label" binding:"required"`
	Subject     string        `bson:"subject" json:"subject" binding:"required"`
	Body        string        `bson:"body" json:"body" binding:"required"`
	Schedule    string        `bson:"schedule" json:"schedule" binding:"required"` // cron expression
	Channels    []string      `bson:"channels" json:"channels" binding:"required"`
	Recipients  []string      `bson:"recipients" json:"recipients" binding:"required"`
	Type        string        `bson:"type" json:"type" binding:"required"`
	ExpireTime  time.Duration `bson:"expireTime" json:"expireTime"`
	CreatedDate time.Time     `bson:"createdDate" json:"createdDate"`
}

func (e *Notification) GetId() string {
	return e.Id
}

func (e *Notification) SetId(id string) {
	e.Id = id
}

type NotificationRecipient struct {
	Id             string              `bson:"_id" json:"id"`
	NotificationId string              `bson:"notificationId" json:"notificationId" binding:"required"`
	Details        NotificationDetails `bson:"details" json:"details" binding:"required"`
	RecipientId    string              `bson:"recipientId" json:"recipientId" binding:"required"`
	Status         string              `bson:"status" json:"status"`
	ExpireTime     time.Duration       `bson:"expireTime" json:"expireTime"`
	CreatedDate    time.Time           `bson:"createdDate" json:"createdDate"`
}

type NotificationDetails struct {
	Subject string `bson:"subject" json:"subject"`
	Body    string `bson:"body" json:"body"`
	Type    string `bson:"type" json:"type"`
}

func (e *NotificationRecipient) GetId() string {
	return e.Id
}

func (e *NotificationRecipient) SetId(id string) {
	e.Id = id
}
