package app

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {

	data, err := json.Marshal(payload)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(status)
	w.Write(data) // #nosec G104
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
	var (
		port      string
		ok        bool
		jwtSecret string
		dbURL     string
	)

	err := godotenv.Load()
	if err != nil {
		log.Println("Using manually set .env", err) // production
	}

	if port, ok = os.LookupEnv("PORT"); !ok {
		port = "8080"
		log.Println("PORT var not found")
	}

	if jwtSecret, ok = os.LookupEnv("SUPER_SECRET"); !ok {
		log.Fatal("require SUPER_SECRET in .env")
	}

	if dbURL, ok = os.LookupEnv("DATABASE_URL"); !ok {
		log.Fatal("require DATABASE_URL in .env")
	}
	dbMode := "?sslmode=disable"

	fmt.Println(dbURL)

	return envVariables{
		port:      port,
		dbURL:     fmt.Sprintf("%s%s", dbURL, dbMode),
		jwtSecret: jwtSecret,
	}
}
