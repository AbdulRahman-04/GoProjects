package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Restro struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Role string `bson:"role" json:"role"`
	RestroName string `bson:"restroname" json:"restroname" binding:"required,min=5,max=30"`
	RestroLocation string `bson:"restrolocation" json:"restrolocation" binding:"required"`
	OwnerName string `bson:"ownername" json:"ownername" binding:"required"`
	OwnerEmail string `bson:"owneremail" json:"owneremail" binding:"required"`
	Password string   `bson:"password" json:"password" binding:"required,min=6"`
	Phone string `bson:"phone" json:"phone" binding:"required,min=10,numeric"`
	RestroLicenseNumber string `bson:"licensenumber" json:"licensenumber" binding:"required"`
	OwnerVerified struct {
		Email bool `bson:"emailVerified" json:"emailVerified"`
	} `bson:"ownerVerified" json:"ownerVerified"`
	OwnerVerifyToken struct {
		Email string `bson:"emailVerifyToken" json:"emailVerifyToken"`
	} `bson:"ownerVerifyToken" json:"ownerVerifyToken"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}