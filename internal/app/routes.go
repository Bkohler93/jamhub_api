package app

import "github.com/go-chi/chi/v5"

func getV1Router(cfg *apiConfig) *chi.Mux {
	v1Router := chi.NewRouter()

	v1Router.Post("/users", cfg.postUsersHandler)
	v1Router.Put("/users", cfg.authMiddleware(cfg.putUsersHandler))

	authRouter := chi.NewRouter()
	authRouter.Post("/login", cfg.postLoginHandler)
	authRouter.Post("/logout", cfg.authMiddleware(cfg.postLogoutHandler))
	authRouter.Post("/refresh", cfg.authRefreshMiddleware(cfg.postRefreshHandler))

	v1Router.Mount("/auth", authRouter)
	return v1Router
}
