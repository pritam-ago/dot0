package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
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
	// Load .env file from project root (two directories up from cmd/server)
	if err := godotenv.Load("../../.env"); err != nil {
		// Try loading from current directory as fallback
		if err := godotenv.Load(); err != nil {
			log.Println("No .env file found, using system environment variables")
		}
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Println("PORT not set, using default port 8080")
		port = "8080"
	}
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
	http.HandleFunc("/get-base-dir/", corsMiddleware(handlers.HandleGetBaseDir))
	http.HandleFunc("/connect-pc/", corsMiddleware(ws.HandlePCConnect))
	http.HandleFunc("/connect-user/", corsMiddleware(ws.HandleUserConnect))

	log.Println("Relay server running on :" + port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
