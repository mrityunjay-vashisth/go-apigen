package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/mrityunjay-vashisth/go-apigen/pkg/generator"
)

type User struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

func main() {
	// 1) Parse the OpenAPI file
	doc, err := generator.ParseOpenAPIFile("openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// 2) Create a map of operationId -> handler function
	ops := generator.OperationMap{
		"listUsers":   listUsersHandler,   // for GET /users
		"getUserById": getUserByIdHandler, // for GET /users/{userId}
	}

	// 3) Build the Gorilla Mux router
	router, err := generator.GenerateMuxRouter(doc, ops)
	if err != nil {
		log.Fatalf("Failed to generate router: %v", err)
	}

	// 4) Start the server
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Handler for "listUsers" (GET /users?limit=<int>)
func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Extract "limit" query param if provided
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if limitStr == "" {
		limit = 10 // default from the OpenAPI
	} else if err != nil {
		http.Error(w, "Invalid 'limit' query param", http.StatusBadRequest)
		return
	}

	// In real code, you'd fetch from a DB. We'll just return some hard-coded data:
	allUsers := []User{
		{UserID: "u1", Name: "Alice"},
		{UserID: "u2", Name: "Bob"},
		{UserID: "u3", Name: "Charlie"},
	}

	// If limit < number of users, truncate the slice
	if limit < len(allUsers) {
		allUsers = allUsers[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allUsers)
}

// Handler for "getUserById" (GET /users/{userId})
func getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	userId := mux.Vars(r)["userId"] // or use generator.GetPathParam(r, "userId")

	// Fake some data in memory
	users := map[string]User{
		"u1": {UserID: "u1", Name: "Alice"},
		"u2": {UserID: "u2", Name: "Bob"},
	}

	user, ok := users[userId]
	if !ok {
		http.Error(w, "User Not Found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
