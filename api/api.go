package api

import (
	"github.com/ecetinerdem/forseer/db"
	"github.com/ecetinerdem/forseer/routes"
	"github.com/go-chi/chi/v5"
)

type Server struct {
	db     *db.DB
	Router *chi.Mux
}

func NewServer(database *db.DB) *Server {
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
