package models

import "github.com/google/uuid"

type User struct {
	UUID         uuid.UUID
	Name         string
	Surname      string
	User_Type    string
	Email        string
	PasswordHash string
}
