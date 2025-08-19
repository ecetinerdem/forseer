package api

import (
	"github.com/ecetinerdem/forseer/database"
	"github.com/ecetinerdem/forseer/middleware"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	db     *database.DB
	Router *chi.Mux
}

func NewServer(database *database.DB) *Server {
	s := &Server{
		db:     database,
		Router: chi.NewRouter(),
	}
	s.setUpRoutes()
	return s
}

func (s *Server) setUpRoutes() *chi.Mux {
	s.Router.Get("/", s.HandleGreeting)
	s.Router.Post("/register", s.HandleCreateUser)
	s.Router.Post("/login", s.HandleLoginUser)

	s.Router.Route("/api/v1", func(r chi.Router) {

		r.Route("/users", func(userRouter chi.Router) {
			userRouter.Use(middleware.UserAuthentication)
			userRouter.Get("/users", s.HandleGetUsers)
			userRouter.Get("/users/{id}", s.HandleGetUserById)
			userRouter.Put("/users/{id}", s.HandleUpdateUser)
			userRouter.Delete("/users/{id}", s.HandleDeleteUserById)
		})

		r.Get("/search/users", s.HandleGetUserByEmail)

	})

	return s.Router
}
