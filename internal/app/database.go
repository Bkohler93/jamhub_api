package app

import (
	"context"
	"database/sql"

	"github.com/bkohler93/jamhubapi/internal/database"
	"github.com/google/uuid"
)

type DB interface {
	GetUserByDisplayName(ctx context.Context, displayName string) (database.User, error)

	CreateUser(ctx context.Context, arg database.CreateUserParams) (database.User, error)

	UpdateUser(ctx context.Context, arg database.UpdateUserParams) (database.User, error)

	GetUserByEmail(ctx context.Context, arg sql.NullString) (database.User, error)

	GetUserByPhone(ctx context.Context, arg sql.NullString) (database.User, error)

	CreateRevokedToken(ctx context.Context, arg database.CreateRevokedTokenParams) (database.RevokedToken, error)

	CreateRoom(ctx context.Context, arg database.CreateRoomParams) (database.Room, error)

	GetRooms(ctx context.Context, limit int32) ([]database.Room, error)

	GetRoomByID(ctx context.Context, uid uuid.UUID) (database.Room, error)

	DeleteRoom(ctx context.Context, uid uuid.UUID) error

	CreatePost(ctx context.Context, arg database.CreatePostParams) (database.Post, error)

	GetPost(ctx context.Context, uid uuid.UUID) (database.Post, error)

	DeletePost(ctx context.Context, uid uuid.UUID) error

	GetRoomPosts(ctx context.Context, uid uuid.UUID) ([]database.Post, error)

	CreateRoomSubscription(ctx context.Context, arg database.CreateRoomSubscriptionParams) (database.RoomSubscription, error)

	DeleteRoomSubscription(ctx context.Context, arg database.DeleteRoomSubscriptionParams) error

	GetAllRoomSubs(ctx context.Context) ([]database.RoomSubscription, error)

	GetRoomRoomSubscriptions(ctx context.Context, uid uuid.UUID) ([]database.RoomSubscription, error)

	GetUserRoomSubscriptions(ctx context.Context, uid uuid.UUID) ([]database.RoomSubscription, error)

	CreatePostVote(ctx context.Context, arg database.CreatePostVoteParams) (database.PostVote, error)

	DeletePostVote(ctx context.Context, arg database.DeletePostVoteParams) error

	GetPostVote(ctx context.Context, arg database.GetPostVoteParams) (database.PostVote, error)

	GetPostPostVotes(ctx context.Context, postID uuid.UUID) ([]database.PostVote, error)

	GetUserRoomsOrderedBySubs(ctx context.Context, arg database.GetUserRoomsOrderedBySubsParams) ([]database.GetUserRoomsOrderedBySubsRow, error)

	GetRoomsOrderedBySubs(ctx context.Context, arg database.GetRoomsOrderedBySubsParams) ([]database.GetRoomsOrderedBySubsRow, error)

	GetNewRoomPosts(ctx context.Context, uid uuid.UUID) ([]database.Post, error)

	GetTopRoomPosts(ctx context.Context, uid uuid.UUID) ([]database.GetTopRoomPostsRow, error)

	GetUserByID(ctx context.Context, uid uuid.UUID) (database.User, error)

	GetRevokedToken(ctx context.Context, uid uuid.UUID) (database.RevokedToken, error)
}
