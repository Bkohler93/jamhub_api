package app

import (
	"github.com/go-chi/chi/v5"
)

func getV1Router(cfg *apiConfig) *chi.Mux {
	v1Router := chi.NewRouter()

	v1Router.Post("/users", cfg.postUsersHandler)
	v1Router.Put("/users", cfg.authMiddleware(cfg.putUsersHandler))
	v1Router.Get("/users/rooms/room_subscriptions", cfg.authMiddleware(cfg.getUserSubscribedRoomsHandler))

	v1Router.Post("/rooms", cfg.authMiddleware(cfg.postRoomsHandler))
	v1Router.Get("/rooms", cfg.getRoomsHandler)
	v1Router.Get("/rooms/{room_id}", cfg.getRoomByIDHandler)
	v1Router.Delete("/rooms/{room_id}", cfg.authMiddleware(cfg.deleteRoomByIDHandler))
	v1Router.Get("/rooms/room_subscriptions", cfg.getRoomsOrderedByRoomSubsHandler)

	v1Router.Post("/posts", cfg.authMiddleware(cfg.postPostsHandler))
	v1Router.Delete("/posts/{post_id}", cfg.authMiddleware(cfg.deletePostsHandler))
	v1Router.Get("/rooms/posts/{room_id}", cfg.getRoomPostsHandler)
	v1Router.Get("/posts/rooms", cfg.getRoomPostsOrderedHandler)

	v1Router.Post("/room_subs", cfg.authMiddleware(cfg.postRoomSubsHandler))
	v1Router.Delete("/room_subs/{room_id}", cfg.authMiddleware(cfg.deleteRoomSubsHandler))
	v1Router.Get("/room_subs", cfg.getAllRoomSubsHandler)
	v1Router.Get("/rooms/room_subs/{room_id}", cfg.getAllRoomRoomSubsHandler)
	v1Router.Get("/users/room_subs", cfg.authMiddleware(cfg.getUserRoomSubsHandler))

	v1Router.Post("/post_votes", cfg.authMiddleware(cfg.postPostVotesHandler))
	v1Router.Delete("/post_votes/{post_id}", cfg.authMiddleware(cfg.deletePostHandler))

	authRouter := chi.NewRouter()
	authRouter.Post("/login", cfg.postLoginHandler)
	authRouter.Post("/logout", cfg.authMiddleware(cfg.postLogoutHandler))
	authRouter.Post("/refresh", cfg.authRefreshMiddleware(cfg.postRefreshHandler))

	v1Router.Mount("/auth", authRouter)
	return v1Router
}
