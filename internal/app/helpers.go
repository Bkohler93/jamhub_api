package app

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.WriteHeader(status)

	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write(data)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	res := map[string]string{
		"error": msg,
	}
	respondJSON(w, status, res)
}

func generateToken(userID string, secret string, expiresAt time.Time, issuer string) (token string, err error) {

	claims := jwt.RegisteredClaims{
		Issuer:  issuer,
		Subject: userID,
		ExpiresAt: &jwt.NumericDate{
			Time: expiresAt,
		},
		IssuedAt: &jwt.NumericDate{
			Time: time.Now(),
		},
		ID: uuid.New().String(),
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err = jwt.SignedString([]byte(secret))
	if err != nil {
		return token, err
	}

	return token, err
}

type envVariables struct {
	port      string
	dbURL     string
	jwtSecret string
}

func getEnvVars() envVariables {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	jwtSecret := os.Getenv("SUPER_SECRET")
	if jwtSecret == "" {
		log.Fatal("require SUPER_SECRET in .env")
	}

	dbURL := os.Getenv("DATABASE_APP_URL")
	if dbURL == "" {
		log.Fatal("require DATABASE_APP_URL in .env")
	}

	return envVariables{
		port:      port,
		dbURL:     os.Getenv("DATABASE_APP_URL"),
		jwtSecret: os.Getenv("SUPER_SECRET"),
	}
}
