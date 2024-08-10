package controller

import (
	"net/http"
	"text/template"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	homeTemplate := template.Must(template.ParseFiles("templates/home.html", "templates/components/navbar.html"))

	authenticated := r.Context().Value("authenticated").(bool)
	username := r.Context().Value("username")

	if username == nil {
		username = ""
	}

	data := struct {
		Authenticated bool
		Username      string
		ActivePage    string
	}{
		Authenticated: authenticated,
		Username:      username.(string),
		ActivePage:    "home",
	}

	// Execute the home template, which now includes the navbar template
	if err := homeTemplate.ExecuteTemplate(w, "home.html", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// // print the home template file
// temp, err := ioutil.ReadFile("templates/components/navbar.html")
// fmt.Println("navbar: ", string(temp))

// homeTemplate, err := template.New("base").ParseFiles("templates/home.html")
// if err != nil {
// 	http.Error(w, "Failed to render template 1", http.StatusInternalServerError)
// }
