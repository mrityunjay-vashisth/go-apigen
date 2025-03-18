# go-apigen

A lightweight Go library that helps generate HTTP server handlers from OpenAPI specifications.

## Overview

`go-apigen` makes it easy to build REST APIs by automatically wiring up HTTP handlers to routes defined in your OpenAPI specification. Rather than manually connecting each endpoint to its handler function, this library reads your OpenAPI spec and creates a fully configured Gorilla Mux router with all your endpoints properly connected.

## Features

- **OpenAPI-Driven Development**: Design your API with an OpenAPI specification and automatically generate the matching router
- **Minimal Boilerplate**: Focus on writing your business logic instead of routing code
- **Type Safety**: Provides a clean mapping between OpenAPI operationIds and handler functions
- **Gorilla Mux Integration**: Built on top of the widely-used Gorilla Mux router
- **Automatic 501 Handlers**: Routes without an implementation get a proper "Not Implemented" response

## Requirements

- Go 1.24 or higher
- OpenAPI 3.0 specification file (YAML or JSON)

## Installation

```bash
go get github.com/mrityunjay-vashisth/go-apigen
```

## Quick Start

1. Define your API in an OpenAPI specification file (e.g., `openapi.yaml`)
2. Parse the OpenAPI file using `generator.ParseOpenAPIFile`
3. Create handler functions for each operationId
4. Map the operationIds to your handler functions using `generator.OperationMap`
5. Generate the router with `generator.GenerateMuxRouter`

```go
package main

import (
	"log"
	"net/http"

	"github.com/mrityunjay-vashisth/go-apigen/pkg/generator"
)

func main() {
	// Parse the OpenAPI file
	doc, err := generator.ParseOpenAPIFile("openapi.yaml")
	if err != nil {
		log.Fatalf("Failed to parse OpenAPI spec: %v", err)
	}

	// Create a map of operationId -> handler function
	ops := generator.OperationMap{
		"listUsers":   listUsersHandler,
		"getUserById": getUserByIdHandler,
	}

	// Build the Gorilla Mux router
	router, err := generator.GenerateMuxRouter(doc, ops)
	if err != nil {
		log.Fatalf("Failed to generate router: %v", err)
	}

	// Start the server
	log.Println("Server listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

// Define your handler functions...
```

## Example

### OpenAPI Specification (openapi.yaml)

```yaml
openapi: 3.0.2
info:
  title: Simple User API
  version: 1.0.0

paths:
  /users:
    get:
      operationId: listUsers
      summary: List users
      parameters:
        - name: limit
          in: query
          schema:
            type: integer
            default: 10
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                type: array
                items:
                  $ref: '#/components/schemas/User'

  /users/{userId}:
    get:
      operationId: getUserById
      parameters:
        - name: userId
          in: path
          required: true
          schema:
            type: string
      responses:
        '200':
          description: OK
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/User'
        '404':
          description: User Not Found

components:
  schemas:
    User:
      type: object
      properties:
        userId:
          type: string
        name:
          type: string
```

### Server Implementation

See the full example in the `cmd/server/main.go` file:

```go
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
```

### Client Implementation

See the client example in `cmd/client/main.go` for how to consume the API.

## API Reference

### generator.OperationMap

`OperationMap` is a map from OpenAPI operationId to an HTTP handler function.

```go
type OperationMap map[string]http.HandlerFunc
```

### generator.ParseOpenAPIFile(filename string) (*openapi3.T, error)

Parses an OpenAPI specification file (YAML or JSON) and returns the parsed document.

### generator.GenerateMuxRouter(doc *openapi3.T, ops OperationMap) (*mux.Router, error)

Generates a Gorilla Mux router based on the provided OpenAPI document and operation map.

### generator.GetPathParam(r *http.Request, key string) string

Helper function to get path parameters from a request.

## Benefits of Using go-apigen

1. **Consistency with API Design**: Your implementation is automatically aligned with your API specification.
2. **Reduced Boilerplate**: No need to manually set up routes for each endpoint.
3. **Better Developer Experience**: Focus on writing the actual handler logic.
4. **Documentation as Code**: Your OpenAPI spec serves as both documentation and the source of truth for routing.
5. **Graceful Handling of Unimplemented Endpoints**: Missing handlers get a proper 501 response.

## License

[MIT License](LICENSE)

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.