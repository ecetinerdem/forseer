package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ecetinerdem/forseer/types"
	"github.com/go-chi/chi/v5"
)

func (s *Server) HandleGreeting(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello Forseer")
}

func (s *Server) HandleGetUsers(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	users, err := s.db.GetUsers(ctx)

	if err != nil {
		http.Error(w, "Could not get users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&users); err != nil {
		http.Error(w, "Could not encode users", http.StatusInternalServerError)
	}

}

func (s *Server) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	var createUser types.RegisterUser

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&createUser); err != nil {
		http.Error(w, "Invalid JSON", http.StatusInternalServerError)
		return
	}

	r.Body.Close()

	user, err := types.NewUser(createUser)

	if err != nil {
		http.Error(w, "Invalid new user", http.StatusInternalServerError)
		return
	}

	newUSer, err := s.db.CreateUser(ctx, user)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newUSer); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (s *Server) handleGetUserById(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()

	userId := chi.URLParam(r, "id")

	if userId == "" {
		http.Error(w, "User Id is required", http.StatusBadRequest)
		return
	}

	foundUser, err := s.db.GetUserById(ctx, userId)

	if err != nil {
		http.Error(w, "User with given id does not exist", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusFound)

	if err := json.NewEncoder(w).Encode(&foundUser); err != nil {
		http.Error(w, "Cannot encode found user", http.StatusInternalServerError)
		return
	}

}

func (s *Server) handleUpdateUSer(w http.ResponseWriter, r *http.Request) {
	var user types.User

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userId := chi.URLParam(r, "id")

	updatedUser, err := s.db.UpdateUser(ctx, userId, &user)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err := json.NewEncoder(w).Encode(updatedUser); err != nil {
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
	}
}
