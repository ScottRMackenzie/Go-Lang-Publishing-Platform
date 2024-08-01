package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db/users"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/email"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/types"
)

func GetUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := users.GetAll(context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUserByIDHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if id == "" {
		errorResponse(w, http.StatusBadRequest, "id is required")
		return
	}

	user, err := users.GetByID(id, context.Background())
	if err != nil {
		errorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	var userReq types.UserAccountCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// check if the username is already taken
	usernameExists, usernameErr := users.CheckUsername(context.Background(), userReq.Username)
	if usernameErr != nil {
		http.Error(w, usernameErr.Error(), http.StatusInternalServerError)
		return
	}
	if usernameExists {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	// check if the email is already taken
	emailExists, emailErr := users.CheckEmail(context.Background(), userReq.Email)
	if emailErr != nil {
		http.Error(w, emailErr.Error(), http.StatusInternalServerError)
		return
	}
	if emailExists {
		http.Error(w, "Email already taken", http.StatusBadRequest)
		return
	}

	newUser, err := users.Create(context.Background(), userReq.Username, userReq.Email, userReq.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email.SendEmail(newUser)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
