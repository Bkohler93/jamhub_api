package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

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

	FileServer(mux)

	mux.Mount("/v1", getV1Router(cfg))

	srv := http.Server{
		Addr:    "localhost:" + env.port,
		Handler: mux,
	}

	fmt.Printf("Listening on localhost:%s\n", env.port)
	log.Fatal(srv.ListenAndServe())
}

// FileServer is serving static files.
func FileServer(router *chi.Mux) {
	root := "./"
	fs := http.FileServer(http.Dir(root))

	router.Get("/docs", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
