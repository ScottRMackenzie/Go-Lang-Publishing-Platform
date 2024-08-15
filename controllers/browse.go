package controller

import (
	"net/http"
	"text/template"
)

func BrowseHandler(w http.ResponseWriter, r *http.Request) {
	browseTemplate := template.Must(template.ParseFiles("templates/browse.html", "templates/components/navbar.html"))

	authenticated := r.Context().Value("authenticated").(bool)

	data := struct {
		Authenticated bool
		ActivePage    string
	}{
		Authenticated: authenticated,
		ActivePage:    "browse",
	}

	if err := browseTemplate.ExecuteTemplate(w, "browse.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
