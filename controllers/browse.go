package controller

import (
	"net/http"
	"os"
	"text/template"
)

func BrowseHandler(w http.ResponseWriter, r *http.Request) {
	baseApiUrl := os.Getenv("BASE_URL_API")

	browseTemplate := template.Must(template.ParseFiles("templates/browse.html", "templates/components/navbar.html", "templates/components/baseURL.html"))

	authenticated := r.Context().Value("authenticated").(bool)

	data := struct {
		Authenticated bool
		ActivePage    string
		BaseApiUrl    string
	}{
		Authenticated: authenticated,
		ActivePage:    "browse",
		BaseApiUrl:    baseApiUrl,
	}

	if err := browseTemplate.ExecuteTemplate(w, "browse.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
