package bj_controller

import (
	"net/http"
	"os"
	"text/template"

	"github.com/ScottRMackenzie/Go-Lang-Publishing-Platform/db/users"
)

func BlackjackFrontendHandler(w http.ResponseWriter, r *http.Request) {
	DOMAIN := os.Getenv("DOMAIN")

	browseTemplate := template.Must(template.ParseFiles("templates/games/blackjack/bj.html", "templates/components/navbar.html", "templates/components/baseURL.html"))

	authenticated := r.Context().Value("authenticated").(bool)

	balance := 0
	if authenticated {
		user, err := users.GetByUsername(r.Context().Value("username").(string), r.Context())
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		balance = user.Balance
	}

	data := struct {
		Authenticated bool
		ActivePage    string
		DOMAIN        string
		Balance       int
	}{
		Authenticated: authenticated,
		ActivePage:    "Games",
		DOMAIN:        DOMAIN,
		Balance:       balance,
	}

	if err := browseTemplate.ExecuteTemplate(w, "bj.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
