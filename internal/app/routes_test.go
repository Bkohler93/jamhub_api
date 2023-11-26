package app

import (
	"testing"

	"github.com/bkohler93/jamhubapi/internal/database"
)

func TestGetV1Router(t *testing.T) {
	mux := getV1Router(&apiConfig{
		db:        &database.Queries{},
		jwtSecret: "",
	})

	if mux == nil {
		t.Errorf("Expected mux to not be nil")
	}
}
