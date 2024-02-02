package app

import (
	"database/sql"
	"testing"
	"time"

	"github.com/bkohler93/jamhubapi/internal/database"
	"github.com/google/uuid"
)

func TestUserModel(t *testing.T) {
	uid := uuid.New()
	email := "new-email"
	pwHash := "hashed-password"
	phone := "555-555-5555"
	displayName := "display-name"
	createdAt := time.Now()
	updatedAt := time.Now()

	expectedU := User{
		ID:           uid,
		Email:        email,
		Phone:        phone,
		PasswordHash: pwHash,
		DisplayName:  displayName,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	dbU := database.User{
		ID: uid,
		Email: sql.NullString{
			String: email,
			Valid:  true,
		},
		Phone: sql.NullString{
			String: phone,
			Valid:  true,
		},
		PasswordHash: pwHash,
		DisplayName:  displayName,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	u := databaseUsertoUser(dbU)

	if u != expectedU {
		t.Errorf("Expected user to equal test user")
	}
}

func TestRoomModel(t *testing.T) {
	id := uuid.New()
	name := "test"
	createdAt := time.Now()
	updatedAt := time.Now()

	expectedRoom := Room{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	dbR := database.Room{
		ID:        id,
		Name:      name,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	r := databaseRoomToRoom(dbR)

	if r != expectedRoom {
		t.Errorf("Expected new room to equal test room")
	}
}

func TestPostModel(t *testing.T) {
	id := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	userID := uuid.New()
	roomID := uuid.New()
	link := "link"

	expectedPost := Post{
		ID:        id,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		UserID:    userID,
		RoomID:    roomID,
		Link:      link,
	}

	dbP := database.Post{
		ID:        id,
		UserID:    userID,
		RoomID:    roomID,
		Link:      link,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	p := databasePostToPost(dbP)

	if p != expectedPost {
		t.Errorf("expected new post to equal test post")
	}
}

func TestRoomSubModel(t *testing.T) {
	roomID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()

	expectedRoomSub := RoomSubscription{
		RoomID:    roomID,
		UserID:    userID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	dbRS := database.RoomSubscription{
		RoomID:    roomID,
		UserID:    userID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	rs := databaseRoomSubscriptionToRoomSubscription(dbRS)

	if rs != expectedRoomSub {
		t.Errorf("expected new room sub to equal test room sub")
	}
}

func TestPostVoteModel(t *testing.T) {
	ID := uuid.New()
	postID := uuid.New()
	userID := uuid.New()
	createdAt := time.Now()
	updatedAt := time.Now()
	isUp := true

	expectedPostVote := PostVote{
		ID:        ID,
		PostID:    postID,
		UserID:    userID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		IsUp:      isUp,
	}

	dbPostVote := database.PostVote{
		ID:        ID,
		PostID:    postID,
		UserID:    userID,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
		IsUp:      isUp,
	}

	pv := databasePostVoteToPostVote(dbPostVote)

	if pv != expectedPostVote {
		t.Errorf("expected new post vote to equal test post vote")
	}
}
