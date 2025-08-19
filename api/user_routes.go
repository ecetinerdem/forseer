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

	newUserResponse, err := s.db.CreateUser(ctx, user)

	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(newUserResponse); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}

}

func (s *Server) HandleLoginUser(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	var loginUser types.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&loginUser); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()
	loginUserResponse, err := s.db.ValidateUser(ctx, loginUser)

	if err != nil {
		http.Error(w, "Unauthorized access", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&loginUserResponse); err != nil {
		http.Error(w, "Could not encode response", http.StatusInternalServerError)
		return
	}

}

func (s *Server) HandleGetUserById(w http.ResponseWriter, r *http.Request) {

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
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&foundUser); err != nil {
		http.Error(w, "Cannot encode found user", http.StatusInternalServerError)
		return
	}

}

func (s *Server) HandleGetUserByEmail(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	userId := chi.URLParam(r, "id")

	if userId == "" {
		http.Error(w, "Only users can search for other users", http.StatusBadRequest)
		return
	}

	email := r.URL.Query().Get("email") //?email=

	if email == "" {
		http.Error(w, "Only users can search for other users", http.StatusBadRequest)
		return
	}

	//if err := json.NewDecoder(r.Body).Decode(&email); err != nil {
	//	http.Error(w, "Invalid JSON request", http.StatusBadRequest)
	//	return
	//}
	//defer r.Body.Close()

	foundUser, err := s.db.GetUserByEmail(ctx, email)

	if err != nil {
		http.Error(w, "Failed to found user in db", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(&foundUser); err != nil {
		http.Error(w, "Failed to encode user", http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {
	var user types.User

	ctx := r.Context()

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid JSON request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	userId := chi.URLParam(r, "id")

	updatedUser, err := s.db.UpdateUser(ctx, userId, &user)

	if err != nil {
		http.Error(w, "Failed to update user from db", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err = json.NewEncoder(w).Encode(updatedUser); err != nil {
		http.Error(w, "Failed to update user encoding", http.StatusInternalServerError)
		return
	}
}

func (s *Server) HandleDeleteUserById(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := chi.URLParam(r, "id")

	err := s.db.DeleteUser(ctx, userId)

	if err != nil {
		http.Error(w, "Delete request cannot be fulfilled", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
