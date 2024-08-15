package controller

import (
	"fmt"
	"net/http"
	"os"
	"text/template"
	"time"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/auth"
	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db/users"
)

var loginTemplate = template.Must(template.ParseFiles("templates/login.html", "templates/components/navbar.html"))

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		authenticated := r.Context().Value("authenticated").(bool)
		if authenticated {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}

		if err := loginTemplate.Execute(w, nil); err != nil {
			http.Error(w, "Failed to render template", http.StatusInternalServerError)
		}
		return
	} else if r.Method == http.MethodPost {
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

		domain := os.Getenv("DOMAIN")
		if domain == "" {
			fmt.Println("DOMAIN environment variable not set, defaulting to localhost")
			domain = "localhost"
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "token",
			Value:    token,
			Path:     "/",
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
			Domain:   domain,
		})

		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {

	domain := os.Getenv("DOMAIN")
	if domain == "" {
		fmt.Println("DOMAIN environment variable not set, defaulting to localhost")
		domain = "localhost"
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "token",
		Value:    "",
		Path:     "/",
		Expires:  time.Now().Add(-1 * time.Hour),
		HttpOnly: true,
		Domain:   domain,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
