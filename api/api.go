package api

import (
	"github.com/ecetinerdem/forseer/database"
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

func (s *Server) setUpRoutes() {
	s.Router.Get("/", s.HandleGreeting)
	s.Router.Get("/users", s.HandleGetUsers)
	s.Router.Get("/users/{id}", s.handleGetUserById)
	s.Router.Get("/users/{id}/search", s.handleGetUserByEmail)
	s.Router.Put("/users/{id}", s.handleUpdateUser)
	s.Router.Post("/users", s.HandleCreateUser)
	s.Router.Delete("/users/{id}", s.handleDeleteUserById)
}
