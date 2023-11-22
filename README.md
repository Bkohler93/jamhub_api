# JamHub

Welcome to JamHub, the place where you can find and post music you like to whatever niche you're into!

## Development

# API Documentation

## User resource

Stores user data. Email or Phone field must be present for identification purposes.

```
{
  id UUID,
  email string,
  phone string,
  password_hash hashed_password,
  display_name string,
  created_at ISO 8601 standard datetime,
  updated_at ISO 8601 standard datetime
}
```

### POST /v1/users

Creates a new user resource.

Request Body

- email or phone or both must be present

```json
{
  "email": "example@gmail.com",
  "password": "password1234",
  "display_name": "display nme",
  "phone": "555-123-4567"
}
```

Response Body

```json
{
  "id": "000-111-222-333",
  "email": "example@gmail.com",
  "phone": "555-444-4444",
  "display_name": "displayname",
  "updated_at": "2023-01-01 12:15:00",
  "created_at": "2023-01-01 12:15:00"
}
```

Response Codes

- `201` on successful creation of user resource
- `400` on failure to provide proper request body, or user already exists in database
- `500` on failure to generate a hashed password

### PUT /v1/users

Updates a user resource. Requires access token to be present in Authorization header.

Request Body

- updates any of the fields present in request body

```json
{
  "email": "example@gmail.com",
  "phone": "555-111-2222",
  "display_name": "displayname",
  "password": "new password"
}
```

Response body

```json
{
  "id": "111-222-333-444",
  "email": "example@gmail.com",
  "phone": "555-333-4444",
  "display_name": "displayName",
  "created_at": "2023-01-01 12:15:00",
  "updated_at": "2023-01-01 12:15:00"
}
```

Response Codes

- `500` failure to encrypt new password or failure to update resource
- `401` - unauthorized to delete the resource.
- `200` if successful update

### GET /v1/users/rooms/room_subscriptions

Retrieves a list of rooms that the user is subscribed to. Requires access token to be present in Authorization header.

Response body

```json
[
  {
    "room_id": "1235-ffff-ddd-sss",
    "room_name": "Beatles hangout",
    "created_at": "2023-01-01 12:15:00",
    "updated_at": "2023-01-01 12:15:00",
    "subscription_count": 1
  },
  {
    "room_id": "1235-ffff-ddd-sss",
    "room_name": "Beatles hangout",
    "created_at": "2023-01-01 12:15:00",
    "updated_at": "2023-01-01 12:15:00",
    "subscription_count": 1
  }
]
```

Response Codes

- `500` failed to retrieve user's subscribed rooms
- `401` - unauthorized to delete the resource.
- `200` subscribed rooms were successfully retrieved from database

## Room resource

Represents the rooms that users subscribe to. Dedicated to a style of music, an artist, whatever the author chooses.

```
{
	id UUID,
	name string,
	created_at ISO 8601 standard datetime,
	updated_at ISO 8601 standard datetime
}
```

### POST /v1/rooms

Creates a new post resource. Requires access token in Authorization header.

Request body

```json
{
  "name": "Beatles lounge"
}
```

Response body

```json
{
  "id": "111-222-333",
  "name": "Beatles lounge",
  "created_at": "2023-01-01 12:15:00",
  "updated_at": "2023-01-01 12:15:00"
}
```

Response Codes

- `500` error creating a room
- `401` - unauthorized to delete the resource.
- `201` successful creation of resource

### GET /v1/rooms

Retrieves all rooms.

Response body

```json
[
  {
    "id": "111-222-333",
    "name": "Beatles lounge",
    "created_at": "2023-01-01 12:15:00",
    "updated_at": "2023-01-01 12:15:00"
  },
  {
    "id": "111-222-333",
    "name": "Beatles lounge",
    "created_at": "2023-01-01 12:15:00",
    "updated_at": "2023-01-01 12:15:00"
  }
]
```

Status Codes

- `500` could not retrieve rooms
- `200` successful retrieval of rooms

### GET /v1/rooms/{room_id}

Retrieves a room with matching id from URL parameter `room_id`.

URL Parameter `room_id` - must be matching to an existing room.

Response body

```json
{
  "id": "111-222-333",
  "name": "Beatles lounge",
  "created_at": "2023-01-01 12:15:00",
  "updated_at": "2023-01-01 12:15:00"
}
```

Status Codes

- `400` invalid `room_id` in URL
- `404` no resource could be found with `room_id`
- `200` successful retrieval of resource.

### DELETE /v1/rooms/{room_id}

Deletes room resource with matching id as `room_id`.

URL parameter `room_id` - matches an existing room resource.

Status Codes

- `400` - invalid `room_id` in URL.
- `200` - on successful deletion
- `204` - no resource to delete with `room_id`

### GET /v1/rooms/room_subscriptions

Retrieves list of room resources ordered by number of subscribers subscribed to room.

Query Parameters

- limit: `integer` designating number of resources to retrieve
- offset: `integer` designating offset to start retrieving from

Response body

```json
[
  {
    "room_id": "5b5f6c2d-8e9a-4f7c-a3d2-1e2d3f4b5c6d",
    "room_name": "The Beatles",
    "created_at": "2023-01-02T12:30:00Z",
    "updated_at": "2023-01-02T12:30:00Z",
    "subscription_count": 3
  },
  {
    "room_id": "5b5f6c2d-8e9a-4f7c-a3d2-1e2d3f4b5c6d",
    "room_name": "The Who Lounge",
    "created_at": "2023-01-02T12:30:00Z",
    "updated_at": "2023-01-02T12:30:00Z",
    "subscription_count": 1
  }
]
```

Status Codes

- `500` failed to retrieve rooms
- `200` on successful retrieval of rooms

## Post resource

User posts links to songs that are then attached to a room. Posts an be up and down-voted.

```
Post {
	ID        UUID
	UserID    UUID
	RoomID    UUID
	Link      string
	CreatedAt ISO 8601 standard datetime
	UpdatedAt ISO 8601 standard datetime
}
```

### POST /v1/posts

Create new post resource. Requires access token in Authorization header.

Request body

```json
{
  "room_id": "111-222-333",
  "link": "spotify.com/songlink1111"
}
```

Response body

```json
{
  "id": "111-222-333",
  "user_id": "111-22-333",
  "room_id": "222-333-fff",
  "link": "spotify.com/songlink1111",
  "created_at": "2023-01-01 12:15:00",
  "updated_at": "2023-01-01 12:15:00"
}
```

Status Codes

- `404` invalid fields in request body
- `401` unauthorized to delete the resource.
- `201` successful resource creation

### DELETE /v1/posts/{post_id}

Deletes post resource with matching id as `post_id`.

URL parameter `post` - matches an existing post resource.

Status Codes

- `400` invalid `post_id` in URL.
- `404` no post with that id is found
- `401` unauthorized to delete the resource.
- `200` resource deleted

### Get /v1/rooms/posts

Retrieves all posts in a given room

URL parameter `room_id` - matches existing room resource.

Response body

```json
[
  {
    "id": "bfd5a0db-700d-451a-a9ef-0e928619c8ee",
    "created_at": "2023-11-21T16:37:41.243565Z",
    "updated_at": "2023-11-21T16:37:41.243565Z",
    "user_id": "97eb361d-d8a2-4f12-9f71-0d2ca33137a4",
    "room_id": "c56704c7-c3bb-42df-99cc-4e632bbf9f73",
    "link": "https://open.spotify.com/track/1lvpyd1lQjutZa6YnAE8aH?si=60da3179f5ad4dc0"
  },
  {
    "id": "bfd5a0db-700d-451a-a9ef-0e928619c8ee",
    "created_at": "2023-11-21T16:37:41.243565Z",
    "updated_at": "2023-11-21T16:37:41.243565Z",
    "user_id": "97eb361d-d8a2-4f12-9f71-0d2ca33137a4",
    "room_id": "c56704c7-c3bb-42df-99cc-4e632bbf9f73",
    "link": "https://open.spotify.com/track/1lvpyd1lQjutZa6YnAE8aH?si=60da3179f5ad4dc0"
  }
]
```

Status Codes

- `400` invalid or missing `room_id` in request body
- `404` resource with given `room_id` not found
- `200` retrieved posts

### Get /v1/posts/rooms

Retrieves an ordered list of posts in a given room.

URL Query Parameter `select` - can be either "top" or "new". Defaults to "new".

Request body

```json
{
	"room_id": "111-222-333-444
}
```

Response body

```json
[
  {
    "post_id": "3f6aaceb-016e-4d2a-844d-03904ce2938b",
    "room_id": "adb0206f-dde6-4d7e-b517-36da3e899751",
    "link": "https://open.spotify.com/track/1lvpyd1lQjutZa6YnAE8aH?si=60da3179f5ad4dc0",
    "created_at": "2023-11-21T16:55:55.200586Z",
    "updated_at": "2023-11-21T16:55:55.200585Z",
    "num_upvotes": 1
  },
  {
    "post_id": "3f6aaceb-016e-4d2a-844d-03904ce2938b",
    "room_id": "adb0206f-dde6-4d7e-b517-36da3e899751",
    "link": "https://open.spotify.com/track/1lvpyd1lQjutZa6YnAE8aH?si=60da3179f5ad4dc0",
    "created_at": "2023-11-21T16:55:55.200586Z",
    "updated_at": "2023-11-21T16:55:55.200585Z",
    "num_upvotes": 1
  }
]
```

Status Codes

- `500` failed to retrieve rooms
- `400` invalid `room_id` provided

## Room Subscription resource

User's subscribe to a room by creating a Room Subscription resource.

```
Post {
	id UUID
	room_id UUID
	user_id UUID
	created_at ISO 8601 standard datetime
	updated_at ISO 8601 standard datetime
}
```

### POST /v1/room_subs

Creates a new room subscription resource. Requires access token in Authorization header of request.

Request body

```json
{
  "room_id": "111-222-333"
}
```

Response body

```json
{
  "id": "111-222-333",
  "room_id": "111-222-444-fff",
  "user_id": "fff-ggg-hhh",
  "created_at": "2023-01-01 12:15:00",
  "updated_at": "2023-01-01 12:15:00"
}
```

Status Codes

- `400` invalid or missing `room_id`
- `401` unauthorized, missing or invalid access token
- `404` no room exists with `room_id`
- `201` successful creation of resource.

### DELETE /v1/room_subs/{room_sub_id}

Deletes a room subscription resource that matches `room_sub_id`. Requires access token in Authorization header of request.

URL parameter `room_sub_id` must be present in URL.

Status Codes

- `401` unauthorized, missing or invalid access token
- `400` invalid or missing `room_sub_id` query param
- `500` failed to delete resource
- `200` resource deleted successfully

### GET /v1/room_subs

Retrieves all room subscription resources from database.

Response body

```json
{
	{
		"id":"1a2b3c4d-5e6f-7a8b-9c1d-2e3f4a5b6c7d",
		"room_id":"4a4d5e1b-2cf4-4a47-b1f7-6d67e3b8a5f1",
		"user_id":"1f4ab37a-788b-4e9e-8a5e-3a7b6b8bb4b1",
		"created_at":"2023-01-01T12:30:00Z",
		"updated_at":"2023-01-01T12:30:00Z"
	},
	{
		"id":"1a2b3c4d-5e6f-7a8b-9c1d-2e3f4a5b6c7d",
		"room_id":"4a4d5e1b-2cf4-4a47-b1f7-6d67e3b8a5f1",
		"user_id":"1f4ab37a-788b-4e9e-8a5e-3a7b6b8bb4b1",
		"created_at":"2023-01-01T12:30:00Z",
		"updated_at":"2023-01-01T12:30:00Z"
	},
}
```

Status Codes

- `500` failed to retrieve resources
- `200` retrieved resources successfully

### GET /v1/rooms/room_subs/{room_id}

Retrieve all room subscriptions for a specified room.

URL parameter `room_id` used to find room subscriptions for.

Response body

```json
[
  {
    "id": "f0b965c5-d100-48a7-9df6-43375b41f4ed",
    "room_id": "adb0206f-dde6-4d7e-b517-36da3e899751",
    "user_id": "3053cb46-a984-471d-afa3-a14169a21b5a",
    "created_at": "2023-11-21T16:55:57.650871Z",
    "updated_at": "2023-11-21T16:55:57.650871Z"
  },
  {
    "id": "f0b965c5-d100-48a7-9df6-43375b41f4ed",
    "room_id": "adb0206f-dde6-4d7e-b517-36da3e899751",
    "user_id": "3053cb46-a984-471d-afa3-a14169a21b5a",
    "created_at": "2023-11-21T16:55:57.650871Z",
    "updated_at": "2023-11-21T16:55:57.650871Z"
  }
]
```

Status Codes

- `400` invalid or missing `room_id` in URL
- `500` error retrieving room subs for given room
- `200` rooms successfully retrieved

### Get /v1/users/room_subs

Retrieves all room subscriptions for a single user. Requires access token in Authorization header in request.

Response body

```json
[
  {
    "id": "f0b965c5-d100-48a7-9df6-43375b41f4ed",
    "room_id": "adb0206f-dde6-4d7e-b517-36da3e899751",
    "user_id": "3053cb46-a984-471d-afa3-a14169a21b5a",
    "created_at": "2023-11-21T16:55:57.650871Z",
    "updated_at": "2023-11-21T16:55:57.650871Z"
  },
  {
    "id": "f0b965c5-d100-48a7-9df6-43375b41f4ed",
    "room_id": "adb0206f-dde6-4d7e-b517-36da3e899751",
    "user_id": "3053cb46-a984-471d-afa3-a14169a21b5a",
    "created_at": "2023-11-21T16:55:57.650871Z",
    "updated_at": "2023-11-21T16:55:57.650871Z"
  }
]
```

Status Codes

- `401` unauthorized, missing or invalid access token
- `500` failed to retrieve room subscriptions
- `200` room subscriptions retrieved successfully

## Post Vote resource

User's vote up or downvotes on a Post.

```
PostVote {
	id UUID
	post_id UUID
	user_id UUID
	created_at ISO 8601 standard datetime
	updated_at ISO 8601 standard datetime
	is_up boolean
}
```

### POST /v1/post_votes

Creates a PostVote resource. Requires access token in Authoriation header of request.

Request body

```json
{
  "is_upvote": true,
  "post_id": "111-222-fff"
}
```

Response body

```json
{
  "id": "111-222-fff",
  "post_id": "222-fff-eee",
  "user_id": "111-333-444",
  "created_at": "2023-01-01 12:15:00",
  "updated_at": "2023-01-01 12:15:00"
}
```

Status Codes

- `401` unauthorized. Missing or invalid access token.
- `400` missing field in request body
- `500` error creating resource
- `201` created resource successfully

### DELETE /v1/post_votes/{post_id}

Deletes a PostVote resource. Requires access token in Authorization header of request.

URL param `post_id` - must be a valid id for a Post in database.

Status Codes

- `401` unauthorized. Missing or invalid access token.
- `400` missing `post_id` in URL parameters
- `500` failed to delete resource
- `200` resource deleted successfully

## Authentication

Uses a access and refresh token system. Upon logging in, access and refresh tokens are sent in a response. When logging out refresh tokens are revoked and a new access token with an expired expiration time is returned.

### POST /v1/login

Logs a user in with username being email OR phone number and a password that matches a hashed password in the database.

Request body

```json
{
  "email": "example@gmail.com",
  "password": "hello1234",
  "phone": "555-111-2222"
}
```

Response body

```json
{
  "id": "111-222-333",
  "email": "example@gmail.com",
  "phone": "555-111-2222",
  "display_name": "hi mom",
  "access_token": "fjdasfjsdlsj",
  "refresh_token": "fdsajfkldsjfskjdkjfdl"
}
```

Status Codes

- `400` missing email or phone number in request body
- `404` no user with email or phone number can be found in database.
- `401` invalid login credentials
- `500` could not generate access or refresh token
- `200` user resource created

### POST /v1/logout

Logs a user out by revoking refresh token and returning an expired access token meant to replace existing access token user currently has. Requires access token in Authorization header of request.

Request body

```json
{
  "refresh_token": "jdslfjldjsfljfjds",
  "access_token": "fdjsljfsldjsfl"
}
```

Response body

```json
{
  "revoked_token": "fdjslfjslkdsjklfjdasj",
  "access_token": "fdjslfjadsjfjk"
}
```

Status Codes

- `401` unauthorized. Invalid or missing access token.
- `500` could not generate new token
- `400` tried to revoke an already revoked token
- `200` successful revocation of refresh token

### POST /v1/refresh

Refreshes a session by returning a new valid access token. Requires valid refresh token in Authorization header.

Response body

```json
{
  "access_token": "fjdslkfjslfjsdalfs"
}
```

Status Codes

- `401` invalid or missing refresh token from Authorization header.
- `500` failed to generate new token
- `200` successfully refreshed user session
