package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Menu struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RestroId primitive.ObjectID `bson:"restroId,omitempty" json:"restroId"`
	MenuName string `bson:"menuname" json:"menuname" binding:"required,min=5,max=25"`
	MenuItems []MenuItem `bson:"menuitems" json:"menuitems" binding:"required"`
	ImageUrl string `bson:"imageurl" json:"imageurl" binding:"required"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}

type MenuItem struct {
	ItemsList int `bson:"itemslist" json:"itemslist" binding:"required"`
	ItemsName string `bson:"itemsname" json:"itemsname" binding:"required"`
	ItemsPrice float64 `bson:"itemsprice" json:"itemsprice" binding:"required"`
}