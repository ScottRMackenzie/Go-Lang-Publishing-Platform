package controller

import (
	"net/http"
	"os"
	"text/template"
)

func BrowseHandler(w http.ResponseWriter, r *http.Request) {
	DOMAIN := os.Getenv("DOMAIN")

	browseTemplate := template.Must(template.ParseFiles("templates/browse.html", "templates/components/navbar.html", "templates/components/baseURL.html"))

	authenticated := r.Context().Value("authenticated").(bool)

	data := struct {
		Authenticated bool
		ActivePage    string
		DOMAIN        string
	}{
		Authenticated: authenticated,
		ActivePage:    "browse",
		DOMAIN:        DOMAIN,
	}

	if err := browseTemplate.ExecuteTemplate(w, "browse.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
