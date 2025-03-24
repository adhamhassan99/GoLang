package main

import (
	"fmt"
	"net/http"
	"time"
)

type Login struct {
	HashedPassword string `json:"hashedPassword"`
	SessionToken   string `json:"sessionToken"`
	CSRFToken      string
}

var users = map[string]Login{}

func handleRegister(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	if len(username) < 8 || len(password) < 8 {
		http.Error(w, "invalid username/password", http.StatusNotAcceptable)
		return
	}

	if _, ok := users[username]; ok {
		http.Error(w, "User already exists", http.StatusConflict)
		return
	}

	hashedPassword, error := hashPassword(password)

	if error != nil {
		fmt.Println("error" + error.Error())
	}

	users[username] = Login{
		HashedPassword: hashedPassword,
	}

	fmt.Println(users[username])
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	password := r.FormValue("password")

	user, ok := users[username]

	if !ok || !checkPasswordAndHash(password, user.HashedPassword) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	token := generateToken(32)
	csrfToken := generateToken(32)

	http.SetCookie(w,
		&http.Cookie{
			Name:     "session_token",
			Value:    token,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: true,
		})

	http.SetCookie(
		w,
		&http.Cookie{
			Name:     "csrf_token",
			Value:    csrfToken,
			Expires:  time.Now().Add(24 * time.Hour),
			HttpOnly: false,
		},
	)

	user.CSRFToken = csrfToken
	user.SessionToken = token
	users[username] = user

	fmt.Fprintln(w, "logged in")
}

func handleLogout(w http.ResponseWriter, r *http.Request) {

}

func protected(w http.ResponseWriter, r *http.Request) {

}

func main() {
	Mux := http.NewServeMux()

	Mux.HandleFunc("/register", handleRegister)
	Mux.HandleFunc("POST /login", handleLogin)
	Mux.HandleFunc("/logout", handleLogout)
	Mux.HandleFunc("/protected", protected)

	http.ListenAndServe(":8080", Mux)
}
