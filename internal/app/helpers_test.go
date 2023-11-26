package app

import (
	"log"
	"testing"

	"github.com/joho/godotenv"
)

func TestGetEnv(t *testing.T) {
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("should have been able to load Env")
	}
	envVars := getEnvVars()

	if envVars.dbURL == "" {
		t.Errorf("Expected dbURL to be populated")
	}

	if envVars.jwtSecret == "" {
		t.Error("Expected Jwt Secret to be populated")
	}

	if envVars.port == "" {
		t.Error("Expected port to be populated")
	}
}
