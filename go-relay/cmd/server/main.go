package main

import (
	"log"
	"net/http"

	"github.com/pritam-ago/go-relay/internal/db"
	"github.com/pritam-ago/go-relay/internal/handlers"
	"github.com/pritam-ago/go-relay/internal/ws"
)

func main() {
	db.Init()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Relay server running...."))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
		w.WriteHeader(http.StatusOK)
	})

	http.HandleFunc("/register-pin", handlers.HandleRegisterPin)
	http.HandleFunc("/check-pin/", handlers.HandleCheckPin)
	http.HandleFunc("/connect-pc/", ws.HandlePCConnect)
	http.HandleFunc("/connect-user/", ws.HandleUserConnect)

	log.Println("🧠 Relay server running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
