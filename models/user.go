package models

type User struct {
	Id       string `json:"id,omitempty" bson:"_id,omitempty"`
	Name     string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}
