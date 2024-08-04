package controller

import (
	"net/http"
	"text/template"
	"time"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/auth"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db/users"
)

var loginTemplate = template.Must(template.ParseFiles("templates/login.html"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		if err := loginTemplate.Execute(w, nil); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
		return
	}

	if r.Method == http.MethodPost {
		var loginReq struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		loginReq.Username = r.FormValue("username")
		loginReq.Password = r.FormValue("password")

		user, err := users.Authenticate(loginReq.Username, loginReq.Password)
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusUnauthorized)
			return
		}

		token, err := auth.GenerateJWT(user.ID, user.Username)
		if err != nil {
			http.Error(w, "Failed to generate token", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	}
}
