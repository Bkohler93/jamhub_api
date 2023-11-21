package app

import (
	"time"

	"github.com/bkohler93/jamhubapi/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone"`
	PasswordHash string    `json:"password_hash"`
	DisplayName  string    `json:"display_name"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func databaseUsertoUser(dbU database.User) User {
	return User{
		ID:           dbU.ID,
		Email:        dbU.Email.String,
		Phone:        dbU.Phone.String,
		PasswordHash: dbU.PasswordHash,
		DisplayName:  dbU.DisplayName,
		CreatedAt:    dbU.CreatedAt,
		UpdatedAt:    dbU.UpdatedAt,
	}
}

type Room struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseRoomToRoom(dbR database.Room) Room {
	return Room{
		ID:        dbR.ID,
		Name:      dbR.Name,
		CreatedAt: dbR.CreatedAt,
		UpdatedAt: dbR.UpdatedAt,
	}
}

type Post struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	RoomID    uuid.UUID `json:"room_id"`
	Link      string    `json:"link"`
}

func databasePostToPost(dbP database.Post) Post {
	return Post{
		ID:        dbP.ID,
		CreatedAt: dbP.CreatedAt,
		UpdatedAt: dbP.UpdatedAt,
		UserID:    dbP.UserID,
		RoomID:    dbP.RoomID,
		Link:      dbP.Link,
	}
}

type RoomSubscription struct {
	ID        uuid.UUID `json:"id"`
	RoomID    uuid.UUID `json:"room_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func databaseRoomSubscriptionToRoomSubscription(dbRS database.RoomSubscription) RoomSubscription {
	return RoomSubscription{
		ID:        dbRS.ID,
		RoomID:    dbRS.RoomID,
		UserID:    dbRS.UserID,
		CreatedAt: dbRS.CreatedAt,
		UpdatedAt: dbRS.UpdatedAt,
	}
}

type PostVote struct {
	ID        uuid.UUID `json:"id"`
	PostID    uuid.UUID `json:"post_id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	IsUp      bool      `json:"is_up"`
}

func databasePostVoteToPostVote(dbPV database.PostVote) PostVote {
	return PostVote{
		ID:        dbPV.ID,
		PostID:    dbPV.PostID,
		UserID:    dbPV.UserID,
		CreatedAt: dbPV.CreatedAt,
		UpdatedAt: dbPV.UpdatedAt,
		IsUp:      dbPV.IsUp,
	}
}
