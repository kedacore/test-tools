package main

import (
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
)

const (
	PORT     = 8080
	TLS_PORT = 4333
)

var value int

type response struct {
	Value   int  `json:"value"`
	Success bool `json:"success"`
}

type application struct {
	auth struct {
		basic struct {
			username string
			password string
		}
		bearer struct {
			token string
		}
	}
}

func (app *application) basicAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			usernameHash := sha256.Sum256([]byte(username))
			passwordHash := sha256.Sum256([]byte(password))
			expectedUsernameHash := sha256.Sum256([]byte(app.auth.basic.username))
			expectedPasswordHash := sha256.Sum256([]byte(app.auth.basic.password))

			usernameMatch := (subtle.ConstantTimeCompare(usernameHash[:], expectedUsernameHash[:]) == 1)
			passwordMatch := (subtle.ConstantTimeCompare(passwordHash[:], expectedPasswordHash[:]) == 1)

			if usernameMatch && passwordMatch {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func (app *application) bearerTokenAuth(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		header := r.Header.Get("Authorization")
		if header == fmt.Sprintf("Bearer %s", app.auth.bearer.token) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func getValue(w http.ResponseWriter, r *http.Request) {

	rsp := response{
		Value:   value,
		Success: true,
	}

	json.NewEncoder(w).Encode(rsp)
}

func setValue(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	number := vars["number"]

	value, _ = strconv.Atoi(number)
	json.NewEncoder(w).Encode(value)
}

func main() {
	fmt.Printf("Running server on port: %d", PORT)

	app := new(application)
	app.auth.basic.username = os.Getenv("AUTH_USERNAME")
	app.auth.basic.password = os.Getenv("AUTH_PASSWORD")
	app.auth.bearer.token = os.Getenv("AUTH_TOKEN")

	r := mux.NewRouter()

	r.HandleFunc("/api/value", getValue).Methods("GET")
	r.HandleFunc("/api/basic/value", app.basicAuth(getValue)).Methods("GET")
	r.HandleFunc("/api/token/value", app.bearerTokenAuth(getValue)).Methods("GET")

	r.HandleFunc("/api/value/{number:[0-9]+}", setValue).Methods("POST")

	http.Handle("/", r)
	http.ListenAndServe(fmt.Sprintf(":%d", PORT), nil)

	value, found := os.LookupEnv("USE_TLS")
	if found && value == "true" {
		http.ListenAndServeTLS(fmt.Sprintf(":%d", TLS_PORT), "/cert/tls.crt", "/cert/tls.key", nil)
	}
}
