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

	s.Router.Route("api/v1", func(r chi.Router) {
		r.Use(middleware.UserAuthentication)
		r.Get("/users", s.HandleGetUsers)
		r.Get("/users/{id}", s.handleGetUserById)
		r.Get("/users/{id}/search", s.handleGetUserByEmail)
		r.Put("/users/{id}", s.handleUpdateUser)
		r.Post("/users", s.HandleCreateUser)
		r.Delete("/users/{id}", s.handleDeleteUserById)
	})

	return s.Router
}
