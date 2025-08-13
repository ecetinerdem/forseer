package api

import (
	"github.com/ecetinerdem/forseer/database"
	"github.com/ecetinerdem/forseer/routes"
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
	s.Router.Get("/", routes.HandleGreeting)
	s.Router.Get("/users", routes.HandleUsers)
}
