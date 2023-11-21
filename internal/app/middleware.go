package app

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type authHandler func(w http.ResponseWriter, r *http.Request, u User)

func (cfg *apiConfig) authMiddleware(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authStr := r.Header.Get("Authorization")
		jwtEncoded := strings.Replace(authStr, "Bearer ", "", 1)

		if jwtEncoded == "" {
			respondError(w, http.StatusUnauthorized, "no access token present")
			return
		}
		claims := jwt.RegisteredClaims{}
		jwt, err := jwt.ParseWithClaims(jwtEncoded, &claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		})
		if err != nil {
			respondError(w, http.StatusUnauthorized, fmt.Sprintf("bad access token: %v", err))
			return
		}

		issuer, err := jwt.Claims.GetIssuer()
		if err != nil || issuer != "jamhub_access" {
			respondError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		id, err := jwt.Claims.GetSubject()
		if err != nil {
			respondError(w, http.StatusUnauthorized, "subject not found in token")
			return
		}

		uid, err := uuid.Parse(id)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "id could not be parsed from subject")
			return
		}

		u, err := cfg.db.GetUserByID(r.Context(), uid)
		if err != nil {
			respondError(w, http.StatusUnauthorized, fmt.Sprintf("no user with id: %s", id))
			return
		}

		user := databaseUsertoUser(u)

		handler(w, r, user)
	}
}

func (cfg *apiConfig) authRefreshMiddleware(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authStr := r.Header.Get("Authorization")
		jwtEncoded := strings.Replace(authStr, "Bearer ", "", 1)

		claims := jwt.RegisteredClaims{}
		jwt, err := jwt.ParseWithClaims(jwtEncoded, &claims, func(t *jwt.Token) (interface{}, error) {
			return []byte(cfg.jwtSecret), nil
		})
		if err != nil {
			respondError(w, http.StatusUnauthorized, "bad token")
			return
		}

		issuer, err := jwt.Claims.GetIssuer()
		if err != nil || issuer != "jamhub_refresh" {
			respondError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		id, err := jwt.Claims.GetSubject()
		if err != nil {
			respondError(w, http.StatusUnauthorized, "subject not found in token")
			return
		}

		uid, err := uuid.Parse(id)
		if err != nil {
			respondError(w, http.StatusInternalServerError, "id could not be parsed from subject")
			return
		}

		tokenUUID := uuid.MustParse(claims.ID)
		_, err = cfg.db.GetRevokedToken(r.Context(), tokenUUID)
		if err == nil {
			respondError(w, http.StatusUnauthorized, "token has been revoked")
			return
		}

		u, err := cfg.db.GetUserByID(r.Context(), uid)
		if err != nil {
			respondError(w, http.StatusUnauthorized, fmt.Sprintf("no user with id: %s", id))
			return
		}

		user := databaseUsertoUser(u)

		handler(w, r, user)
	}
}
