package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/ecetinerdem/forseer/db"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

type Server struct {
	db     *db.DB
	router *chi.Mux
}

func NewServer(database *db.DB) *Server {
	s := &Server{
		db:     database,
		router: chi.NewRouter(),
	}
	s.setUpRoutes()
	return s
}

func (s *Server) setUpRoutes() {
	s.router.Get("/", s.handleGreeting)
	s.router.Get("/users", s.handleUsers)
}

func (s *Server) handleGreeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello")
}

func (s *Server) handleUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hanle users")
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	database, err := db.NewDB()

	if err != nil {
		log.Fatal("Database connection error: ", err)
	}

	defer database.Close()

	server := NewServer(database)
	PORT := os.Getenv("PORT")
	log.Println("Server starting on :8080")
	log.Fatal(http.ListenAndServe(":"+PORT, server.router))

}
