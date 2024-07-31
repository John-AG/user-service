package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
	"user-service/internal/storage"
	"user-service/internal/user"
)

func TestAddUser(t *testing.T) {
	store = storage.New()
	handler := http.HandlerFunc(addUser)

	newUser := user.User{
		FirstName: "Alice",
		LastName:  "Bob",
		Nickname:  "AB123",
		Password:  "supersecurepassword",
		Email:     "alice@bob.com",
		Country:   "UK",
	}

	body, _ := json.Marshal(newUser)
	req, err := http.NewRequest("POST", "/users", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var returnedUser user.User
	if err := json.NewDecoder(rr.Body).Decode(&returnedUser); err != nil {
		t.Fatal(err)
	}

	if returnedUser.FirstName != newUser.FirstName || returnedUser.LastName != newUser.LastName {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedUser, newUser)
	}
}

func TestListUsers(t *testing.T) {
	store = storage.New()
	handler := http.HandlerFunc(listUsers)

	users := []user.User{
		{ID: "1", FirstName: "Alice", LastName: "Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", FirstName: "Bob", LastName: "Jones", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, u := range users {
		store.AddUser(u)
	}

	req, err := http.NewRequest("GET", "/users?page=1&pageSize=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedUsers []user.User
	if err := json.NewDecoder(rr.Body).Decode(&returnedUsers); err != nil {
		t.Fatal(err)
	}

	if len(returnedUsers) != 1 || returnedUsers[0].FirstName != "Alice" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedUsers, users[0])
	}
}

func TestUpdateUser(t *testing.T) {
	store = storage.New()
	handler := http.HandlerFunc(specificUserHandler)

	originalUser := user.User{ID: "1", FirstName: "Alice", LastName: "Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	store.AddUser(originalUser)

	updatedUser := user.User{FirstName: "AliceUpdated", LastName: "Smith", Nickname: "AS123", Password: "newpassword", Email: "alice.updated@example.com", Country: "US"}
	body, _ := json.Marshal(updatedUser)
	req, err := http.NewRequest("PUT", "/users/1", bytes.NewBuffer(body))
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedUser user.User
	if err := json.NewDecoder(rr.Body).Decode(&returnedUser); err != nil {
		t.Fatal(err)
	}

	if returnedUser.FirstName != "AliceUpdated" || returnedUser.Email != "alice.updated@example.com" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedUser, updatedUser)
	}
}

func TestDeleteUser(t *testing.T) {
	store = storage.New()
	handler := http.HandlerFunc(specificUserHandler)

	userToDelete := user.User{ID: "1", FirstName: "Alice", LastName: "Smith", CreatedAt: time.Now(), UpdatedAt: time.Now()}
	store.AddUser(userToDelete)

	req, err := http.NewRequest("DELETE", "/users/1", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}

	_, exists := store.GetUser("1")
	if exists {
		t.Errorf("user was not deleted")
	}
}

func TestListUsersWithFilter(t *testing.T) {
	store = storage.New()
	handler := http.HandlerFunc(listUsers)

	users := []user.User{
		{ID: "1", FirstName: "Alice", LastName: "Smith", Country: "UK", CreatedAt: time.Now(), UpdatedAt: time.Now()},
		{ID: "2", FirstName: "Bob", LastName: "Jones", Country: "US", CreatedAt: time.Now(), UpdatedAt: time.Now()},
	}

	for _, u := range users {
		store.AddUser(u)
	}

	req, err := http.NewRequest("GET", "/users?page=1&pageSize=10&country=UK", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var returnedUsers []user.User
	if err := json.NewDecoder(rr.Body).Decode(&returnedUsers); err != nil {
		t.Fatal(err)
	}

	if len(returnedUsers) != 1 || returnedUsers[0].Country != "UK" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedUsers, []user.User{users[0]})
	}

	req, err = http.NewRequest("GET", "/users?page=1&pageSize=10&country=US", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if err := json.NewDecoder(rr.Body).Decode(&returnedUsers); err != nil {
		t.Fatal(err)
	}

	if len(returnedUsers) != 1 || returnedUsers[0].Country != "US" {
		t.Errorf("handler returned unexpected body: got %v want %v", returnedUsers, []user.User{users[1]})
	}
}
