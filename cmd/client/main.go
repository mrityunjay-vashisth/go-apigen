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
	// 1) Call GET /users?limit=2
	resp, err := http.Get("http://localhost:8080/users?limit=2")
	if err != nil {
		log.Fatalf("Failed to call /users: %v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Fatalf("GET /users status %d", resp.StatusCode)
	}

	var users []User
	if err := json.NewDecoder(resp.Body).Decode(&users); err != nil {
		log.Fatalf("Error decoding /users response: %v", err)
	}
	fmt.Printf("GET /users?limit=2 => %d users\n", len(users))

	// 2) Call GET /users/u1
	single, err := http.Get("http://localhost:8080/users/u1")
	if err != nil {
		log.Fatalf("Failed to call /users/u1: %v", err)
	}
	defer single.Body.Close()
	if single.StatusCode == http.StatusOK {
		var u User
		if err := json.NewDecoder(single.Body).Decode(&u); err != nil {
			log.Fatalf("Error decoding /users/u1 response: %v", err)
		}
		fmt.Printf("GET /users/u1 => user: %v\n", u)
	} else if single.StatusCode == http.StatusNotFound {
		fmt.Println("GET /users/u1 => 404 Not Found")
	} else {
		fmt.Printf("GET /users/u1 => status %d\n", single.StatusCode)
	}
}
