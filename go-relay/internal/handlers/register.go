package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pritam-ago/go-relay/internal/db"
	"github.com/pritam-ago/go-relay/internal/models"
)

type RegisterPinRequest struct {
	PIN string `json:"pin"`
}

type RegisterPinResponse struct {
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at"`
}

func HandleRegisterPin(w http.ResponseWriter, r *http.Request) {
	log.Printf("Register PIN request received: %s %s", r.Method, r.URL.Path)

	if r.Method != http.MethodPost {
		log.Printf("Invalid method for register PIN: %s", r.Method)
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body RegisterPinRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if body.PIN == "" {
		log.Printf("Empty PIN in request body")
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Attempting to register PIN: %s", body.PIN)

	session, err := models.CreateSession(db.DB, body.PIN)
	if err != nil {
		log.Printf("Failed to create session for PIN %s: %v", body.PIN, err)
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully created session for PIN: %s", body.PIN)
	log.Printf("Session details - PIN: %s, CreatedAt: %v, ExpiresAt: %v",
		session.PIN, session.CreatedAt, session.ExpiresAt)

	resp := RegisterPinResponse{
		Message:   "PIN registered successfully",
		ExpiresAt: session.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		log.Printf("Failed to encode response for PIN %s: %v", body.PIN, err)
	} else {
		log.Printf("Successfully sent registration response for PIN: %s", body.PIN)
	}
}
