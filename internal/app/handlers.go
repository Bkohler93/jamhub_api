package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bkohler93/jamhubapi/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (cfg *apiConfig) postUsersHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := struct {
		Email       string `json:"email"`
		Password    string `json:"password"`
		DisplayName string `json:"display_name"`
		Phone       string `json:"phone"`
	}{}

	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&reqBody)

	if reqBody.Email == "" && reqBody.Phone == "" {
		respondError(w, http.StatusBadRequest, "user requires email or phone number to register")
		return
	}

	_, err := cfg.db.GetUserByDisplayName(r.Context(), reqBody.DisplayName)
	if err == nil {
		respondError(w, http.StatusInternalServerError, "user already exists with that display name")
		return
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not hash password")
		return
	}

	isValidEmail, isValidPhone := true, true
	if reqBody.Email == "" {
		isValidEmail = false
	}

	if reqBody.Phone == "" {
		isValidPhone = false
	}

	user, err := cfg.db.CreateUser(r.Context(), database.CreateUserParams{
		ID:           uuid.New(),
		Email:        sql.NullString{String: reqBody.Email, Valid: isValidEmail},
		Phone:        sql.NullString{String: reqBody.Phone, Valid: isValidPhone},
		PasswordHash: string(pwHash),
		DisplayName:  reqBody.DisplayName,
		UpdatedAt:    time.Now(),
		CreatedAt:    time.Now(),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "failed to create user")
	}

	resBody := struct {
		ID          uuid.UUID `json:"id"`
		Email       string    `json:"email"`
		Phone       string    `json:"phone"`
		DisplayName string    `json:"display_name"`
		UpdatedAt   time.Time `json:"updated_at"`
		CreatedAt   time.Time `json:"created_at"`
	}{
		ID:          user.ID,
		Email:       user.Email.String,
		Phone:       user.Phone.String,
		DisplayName: user.DisplayName,
		UpdatedAt:   user.UpdatedAt,
		CreatedAt:   user.CreatedAt,
	}

	respondJSON(w, http.StatusCreated, resBody)
}

func (cfg *apiConfig) putUsersHandler(w http.ResponseWriter, r *http.Request, u User) {
	reqBody := struct {
		Email       string `json:"email"`
		Phone       string `json:"phone"`
		DisplayName string `json:"display_name"`
		Password    string `json:"password"`
	}{}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.Decode(&reqBody)

	emailNotNull := true
	if reqBody.Email == "" {
		emailNotNull = false
	}
	phoneNotNull := true
	if reqBody.Phone == "" {
		phoneNotNull = false
	}

	var passwordHash []byte
	passwordNotNull := false
	var err error
	if reqBody.Password != "" {
		passwordNotNull = true
		passwordHash, err = bcrypt.GenerateFromPassword([]byte(reqBody.Password), bcrypt.DefaultCost)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "failed to encrypt password")
			return
		}
	}

	displayNameNotNull := true
	if reqBody.DisplayName == "" {
		displayNameNotNull = false
	}

	user, err := cfg.db.UpdateUser(r.Context(), database.UpdateUserParams{
		Email:        sql.NullString{String: reqBody.Email, Valid: emailNotNull},
		Phone:        sql.NullString{String: reqBody.Phone, Valid: phoneNotNull},
		UpdatedAt:    time.Now(),
		ID:           u.ID,
		PasswordHash: sql.NullString{String: string(passwordHash), Valid: passwordNotNull},
		DisplayName:  sql.NullString{String: reqBody.DisplayName, Valid: displayNameNotNull},
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to update user in database: %v", err))
		return
	}

	u = databaseUsertoUser(user)
	respondJSON(w, http.StatusOK, u)
}

func (cfg *apiConfig) postLoginHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Phone    string `json:"phone"`
	}{}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&reqBody)

	if reqBody.Email == "" && reqBody.Phone == "" {
		respondError(w, http.StatusBadRequest, "require email or phone to login")
		return
	}
	var user database.User
	var err error
	if reqBody.Email != "" {
		user, err = cfg.db.GetUserByEmail(r.Context(), sql.NullString{String: reqBody.Email, Valid: true})
		if err != nil {
			respondError(w, http.StatusNotFound, "no user with that email exists")
			return
		}
	} else {
		user, err = cfg.db.GetUserByPhone(r.Context(), sql.NullString{String: reqBody.Phone, Valid: true})
		if err != nil {
			respondError(w, http.StatusNotFound, "no user with that phone number exists")
			return
		}
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(reqBody.Password))
	if err != nil {
		respondError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	accessToken, err := generateToken(user.ID.String(), cfg.jwtSecret, time.Now().Add(time.Minute*15), "jamhub_access")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error generating access token")
	}
	refreshToken, err := generateToken(user.ID.String(), cfg.jwtSecret, time.Now().Add(time.Hour*99999), "jamhub_refresh")
	if err != nil {
		respondError(w, http.StatusInternalServerError, "error generating refresh token")
	}

	resBody := struct {
		ID           uuid.UUID `json:"id"`
		Email        string    `json:"email"`
		Phone        string    `json:"phone"`
		DisplayName  string    `json:"display_name"`
		AccessToken  string    `json:"access_token"`
		RefreshToken string    `json:"refresh_token"`
	}{
		ID:           user.ID,
		Email:        user.Email.String,
		Phone:        user.Phone.String,
		DisplayName:  user.DisplayName,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	respondJSON(w, http.StatusOK, resBody)
}

func (cfg *apiConfig) postLogoutHandler(w http.ResponseWriter, r *http.Request, u User) {
	reqBody := struct {
		RefreshToken string `json:"refresh_token"`
		AccessToken  string `json:"access_token"`
	}{}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&reqBody)

	if reqBody.RefreshToken == "" {
		respondError(w, http.StatusBadRequest, "refresh token required to logout")
		return
	}

	//parse refresh token to retrieve ID
	claims := jwt.RegisteredClaims{}
	refreshJwt, err := jwt.ParseWithClaims(reqBody.RefreshToken, &claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondError(w, http.StatusUnauthorized, "bad token")
		return
	}

	issuer, err := refreshJwt.Claims.GetIssuer()
	if err != nil || issuer != "jamhub_refresh" {
		respondError(w, http.StatusUnauthorized, "bad refresh token")
		return
	}

	tokenUID := claims.ID
	tokenUUID := uuid.MustParse(tokenUID)

	//generate access token to change expiration date to now
	mClaims := jwt.MapClaims{}
	accessJwt, err := jwt.ParseWithClaims(reqBody.AccessToken, &mClaims, func(t *jwt.Token) (interface{}, error) {
		return []byte(cfg.jwtSecret), nil
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not parse access token")
		return
	}

	mClaims["exp"] = time.Now().UTC().Unix()

	newAccessToken := jwt.NewWithClaims(accessJwt.Method, mClaims)
	newAccessTokenString, err := newAccessToken.SignedString([]byte(cfg.jwtSecret))
	if err != nil {
		respondError(w, http.StatusInternalServerError, "could not create new access token")
		return
	}
	tok, err := cfg.db.CreateRevokedToken(r.Context(), database.CreateRevokedTokenParams{
		ID:        tokenUUID,
		UserID:    u.ID,
		RevokedAt: time.Now(),
	})
	if err != nil {
		log.Printf("tried to revoke an already revoked token: %v", tokenUUID)
		respondError(w, http.StatusBadRequest, "user is already logged out")
		return
	}

	resBody := struct {
		Token       uuid.UUID `json:"revoked_token"`
		AccessToken string    `json:"access_token"`
	}{
		Token:       tok.ID,
		AccessToken: newAccessTokenString,
	}

	respondJSON(w, http.StatusOK, resBody)
}

func (cfg *apiConfig) postRefreshHandler(w http.ResponseWriter, r *http.Request, u User) {
	accessToken, err := generateToken(u.ID.String(), cfg.jwtSecret, time.Now().Add(time.Minute*15), "jamhub_access")
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to generate new token: %v", err))
		return
	}

	resBody := struct {
		AccessToken string `json:"access_token"`
	}{
		AccessToken: accessToken,
	}

	respondJSON(w, http.StatusOK, resBody)
}

func (cfg *apiConfig) postRoomsHandler(w http.ResponseWriter, r *http.Request, u User) {
	reqBody := struct {
		Name string `json:"name"`
	}{}
	defer r.Body.Close()

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&reqBody)

	rm, err := cfg.db.CreateRoom(r.Context(), database.CreateRoomParams{
		ID:        uuid.New(),
		Name:      reqBody.Name,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("error creating room, %v", err))
		return
	}
	room := databaseRoomToRoom(rm)

	respondJSON(w, http.StatusCreated, room)
}

func (cfg *apiConfig) getRoomsHandler(w http.ResponseWriter, r *http.Request) {
	var limit int
	l := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(l)
	if err != nil {
		limit = 10
	}
	rooms, err := cfg.db.GetRooms(r.Context(), int32(limit))
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("could not get rooms - %v", err))
		return
	}

	respondJSON(w, http.StatusOK, rooms)
}

func (cfg *apiConfig) getRoomByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "room_id")
	uid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid room id")
		return
	}

	rm, err := cfg.db.GetRoomByID(r.Context(), uid)
	if err != nil {
		respondError(w, http.StatusNotFound, "failed to find room with that id")
		return
	}

	respondJSON(w, http.StatusOK, rm)
}

func (cfg *apiConfig) deleteRoomByIDHandler(w http.ResponseWriter, r *http.Request, u User) {
	id := chi.URLParam(r, "room_id")
	uid, err := uuid.Parse(id)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid room id")
	}

	err = cfg.db.DeleteRoom(r.Context(), uid)
	if err != nil {
		respondError(w, http.StatusOK, "no room to delete with that id")
	}

	respondJSON(w, http.StatusOK, nil)
}

func (cfg *apiConfig) postPostsHandler(w http.ResponseWriter, r *http.Request, u User) {
	reqBody := struct {
		RoomID string `json:"room_id"`
		Link   string `json:"link"`
	}{}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()

	decoder.Decode(&reqBody)
	if reqBody.RoomID == "" || reqBody.Link == "" {
		respondError(w, http.StatusBadRequest, "Expected room_id and link in request body")
		return
	}

	roomUID, err := uuid.Parse(reqBody.RoomID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid room_id")
		return
	}

	p, err := cfg.db.CreatePost(r.Context(), database.CreatePostParams{
		ID:        uuid.New(),
		UserID:    u.ID,
		RoomID:    roomUID,
		Link:      reqBody.Link,
		UpdatedAt: time.Now(),
		CreatedAt: time.Now(),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("could not create post - %v", err))
		return
	}

	post := databasePostToPost(p)

	respondJSON(w, http.StatusCreated, post)
}

func (cfg *apiConfig) deletePostsHandler(w http.ResponseWriter, r *http.Request, u User) {
	postID := chi.URLParam(r, "post_id")
	if postID == "" {
		respondError(w, http.StatusBadRequest, "requires post id")
		return
	}
	postUID, err := uuid.Parse(postID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid post id")
		return
	}

	p, err := cfg.db.GetPost(r.Context(), postUID)
	if err != nil {
		respondError(w, http.StatusNotFound, "no post with that id was found")
		return
	}

	if p.UserID != u.ID {
		respondError(w, http.StatusUnauthorized, "user not authorized to delete this post")
		return
	}

	err = cfg.db.DeletePost(r.Context(), p.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("error deleting post - %v", err))
	}

	respondJSON(w, http.StatusOK, nil)
}

func (cfg *apiConfig) getRoomPostsHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := struct {
		RoomID string `json:"room_id"`
	}{}
	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&reqBody)

	if reqBody.RoomID == "" {
		respondError(w, http.StatusBadRequest, "requires room_id to retrieve posts for room")
		return
	}
	roomUID, err := uuid.Parse(reqBody.RoomID)
	if err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("invalid room_id - %v", err))
		return
	}

	ps, err := cfg.db.GetRoomPosts(r.Context(), roomUID)
	if err != nil {
		respondError(w, http.StatusNotFound, "no posts found in room with room_id")
		return
	}

	var posts []Post
	for _, p := range ps {
		posts = append(posts, databasePostToPost(p))
	}

	respondJSON(w, http.StatusOK, posts)
}

func (cfg *apiConfig) postRoomSubsHandler(w http.ResponseWriter, r *http.Request, u User) {
	reqBody := struct {
		RoomID string `json:"room_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	decoder.Decode(&reqBody)

	if reqBody.RoomID == "" {
		respondError(w, http.StatusBadRequest, "expected room_id in request - ")
		return
	}

	roomUID, err := uuid.Parse(reqBody.RoomID)
	if err != nil {
		respondError(w, http.StatusBadGateway, fmt.Sprintf("invalid room_id - %v", err))
		return
	}

	_, err = cfg.db.GetRoomByID(r.Context(), roomUID)
	if err != nil {
		respondError(w, http.StatusNotFound, fmt.Sprintf("no room with that id exists - %v", err))
		return
	}

	rms, err := cfg.db.CreateRoomSubscription(r.Context(), database.CreateRoomSubscriptionParams{
		ID:        uuid.New(),
		RoomID:    roomUID,
		UserID:    u.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("server could not create room subscription - %v", err))
		return
	}

	roomSub := databaseRoomSubscriptionToRoomSubscription(rms)

	respondJSON(w, http.StatusCreated, roomSub)
}

func (cfg *apiConfig) deleteRoomSubsHandler(w http.ResponseWriter, r *http.Request, u User) {
	roomSubID := chi.URLParam(r, "room_sub_id")
	roomSubUUID, err := uuid.Parse(roomSubID)
	if err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("expected valid id query parameter - %v", err))
		return
	}

	if err = cfg.db.DeleteRoomSubscription(r.Context(), database.DeleteRoomSubscriptionParams{
		RoomID: roomSubUUID,
		UserID: u.ID,
	}); err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete room susbcription - %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) getAllRoomSubsHandler(w http.ResponseWriter, r *http.Request) {
	rms, err := cfg.db.GetAllRoomSubs(r.Context())
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to retrieve room subs - %v", err))
		return
	}

	var roomSubs []RoomSubscription
	for _, rm := range rms {
		roomSubs = append(roomSubs, databaseRoomSubscriptionToRoomSubscription(rm))
	}

	respondJSON(w, http.StatusOK, roomSubs)
}

func (cfg *apiConfig) getAllRoomRoomSubsHandler(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "room_id")
	roomUID, err := uuid.Parse(roomID)
	if err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("expected valid room id in url - %v", err))
		return
	}
	rms, err := cfg.db.GetRoomRoomSubscriptions(r.Context(), roomUID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("error retrieving room subs for given room - %v", err))
		return
	}

	var roomSubs []RoomSubscription
	for _, rm := range rms {
		roomSubs = append(roomSubs, databaseRoomSubscriptionToRoomSubscription(rm))
	}

	respondJSON(w, http.StatusOK, roomSubs)
}

func (cfg *apiConfig) getUserRoomSubsHandler(w http.ResponseWriter, r *http.Request, u User) {
	rms, err := cfg.db.GetUserRoomSusbcriptions(r.Context(), u.ID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to retrieve room subs for user - %v", err))
		return
	}

	var roomSubs []RoomSubscription
	for _, rm := range rms {
		roomSubs = append(roomSubs, databaseRoomSubscriptionToRoomSubscription(rm))
	}

	respondJSON(w, http.StatusOK, roomSubs)
}

func (cfg *apiConfig) postPostVotesHandler(w http.ResponseWriter, r *http.Request, u User) {
	reqBody := struct {
		IsUpvote bool   `json:"is_upvote"`
		PostID   string `json:"post_id"`
	}{}

	defer r.Body.Close()
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&reqBody)
	if err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("requires 'isUpvote' field in request body - %v", err))
		return
	}

	postUID, err := uuid.Parse(reqBody.PostID)
	if err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("invalid post_id - %v", err))
		return
	}

	pv, err := cfg.db.CreatePostVote(r.Context(), database.CreatePostVoteParams{
		ID:        uuid.New(),
		PostID:    postUID,
		UserID:    u.ID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		IsUp:      reqBody.IsUpvote,
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("error creating post vote - %v", err))
		return
	}
	postVote := databasePostVoteToPostVote(pv)

	respondJSON(w, http.StatusCreated, postVote)
}

func (cfg *apiConfig) deletePostHandler(w http.ResponseWriter, r *http.Request, u User) {
	postID := chi.URLParam(r, "post_id")
	if postID == "" {
		respondError(w, http.StatusBadRequest, "requires 'post_id' in URL")
		return
	}
	postUID, err := uuid.Parse(postID)
	if err != nil {
		respondError(w, http.StatusBadRequest, fmt.Sprintf("expected valid post_id - %v", err))
		return
	}

	err = cfg.db.DeletePostVote(r.Context(), postUID)
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("error deleting post vote - %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (cfg *apiConfig) getUserSubscribedRoomsHandler(w http.ResponseWriter, r *http.Request, u User) {
	var l int
	var o int

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	l, err := strconv.Atoi(limit)
	if err != nil {
		l = 10
	}

	o, err = strconv.Atoi(offset)
	if err != nil {
		o = 0
	}

	rms, err := cfg.db.GetUserRoomsOrderedBySubs(r.Context(), database.GetUserRoomsOrderedBySubsParams{
		UserID: u.ID,
		Limit:  int32(l),
		Offset: int32(o),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("error retrieving user subscribed rooms - %v", err))
		return
	}

	respondJSON(w, http.StatusOK, rms)
}

func (cfg *apiConfig) getRoomsOrderedByRoomSubsHandler(w http.ResponseWriter, r *http.Request) {
	var l int
	var o int

	limit := r.URL.Query().Get("limit")
	offset := r.URL.Query().Get("offset")

	l, err := strconv.Atoi(limit)
	if err != nil {
		l = 10
	}

	o, err = strconv.Atoi(offset)
	if err != nil {
		o = 0
	}
	rms, err := cfg.db.GetRoomsOrderedBySubs(r.Context(), database.GetRoomsOrderedBySubsParams{
		Limit:  int32(l),
		Offset: int32(o),
	})
	if err != nil {
		respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to retrieve rooms - %v", err))
		return
	}

	respondJSON(w, http.StatusOK, rms)
}

func (cfg *apiConfig) getRoomPostsOrderedHandler(w http.ResponseWriter, r *http.Request) {
	reqBody := struct {
		RoomID string `json:"room_id"`
	}{}

	decoder := json.NewDecoder(r.Body)
	decoder.Decode(&reqBody)
	defer r.Body.Close()

	if reqBody.RoomID == "" {
		respondError(w, http.StatusBadRequest, "failed to provide room_id with request")
		return
	}

	roomUID, err := uuid.Parse(reqBody.RoomID)
	if err != nil {
		respondError(w, http.StatusBadRequest, "invalid room_id provided")
		return
	}

	orderedOption := r.URL.Query().Get("select")
	if orderedOption == "" {
		orderedOption = "new"
	}

	if orderedOption == "new" {
		rms, err := cfg.db.GetNewRoomPosts(r.Context(), roomUID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to retrieve new room posts - %v", err))
			return
		}

		respondJSON(w, http.StatusOK, rms)
	} else {
		rms, err := cfg.db.GetTopRoomPosts(r.Context(), roomUID)
		if err != nil {
			respondError(w, http.StatusInternalServerError, fmt.Sprintf("failed to retrieve top room posts - %v", err))
			return
		}

		respondJSON(w, http.StatusOK, rms)
	}
}
