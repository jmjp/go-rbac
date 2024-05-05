package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
)

// ParseBody parses the body of an HTTP request and unmarshals it into the provided value.
//
// It takes a pointer to an http.Request and a value of any type as parameters.
// It returns an error.
func getBody(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func sendJson(w http.ResponseWriter, code int, body interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

func sendString(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	body := map[string]string{"message": message}
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(body)
}

func sendStatus(w http.ResponseWriter, code int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(code)
}

func configureCookie(name, value string, expires time.Time) *http.Cookie {
	return &http.Cookie{
		Name:     name,
		Value:    value,
		Expires:  expires,
		HttpOnly: true,
		Path:     "/",
		Domain:   os.Getenv("host"),
		Secure:   os.Getenv("env") == "prod",
	}
}
