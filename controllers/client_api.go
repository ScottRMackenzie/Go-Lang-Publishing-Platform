package controller

import (
	"fmt"
	"io"
	"net/http"
)

func ClientAPIHandler(w http.ResponseWriter, r *http.Request) {
	path := r.PathValue("path")
	token := r.Context().Value("token")

	fmt.Println("Token:", token)

	if token == nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	apiPort := 2337

	apiUrl := fmt.Sprintf("http://localhost:%d/%s", apiPort, "api/"+path)
	fmt.Println("API URL:", apiUrl)
	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		http.Error(w, "Failed to create request", http.StatusInternalServerError)
		return
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.(string)))

	resp, err := http.DefaultClient.Do(req)

	if resp.StatusCode == http.StatusUnauthorized {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	if err != nil {
		http.Error(w, "Failed to fetch data from API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "Failed to read response body", http.StatusInternalServerError)
		return
	}

	fmt.Println("Response:", string(body))

	w.Write(body)
}

func CheckIfAuthenticatedHandler(w http.ResponseWriter, r *http.Request) {
	authenticated := r.Context().Value("authenticated").(bool)
	if authenticated {
		w.Write([]byte("true"))
	} else {
		w.Write([]byte("false"))
	}
}
