package models

import "github.com/google/uuid"

type User struct {
	Id           uuid.UUID
	User_Name    string
	User_Surname string
	User_Type    string
	Email        string
	PasswordHash string
}
