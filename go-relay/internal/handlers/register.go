package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/pritam-ago/go-relay/internal/models"
	"github.com/pritam-ago/go-relay/internal/db"
)

type RegisterPinRequest struct {
	PIN string `json:"pin"`
}

type RegisterPinResponse struct {
	Message   string    `json:"message"`
	ExpiresAt time.Time `json:"expires_at"`
}

func HandleRegisterPin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var body RegisterPinRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil || body.PIN == "" {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	session, err := models.CreateSession(db.DB, body.PIN)
	if err != nil {
		log.Println("❌ DB error:", err) // 👈 Add this
		http.Error(w, "Failed to create session", http.StatusInternalServerError)
		return
	}

	resp := RegisterPinResponse{
		Message:   "PIN registered successfully",
		ExpiresAt: session.ExpiresAt,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
