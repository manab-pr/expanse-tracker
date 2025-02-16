package models

import "time"

type PasswordReset struct {
	Email     string    `bson:"email"`
	Token     string    `bson:"token"`
	ExpiresAt time.Time `bson:"expires_at"`
}
