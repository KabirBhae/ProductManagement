package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Product struct {
	ProductID primitive.ObjectID `json:"productID" bson:"_id"`
	Name      string             `json:"name" validate:"required"`
	Price     float32            `json:"price" validate:"required"`
	Quantity  int                `json:"quantity" validate:"required"`

	SellerID       primitive.ObjectID `json:"sellerID"`
	SellerUserName string             `json:"sellerUserName"`
}
