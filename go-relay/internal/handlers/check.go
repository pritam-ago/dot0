package handlers

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/pritam-ago/go-relay/internal/db"
	"github.com/pritam-ago/go-relay/internal/models"
)

type CheckPinResponse struct {
	Valid        bool      `json:"valid"`
	ExpiresAt    time.Time `json:"expires_at,omitempty"`
	PCConnected  bool      `json:"pc_connected,omitempty"`
	ErrorMessage string    `json:"error,omitempty"`
}

type GetBaseDirResponse struct {
	BaseDirectory string `json:"base_directory,omitempty"`
	ErrorMessage  string `json:"error,omitempty"`
}

func HandleCheckPin(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract PIN from URL
	pin := strings.TrimPrefix(r.URL.Path, "/check-pin/")
	if pin == "" {
		http.Error(w, "PIN is required", http.StatusBadRequest)
		return
	}

	session, err := models.GetSession(db.DB, pin)
	if err != nil {
		resp := CheckPinResponse{
			Valid:        false,
			ErrorMessage: "PIN not found",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	// Check expiration
	if time.Now().After(session.ExpiresAt) {
		resp := CheckPinResponse{
			Valid:        false,
			ErrorMessage: "PIN expired",
		}
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := CheckPinResponse{
		Valid:       true,
		ExpiresAt:   session.ExpiresAt,
		PCConnected: session.PCConnected,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func HandleGetBaseDir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract PIN from URL
	pin := strings.TrimPrefix(r.URL.Path, "/get-base-dir/")
	if pin == "" {
		http.Error(w, "PIN is required", http.StatusBadRequest)
		return
	}

	session, err := models.GetSession(db.DB, pin)
	if err != nil {
		resp := GetBaseDirResponse{
			ErrorMessage: "PIN not found",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	if session.BaseDirectory == nil {
		resp := GetBaseDirResponse{
			ErrorMessage: "Base directory not set",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}

	resp := GetBaseDirResponse{
		BaseDirectory: *session.BaseDirectory,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
