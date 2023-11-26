package app

// auth_test.go

import (
	"context"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/bkohler93/jamhubapi/internal/database"
	"github.com/google/uuid"
)

const (
	tokenSecret = "super-secret"
)

type MockDB struct {
}

// CreatePost implements DB.
func (MockDB) CreatePost(ctx context.Context, arg database.CreatePostParams) (database.Post, error) {
	return database.Post{}, nil
}

// CreatePostVote implements DB.
func (MockDB) CreatePostVote(ctx context.Context, arg database.CreatePostVoteParams) (database.PostVote, error) {
	return database.PostVote{}, nil
}

// CreateRevokedToken implements DB.
func (MockDB) CreateRevokedToken(ctx context.Context, arg database.CreateRevokedTokenParams) (database.RevokedToken, error) {
	return database.RevokedToken{}, nil
}

// CreateRoom implements DB.
func (MockDB) CreateRoom(ctx context.Context, arg database.CreateRoomParams) (database.Room, error) {
	return database.Room{}, nil
}

// CreateRoomSubscription implements DB.
func (MockDB) CreateRoomSubscription(ctx context.Context, arg database.CreateRoomSubscriptionParams) (database.RoomSubscription, error) {
	return database.RoomSubscription{}, nil
}

// CreateUser implements DB.
func (MockDB) CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error) {
	return database.User{}, nil
}

// DeletePost implements DB.
func (MockDB) DeletePost(ctx context.Context, uid uuid.UUID) error {
	return nil
}

// DeletePostVote implements DB.
func (MockDB) DeletePostVote(ctx context.Context, uid uuid.UUID) error {
	return nil
}

// DeleteRoom implements DB.
func (MockDB) DeleteRoom(ctx context.Context, uid uuid.UUID) error {
	return nil
}

// DeleteRoomSubscription implements DB.
func (MockDB) DeleteRoomSubscription(ctx context.Context, arg database.DeleteRoomSubscriptionParams) error {
	return nil
}

// GetAllRoomSubs implements DB.
func (MockDB) GetAllRoomSubs(ctx context.Context) ([]database.RoomSubscription, error) {
	return []database.RoomSubscription{}, nil
}

// GetNewRoomPosts implements DB.
func (MockDB) GetNewRoomPosts(ctx context.Context, uid uuid.UUID) ([]database.Post, error) {
	return []database.Post{}, nil
}

// GetPost implements DB.
func (MockDB) GetPost(ctx context.Context, uid uuid.UUID) (database.Post, error) {
	return database.Post{}, nil
}

// GetRevokedToken implements DB.
func (MockDB) GetRevokedToken(ctx context.Context, uid uuid.UUID) (database.RevokedToken, error) {
	return database.RevokedToken{}, errors.New("no token received")
}

// GetRoomByID implements DB.
func (MockDB) GetRoomByID(ctx context.Context, uid uuid.UUID) (database.Room, error) {
	return database.Room{}, nil
}

// GetRoomPosts implements DB.
func (MockDB) GetRoomPosts(ctx context.Context, uid uuid.UUID) ([]database.Post, error) {
	return []database.Post{}, nil
}

// GetRoomRoomSubscriptions implements DB.
func (MockDB) GetRoomRoomSubscriptions(ctx context.Context, uid uuid.UUID) ([]database.RoomSubscription, error) {
	return []database.RoomSubscription{}, nil
}

// GetRooms implements DB.
func (MockDB) GetRooms(ctx context.Context, limit int32) ([]database.Room, error) {
	return []database.Room{}, nil
}

// GetRoomsOrderedBySubs implements DB.
func (MockDB) GetRoomsOrderedBySubs(ctx context.Context, arg database.GetRoomsOrderedBySubsParams) ([]database.GetRoomsOrderedBySubsRow, error) {
	return []database.GetRoomsOrderedBySubsRow{}, nil
}

// GetTopRoomPosts implements DB.
func (MockDB) GetTopRoomPosts(ctx context.Context, uid uuid.UUID) ([]database.GetTopRoomPostsRow, error) {
	return []database.GetTopRoomPostsRow{}, nil
}

// GetUserByDisplayName implements DB.
func (MockDB) GetUserByDisplayName(ctx context.Context, displayName string) (database.User, error) {
	return database.User{}, errors.New("no user found")
}

// GetUserByEmail implements DB.
func (MockDB) GetUserByEmail(ctx context.Context, arg sql.NullString) (database.User, error) {
	return database.User{}, nil
}

// GetUserByID implements DB.
func (MockDB) GetUserByID(ctx context.Context, uid uuid.UUID) (database.User, error) {
	return database.User{}, nil
}

// GetUserByPhone implements DB.
func (MockDB) GetUserByPhone(ctx context.Context, arg sql.NullString) (database.User, error) {
	return database.User{}, nil
}

// GetUserRoomSubscriptions implements DB.
func (MockDB) GetUserRoomSubscriptions(ctx context.Context, uid uuid.UUID) ([]database.RoomSubscription, error) {
	return []database.RoomSubscription{}, nil
}

// GetUserRoomsOrderedBySubs implements DB.
func (MockDB) GetUserRoomsOrderedBySubs(ctx context.Context, arg database.GetUserRoomsOrderedBySubsParams) ([]database.GetUserRoomsOrderedBySubsRow, error) {
	return []database.GetUserRoomsOrderedBySubsRow{}, nil
}

// UpdateUser implements DB.
func (MockDB) UpdateUser(ctx context.Context, arg database.UpdateUserParams) (database.User, error) {
	return database.User{}, nil
}

func generateTokenSafe(id string, expiration time.Time, issuer string) string {
	token, _ := generateToken(id, tokenSecret, expiration, issuer)
	return token
}

func TestAuthMiddleware(t *testing.T) {
	// Set up a test case
	testCases := []struct {
		name           string
		authToken      string
		expectedStatus int
	}{
		{
			name:           "Invalid Authorization Header",
			authToken:      generateTokenSafe("invalid-id", time.Now().Add(time.Hour), "jamhub_access"),
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "No Token present",
			authToken:      "",
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:           "Valid Authorization Token",
			authToken:      generateTokenSafe(uuid.New().String(), time.Now().Add(time.Minute), "jamhub_access"),
			expectedStatus: http.StatusOK,
		},
		// Add more test cases for different scenarios
	}

	// Initialize the mock config
	mockConfig := NewConfig(MockDB{}, tokenSecret)

	// Iterate through test cases
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the Authorization header
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", "Bearer "+tc.authToken)

			rr := httptest.NewRecorder()

			handler := func(w http.ResponseWriter, r *http.Request, user User) {
				// Implement your mock handler logic here
			}
			authHandler := http.HandlerFunc(mockConfig.authMiddleware(handler))

			authHandler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("Expected status %d, but got %d", tc.expectedStatus, rr.Code)
			}
		})
	}
}

func TestAuthRefreshMiddleware(t *testing.T) {
	testCases := []struct {
		name         string
		authToken    string
		expectedCode int
	}{
		{
			name:         "Valid Token results in Status OK",
			authToken:    generateTokenSafe(uuid.New().String(), time.Now().Add(time.Hour), "jamhub_refresh"),
			expectedCode: 200,
		},
		{
			name:         "Access token results in Unauthorized",
			authToken:    generateTokenSafe(uuid.New().String(), time.Now().Add(time.Hour), "jamhub_access"),
			expectedCode: 401,
		},
		{
			name:         "Expired token results in Unauthorized",
			authToken:    generateTokenSafe(uuid.New().String(), time.UnixMilli(0), "jamhub_refresh"),
			expectedCode: 401,
		},
	}

	mockConfig := NewConfig(MockDB{}, tokenSecret)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request with the Authorization header
			req, err := http.NewRequest("GET", "/", nil)
			if err != nil {
				t.Fatal(err)
			}

			req.Header.Set("Authorization", "Bearer "+tc.authToken)

			rr := httptest.NewRecorder()

			handler := func(w http.ResponseWriter, r *http.Request, user User) {
				// Implement your mock handler logic here
			}
			authHandler := http.HandlerFunc(mockConfig.authRefreshMiddleware(handler))

			authHandler.ServeHTTP(rr, req)

			if rr.Code != tc.expectedCode {
				t.Errorf("Expected status %d, but got %d", tc.expectedCode, rr.Code)
			}
		})
	}
}
