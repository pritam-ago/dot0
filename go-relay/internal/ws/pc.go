package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // allow any origin for now
	},
}

// ActivePCConnections holds active PC WebSocket sessions by PIN
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

	// Keep PC socket alive (or handle disconnects)
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println("💀 PC disconnected:", err)
			delete(ActivePCConnections, pin)
			break
		}
	}
}
