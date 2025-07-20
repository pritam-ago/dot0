package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pritam-ago/go-relay/internal/db"
	"github.com/pritam-ago/go-relay/internal/models"
)

func relayMessages(from, to *websocket.Conn, tag string) {
	for {
		msgType, msg, err := from.ReadMessage()
		if err != nil {
			log.Println("🔌 Disconnected ["+tag+"]:", err)
			from.Close()
			to.Close()
			break
		}

		err = to.WriteMessage(msgType, msg)
		if err != nil {
			log.Println("❌ Write error ["+tag+"]:", err)
			from.Close()
			to.Close()
			break
		}
	}
}

func HandleUserConnect(w http.ResponseWriter, r *http.Request) {
	pin := r.URL.Path[len("/connect-user/"):]
	if pin == "" {
		http.Error(w, "Missing PIN", http.StatusBadRequest)
		return
	}

	pcConn, ok := ActivePCConnections[pin]
	if !ok {
		http.Error(w, "PC not connected for this PIN", http.StatusNotFound)
		return
	}

	// Fetch session from database BEFORE WebSocket upgrade
	session, err := models.GetSession(db.DB, pin)
	if err != nil || session.BaseDirectory == nil {
		http.Error(w, "PIN not registered or missing base directory", http.StatusForbidden)
		return
	}

	if session.UserConnected {
		http.Error(w, "User already connected to this PIN", http.StatusForbidden)
		return
	}

	userConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Failed to upgrade user socket:", err)
		return
	}

	// Mark user as connected in the database
	models.MarkUserConnected(db.DB, pin, true)
	log.Println("User connected to PIN:", pin)

	// Store user connection
	ActiveUserConnections[pin] = userConn

	// Cleanup on disconnect
	defer func() {
		log.Println("User disconnected from PIN:", pin)
		models.MarkUserConnected(db.DB, pin, false)
		delete(ActiveUserConnections, pin)
		userConn.Close()
	}()

	// Start relaying messages
	go relayMessages(userConn, pcConn, "user→pc")
	relayMessages(pcConn, userConn, "pc→user") // this will block until PC disconnects or closes socket
}
