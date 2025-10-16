package models

import "github.com/google/uuid"

type User struct {
	UUID         uuid.UUID
	Name         string
	Surname      string
	Role_id      int
	Email        string
	PasswordHash string
}
