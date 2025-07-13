package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Order struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId primitive.ObjectID `bson:"userId,omitempty" json:"userId"`
	OrderName string `bson:"ordername" json:"ordername"`
	OrderItems []OrderItem `bson:"orderitems" json:"orderitems" binding:"required"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type OrderItem struct {
	Quantity   int     `bson:"quantity" json:"quantity" binding:"required"`
	ItemsName  string  `bson:"itemsname" json:"itemsname" binding:"required"`
	ItemsPrice float64 `bson:"itemsprice" json:"itemsprice" binding:"required"`
}

