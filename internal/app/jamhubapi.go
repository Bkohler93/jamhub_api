package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	"github.com/bkohler93/jamhubapi/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	db        *database.Queries
	jwtSecret string
}

func RunApp() {
	env := getEnvVars()

	db, err := sql.Open("postgres", env.dbURL)
	if err != nil {
		log.Fatal("could not establish db connection", err)
	}

	cfg := &apiConfig{
		db:        database.New(db),
		jwtSecret: env.jwtSecret,
	}

	mux := chi.NewRouter()
	mux.Use(cors.AllowAll().Handler)

	mux.Mount("/v1", getV1Router(cfg))

	srv := http.Server{
		Addr:    "localhost:" + env.port,
		Handler: mux,
	}

	fmt.Printf("Listening on localhost:%s\n", env.port)
	log.Fatal(srv.ListenAndServe())
}
