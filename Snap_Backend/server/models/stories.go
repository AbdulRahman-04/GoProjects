package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Story struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	UserId primitive.ObjectID `bson:"userId,omitempty" json:"userId"`
	ImageUrl string `bson:"imageUrl,omitempty" json:"imageUrl" validate:"required"`
	Text string `bson:"storytext,omitempty" json:"storytext" validate:"required,max=100"`
	Song string `bson:"storysong,omitempty" json:"storysong" validate:"required"`
	FilterName string `bson:"filtername,omitempty" json:"filtername" validate:"required,oneof=Dark Bright Blueish Aesthetic"`
	Location string `bson:"location,omitempty" json:"location" validate:"required"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}