package controller

import (
	"net/http"
	"text/template"
)

var homeTemplate = template.Must(template.ParseFiles("templates/home.html"))

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	authenticated := r.Context().Value("authenticated").(bool)
	username := r.Context().Value("username")

	data := struct {
		Authenticated bool
		Username      string
	}{
		Authenticated: authenticated,
		Username:      username.(string),
	}

	if err := homeTemplate.Execute(w, data); err != nil {
		http.Error(w, "Failed to render template", http.StatusInternalServerError)
	}
}
