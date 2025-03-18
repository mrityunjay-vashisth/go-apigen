package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/mrityunjay-vashisth/go-apigen/pkg/generator"
)

type User struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

func main() {
	// 1) Parse the OpenAPI spec
	doc, err := generator.ParseOpenAPIFile("openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to parse OpenAPI: %v", err)
	}

	// 2) Define a map of operationId -> RouteDefinition
	ops := generator.OperationMap{
		"listUsers": generator.RouteDefinition{
			Handler:     listUsersHandler,
			Middlewares: []mux.MiddlewareFunc{logOperationMiddleware("listUsers")},
		},
		"getUserById": generator.RouteDefinition{
			Handler:     getUserByIdHandler,
			Middlewares: []mux.MiddlewareFunc{logOperationMiddleware("getUserById")},
		},
	}

	// 3) Create global middlewares (applies to all routes)
	// e.g. a "requestLoggerMiddleware" that logs method & path
	global := []mux.MiddlewareFunc{
		requestLoggerMiddleware,
	}

	// 4) Generate the router
	router, err := generator.GenerateMuxRouter(doc, ops, global...)
	if err != nil {
		log.Fatalf("Failed to generate router: %v", err)
	}

	// 5) Start the HTTP server
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// =======================
// SAMPLE HANDLERS
// =======================

func listUsersHandler(w http.ResponseWriter, r *http.Request) {
	// Extract "limit" query param
	limitStr := r.URL.Query().Get("limit")
	limit, err := strconv.Atoi(limitStr)
	if limitStr == "" {
		limit = 10
	} else if err != nil {
		http.Error(w, "invalid limit param", http.StatusBadRequest)
		return
	}

	// Some sample data
	allUsers := []User{
		{UserID: "u1", Name: "Alice"},
		{UserID: "u2", Name: "Bob"},
		{UserID: "u3", Name: "Charlie"},
	}

	// Truncate if limit is smaller
	if limit < len(allUsers) {
		allUsers = allUsers[:limit]
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(allUsers)
}

func getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	userID := mux.Vars(r)["userId"]

	// Some sample data
	users := map[string]User{
		"u1": {UserID: "u1", Name: "Alice"},
		"u2": {UserID: "u2", Name: "Bob"},
	}

	user, ok := users[userID]
	if !ok {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// =======================
// SAMPLE MIDDLEWARES
// =======================

// requestLoggerMiddleware logs the method, path, and duration for every request
func requestLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("[GLOBAL LOG] %s %s took %v", r.Method, r.URL.Path, time.Since(start))
	})
}

// logOperationMiddleware logs that a specific operationId was called
func logOperationMiddleware(operationId string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[OP LOG] operationId=%s called", operationId)
			next.ServeHTTP(w, r)
		})
	}
}
