package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/auth"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db/users"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/email"
	"golang.org/x/crypto/bcrypt"
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
	type UserAccountCreationRequest struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var userReq UserAccountCreationRequest
	if err := json.NewDecoder(r.Body).Decode(&userReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	fmt.Printf("email: %s, username: %s, password: %s\n", userReq.Email, userReq.Username, userReq.Password)

	// if err := r.ParseForm(); err != nil {
	// 	http.Error(w, "Invalid request payload", http.StatusBadRequest)
	// 	return
	// }
	// userReq.Username = r.FormValue("username")
	// userReq.Email = r.FormValue("email")
	// userReq.Password = r.FormValue("password")

	// check if the username is already taken
	usernameExists, usernameErr := users.DoesUsernameExist(context.Background(), userReq.Username)
	if usernameErr != nil {
		http.Error(w, usernameErr.Error(), http.StatusInternalServerError)
		return
	}
	if usernameExists {
		http.Error(w, "Username already taken", http.StatusBadRequest)
		return
	}

	// check if the email is already taken
	emailExists, emailErr := users.DoesEmailExist(context.Background(), userReq.Email)
	if emailErr != nil {
		http.Error(w, emailErr.Error(), http.StatusInternalServerError)
		return
	}
	if emailExists {
		http.Error(w, "Email already taken", http.StatusBadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userReq.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	userReq.Password = string(hashedPassword)

	newUser, err := users.Create(context.Background(), userReq.Username, userReq.Email, userReq.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	emailVerificationToken, err := users.GenerateEmailVerificationToken(context.Background(), newUser.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	email.SendEmailVerification(newUser, emailVerificationToken)

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))

	// http.Redirect(w, r, "/create-account/success", http.StatusOK)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	type LoginResponse struct {
		Token string `json:"token"`
	}

	var loginReq LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&loginReq); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	user, err := users.GetByUsername(loginReq.Username, context.Background())
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusInternalServerError)
		return
	}

	hashedPassword, err := users.GetHashedPasswordByUsername(loginReq.Username, context.Background())
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(loginReq.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusBadRequest)
		return
	}

	token, err := auth.GenerateJWT(user.ID, user.Username)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}
