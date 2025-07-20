package main

import (
	"log"
	"net/http"

	"github.com/pritam-ago/go-relay/internal/db"
	"github.com/pritam-ago/go-relay/internal/handlers"
	"github.com/pritam-ago/go-relay/internal/ws"
)

// CORS middleware
func corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func main() {
	db.Init()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Relay server running...."))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/register-pin", corsMiddleware(handlers.HandleRegisterPin))
	http.HandleFunc("/check-pin/", corsMiddleware(handlers.HandleCheckPin))
	http.HandleFunc("/connect-pc/", corsMiddleware(ws.HandlePCConnect))
	http.HandleFunc("/connect-user/", corsMiddleware(ws.HandleUserConnect))

	log.Println("🧠 Relay server running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
