package main

import (
	"log"
	"net/http"
	"github.com/pritam-ago/go-relay/db"
)

func main() {
	db.Init()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Relay server is up and running"))
	})

	log.Println("🧠 Relay server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
