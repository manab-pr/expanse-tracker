package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID            primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Name          *string            `json:"name" validate:"min=2,max=20"`
	Email         *string            `json:"email" validate:"email,required"`
	Password      *string            `json:"password" validate:"min=5,max=20"`
	Token         *string            `json:"token"`
	Refresh_Token *string            `json:"refresh_token"`
	Created_at    time.Time          `json:"created_at"`
	Updated_at    time.Time          `json:"updated_at"`
	User_Id       string             `json:"user_id"`
}
