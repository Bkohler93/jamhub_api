package app

import (
	"github.com/bkohler93/jamhubapi/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"password_hash"`
	DisplayName  string    `json:"display_name"`
}

func databaseUsertoUser(dbU database.User) User {
	return User{
		ID:           dbU.ID,
		Email:        dbU.Email.String,
		Phone:        dbU.Phone.String,
		PasswordHash: dbU.PasswordHash,
		DisplayName:  dbU.DisplayName,
	}
}
