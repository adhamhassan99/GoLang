package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type User struct {
	Name string `json:"name"`
}

var userCache = make(map[int]User)
var cacheMutex sync.RWMutex

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleRoot)

	mux.HandleFunc("POST /users", handleCreateUser)
	mux.HandleFunc("GET /users/{id}", handleGetUsers)

	fmt.Println("server listening on 8000")
	http.ListenAndServe(":8000", mux)
}
func handleRoot(w http.ResponseWriter, t *http.Request) {
	fmt.Fprintf(w, "w")
}

func handleCreateUser(w http.ResponseWriter, r *http.Request) {
	var user User

	error := json.NewDecoder(r.Body).Decode(&user)

	if error != nil {
		http.Error(
			w,
			error.Error(),
			http.StatusBadRequest,
		)

		return
	}

	if user.Name == "" {

		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	cacheMutex.Lock()
	userCache[len(userCache)+1] = user
	cacheMutex.Unlock()

	w.WriteHeader(http.StatusNoContent)
}

func handleGetUsers(w http.ResponseWriter, r *http.Request) {

	id, error := strconv.Atoi(r.PathValue("id"))

	if error != nil {
		http.Error(
			w,
			"id is required",
			http.StatusBadRequest,
		)
		return
	}

	cacheMutex.Lock()
	user, ok := userCache[id]
	cacheMutex.Unlock()

	if !ok {
		http.Error(
			w,
			"User not found",
			http.StatusNotFound,
		)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	j, err := json.Marshal(user)
	if err != nil {
		http.Error(
			w,
			"error converting json",
			http.StatusNotFound,
		)
		return
	}

	w.Write(j)
}
