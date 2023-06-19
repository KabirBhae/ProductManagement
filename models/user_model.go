package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	ID       primitive.ObjectID `bson:"_id"`
	Name     string             `json:"name,omitempty" validate:"required"`
	Username string             `json:"username" validate:"required"`
	Email    string             `json:"email,omitempty" validate:"required"`
	Password string             `json:"password,omitempty" validate:"required"`
	Balance  float32            `json:"balance"`
	Status   string             `json:"status"`
	IsAdmin  bool               `json:"isAdmin"`
}
