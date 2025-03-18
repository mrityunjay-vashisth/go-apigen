package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type User struct {
	UserID string `json:"userId"`
	Name   string `json:"name"`
}

func main() {
	// We assume the server is running at localhost:8080

	// 1) Call GET /users?limit=2
	usersURL := "http://localhost:8080/users?limit=3"
	usersResp, err := http.Get(usersURL)
	if err != nil {
		log.Fatalf("Error calling %s: %v", usersURL, err)
	}
	defer usersResp.Body.Close()

	if usersResp.StatusCode != http.StatusOK {
		log.Fatalf("GET /users returned status %d", usersResp.StatusCode)
	}

	var users []User
	if err := json.NewDecoder(usersResp.Body).Decode(&users); err != nil {
		log.Fatalf("Failed to decode /users response: %v", err)
	}

	fmt.Printf("GET /users?limit=2 -> %d users:\n", len(users))
	for _, u := range users {
		fmt.Printf("  %v\n", u)
	}

	// 2) Call GET /users/u1
	userURL := "http://localhost:8080/users/u2"
	singleResp, err := http.Get(userURL)
	if err != nil {
		log.Fatalf("Error calling %s: %v", userURL, err)
	}
	defer singleResp.Body.Close()

	if singleResp.StatusCode == http.StatusNotFound {
		fmt.Println("GET /users/u1 returned 404 (User Not Found)")
	} else if singleResp.StatusCode == http.StatusOK {
		var user User
		if err := json.NewDecoder(singleResp.Body).Decode(&user); err != nil {
			log.Fatalf("Failed to decode /users/u1 response: %v", err)
		}
		fmt.Printf("GET /users/u1 -> user: %#v\n", user)
	} else {
		fmt.Printf("GET /users/u1 returned status %d\n", singleResp.StatusCode)
	}
}
