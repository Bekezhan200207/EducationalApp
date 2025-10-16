package models

import (
	"time"

	"github.com/google/uuid"
)

type Session struct {
	ID           int       `json:"id" db:"id"`
	UserUUID     uuid.UUID `json:"user_uuid" db:"user_uuid"`
	RefreshToken string    `json:"refresh_token" db:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
}
