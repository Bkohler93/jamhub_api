package app

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bkohler93/jamhubapi/internal/database"
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
