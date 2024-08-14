package models

type User struct {
	Id       string `bson:"_id" json:"id"`
	Username string `bson:"username" json:"username" binding:"required"`
	Email    string `bson:"email" json:"email" binding:"required,email"`
	Password string `bson:"password" json:"password" binding:"required"`
}

func (e *User) GetId() string {
	return e.Id
}

func (e *User) SetId(id string) {
	e.Id = id
}

type IEntity interface {
	GetId() string
	SetId(id string)
}
