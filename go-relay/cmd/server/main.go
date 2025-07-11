package main

import (
	"log"
	"net/http"

	"github.com/pritam-ago/go-relay/internal/db"
	"github.com/pritam-ago/go-relay/internal/handlers"
)

func main() {
	db.Init()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Relay server running"))
	})

	http.HandleFunc("/register-pin", handlers.HandleRegisterPin)

	log.Println("🧠 Relay server running on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
