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
	NotificationStatusPartial = "partial"
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

type NotificationDelivery struct {
	Id             string        `bson:"_id" json:"id"`
	NotificationId string        `bson:"notificationId" json:"notificationId"`
	Subject        string        `bson:"subject" json:"subject"`
	Body           string        `bson:"body" json:"body"`
	Type           string        `bson:"type" json:"type"`
	Channels       []string      `bson:"channels" json:"channels"`
	Recipients     []string      `bson:"recipients" json:"recipients" binding:"required"`
	Status         string        `bson:"status" json:"status"`
	Message        string        `bson:"errorMessage" json:"errorMessage"`
	ExpireTime     time.Duration `bson:"expireTime" json:"expireTime"`
	CreatedDate    time.Time     `bson:"createdDate" json:"createdDate"`
	UpdatedDate    time.Time     `bson:"updatedDate" json:"updatedDate"`
}

func (e *NotificationDelivery) GetId() string {
	return e.Id
}

func (e *NotificationDelivery) SetId(id string) {
	e.Id = id
}

type NotificationRecipient struct {
	Id             string    `bson:"_id" json:"id"`
	NotificationId string    `bson:"notificationId" json:"notificationId" binding:"required"`
	DeliveryId     string    `bson:"deliveryId" json:"deliveryId" binding:"required"`
	RecipientId    string    `bson:"recipientId" json:"recipientId" binding:"required"`
	Status         string    `bson:"status" json:"status"`
	Message        string    `bson:"errorMessage" json:"errorMessage"`
	CreatedDate    time.Time `bson:"createdDate" json:"createdDate"`
	SentDate       time.Time `bson:"sentDate" json:"sentDate"`
}

type NotificationDetails struct {
}

func (e *NotificationRecipient) GetId() string {
	return e.Id
}

func (e *NotificationRecipient) SetId(id string) {
	e.Id = id
}
