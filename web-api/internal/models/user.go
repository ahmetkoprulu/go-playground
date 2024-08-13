package models

type User struct {
	Entity
	Username string `bson:"username" json:"username"`
	Email    string `bson:"email" json:"email"`
	Password string `bson:"password" json:"password"`
}

func (e *Entity) GetId() string {
	return e.Id
}

func (e *Entity) SetId(id string) {
	e.Id = id
}

type Entity struct {
	Id string `bson:"_id" json:"_id"`
}

type IEntity interface {
	GetId() string
	SetId(id string)
}
