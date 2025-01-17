package main

import (
	"Groq2API/initialize/auth"
	"Groq2API/initialize/stream"
	"Groq2API/initialize/user"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type ChatCompletionRequest struct {
	RefreshToken string `json:"refresh_token"`
}

func chatCompletionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var req ChatCompletionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jwt, err := auth.FetchJWT(req.RefreshToken)
	if err != nil {
		log.Printf("Error fetching JWT: %v", err)
		http.Error(w, "Failed to fetch JWT", http.StatusInternalServerError)
		return
	}

	orgID, err := user.FetchUserProfile(jwt)
	if err != nil {
		log.Printf("Error fetching user profile: %v", err)
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}

	response, err := stream.FetchStream(jwt, orgID) // Make sure to adjust the FetchStream function to return the response instead of printing it.
	if err != nil {
		log.Printf("Error fetching stream: %v", err)
		http.Error(w, "Failed to fetch stream", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response) 

func main() {
	http.HandleFunc("/v1/chat/completions", chatCompletionsHandler)

	fmt.Println("Server is listening on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
