package ws

import (
	"log"
	"net/http"

	"github.com/pritam-ago/go-relay/internal/db"
	"github.com/pritam-ago/go-relay/internal/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ActivePCConnections = make(map[string]*websocket.Conn)

func HandlePCConnect(w http.ResponseWriter, r *http.Request) {
	pin := r.URL.Path[len("/connect-pc/"):]
	if pin == "" {
		http.Error(w, "Missing PIN", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ WebSocket upgrade failed:", err)
		return
	}

	log.Println("✅ PC connected with PIN:", pin)
	ActivePCConnections[pin] = conn

	// Update DB: PC is now connected
	err = models.UpdatePCConnectionStatus(db.DB, pin, true)
	if err != nil {
		log.Println("⚠️ Failed to update DB:", err)
	}

	// Listen for disconnect
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("💀 PC disconnected:", err)
			delete(ActivePCConnections, pin)

			// Mark PC as disconnected
			_ = models.UpdatePCConnectionStatus(db.DB, pin, false)
			break
		}
	}
}
