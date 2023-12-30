package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bkohler93/jamhubapi/internal/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/lib/pq"
)

// type apiConfig struct {
// 	db        *database.Queries
// 	jwtSecret string
// }

type apiConfig struct {
	db        DB
	jwtSecret string
}

func NewConfig(db DB, jwtSecret string) *apiConfig {
	return &apiConfig{
		db:        db,
		jwtSecret: jwtSecret,
	}
}

func RunApp() {
	env := getEnvVars()

	db, err := sql.Open("postgres", env.dbURL)
	if err != nil {
		log.Fatal("could not establish db connection", err)
	}

	cfg := NewConfig(database.New(db), env.jwtSecret)

	// cfg := &apiConfig{
	// 	db:        database.New(db),
	// 	jwtSecret: env.jwtSecret,
	// }

	mux := chi.NewRouter()
	mux.Use(cors.AllowAll().Handler)
	mux.Get("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("Jamhub API is ready for requests."))
	})
	mux.Get("/db-connection/ready", func(w http.ResponseWriter, r *http.Request) {
		err = db.Ping()
		if err != nil {
			w.WriteHeader(500)
			w.Write([]byte(fmt.Sprintf("db connection failed: %s", err)))
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("db connected and ready"))
	})
	mux.Mount("/v1", getV1Router(cfg))

	srv := http.Server{
		Addr:              ":" + env.port,
		Handler:           mux,
		ReadHeaderTimeout: time.Second * 10,
	}

	fmt.Printf("Listening on port:%s\n", env.port)
	log.Fatal(srv.ListenAndServe())
}
