package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"user-service/internal/storage"
	"user-service/internal/user"

	"github.com/google/uuid"
)

var store *storage.Storage

func main() {
	store = storage.New()

	http.HandleFunc("/health", healthCheckHandler)
	http.HandleFunc("/users", usersCollectionHandler)
	http.HandleFunc("/users/", specificUserHandler)

	log.Println("Starting server on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("could not start server: %s\n", err.Error())
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func usersCollectionHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		listUsers(w, r)
	case http.MethodPost:
		addUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func specificUserHandler(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/users/")
	if id == "" || id == "/" {
		http.Error(w, "User ID required", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodPut:
		updateUser(w, r, id)
	case http.MethodDelete:
		deleteUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func addUser(w http.ResponseWriter, r *http.Request) {
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("Error decoding user: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	u.ID = generateID()
	u.CreatedAt = time.Now()
	u.UpdatedAt = time.Now()

	store.AddUser(u)
	log.Printf("User added: %v", u)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(u)
}

func listUsers(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	page, _ := strconv.Atoi(query.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(query.Get("pageSize"))
	if pageSize < 1 {
		pageSize = 10
	}

	countryFilter := query.Get("country")

	users := store.ListUsers()

	var filteredUsers []user.User
	for _, u := range users {
		if countryFilter != "" && u.Country != countryFilter {
			continue
		}
		filteredUsers = append(filteredUsers, u)
	}

	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(filteredUsers) {
		start = len(filteredUsers)
	}
	if end > len(filteredUsers) {
		end = len(filteredUsers)
	}

	paginatedUsers := filteredUsers[start:end]

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(paginatedUsers)
}

func updateUser(w http.ResponseWriter, r *http.Request, id string) {
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		log.Printf("Error decoding user: %v", err)
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	existingUser, exists := store.GetUser(id)
	if !exists {
		log.Printf("User not found: %s", id)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	u.ID = existingUser.ID
	u.CreatedAt = existingUser.CreatedAt
	u.UpdatedAt = time.Now()

	store.UpdateUser(u)
	log.Printf("User updated: %v", u)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(u)
}

func deleteUser(w http.ResponseWriter, r *http.Request, id string) {
	_, exists := store.GetUser(id)
	if !exists {
		log.Printf("User not found: %s", id)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	store.DeleteUser(id)
	log.Printf("User deleted: %s", id)
	w.WriteHeader(http.StatusNoContent)
}

func generateID() string {
	return uuid.New().String()
}
