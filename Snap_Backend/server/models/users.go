package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID       primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Role     string             `bson:"role,omitempty" json:"role"`
	UserName string             `bson:"username,omitempty" json:"username" validate:"required,min=3,max=50"`
	Email    string             `bson:"email,omitempty" json:"email" validate:"required,email"`
	Password string             `bson:"password,omitempty" json:"password" validate:"required,min=6"`
	Phone    string             `bson:"phone,omitempty" json:"phone" validate:"required,min=10"`
	Age      string             `bson:"age,omitempty" json:"age" validate:"required"`
	Location string             `bson:"location,omitempty" json:"location" validate:"required"`

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
