package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Role string `bson:"role" json:"role"`
	UserName string `bson:"username" json:"username" binding:"required"`
	Email string  `bson:"email" json:"email" binding:"required"`
	Password string `bson:"password" json:"password" binding:"required,min=6"`
	Phone string `bson:"phone" json:"phone" binding:"required,min=10,numeric"`
	Location string `bson:"location" json:"location" binding:"required,oneof=MehdiPatnam Tolichowki Talabkatta Chadulalbaradari Shahalibanda Sayeedabad Dabeerpura Kalapatther Khilwath Charminar"`
	Language string `bson:"language" json:"language" binding:"required,oneof=Hindi English Urdu Telugu"`
	UserVerified struct {
		Email bool `bson:"emailVerified" json:"emailVerified"`
	} `bson:"userVerified" json:"userVerified"`
	UserVerifyToken struct {
		Email string `bson:"emailVerifyToken" json:"emailVerifyToken"`
		Phone string `bson:"phoneVerifyToken" json:"phoneVerifyToken"`
	} `bson:"userVerifyToken" json:"userVerifyToken"`
	CreatedAt time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time `bson:"updated_at" json:"updated_at"`
}