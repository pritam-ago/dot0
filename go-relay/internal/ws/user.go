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
			log.Printf("🔌 Disconnected [%s]: %v", tag, err)
			from.Close()
			to.Close()
			break
		}

		log.Printf("📤 Relaying message [%s]: %s", tag, string(msg))

		err = to.WriteMessage(msgType, msg)
		if err != nil {
			log.Printf("❌ Write error [%s]: %v", tag, err)
			from.Close()
			to.Close()
			break
		}

		log.Printf("✅ Message relayed successfully [%s]", tag)
	}
}

func HandleUserConnect(w http.ResponseWriter, r *http.Request) {
	pin := r.URL.Path[len("/connect-user/"):]
	log.Printf("🔍 User connection attempt for PIN: %s", pin)

	if pin == "" {
		log.Printf("❌ Missing PIN in request path: %s", r.URL.Path)
		http.Error(w, "Missing PIN", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 Checking if PC is connected for PIN: %s", pin)
	pcConn, ok := ActivePCConnections[pin]
	if !ok {
		log.Printf("❌ PC not connected for PIN: %s. Active PC connections: %v", pin, getActivePCConnectionPins())
		http.Error(w, "PC not connected for this PIN", http.StatusNotFound)
		return
	}
	log.Printf("✅ PC connection found for PIN: %s", pin)

	// Fetch session from database BEFORE WebSocket upgrade
	log.Printf("🔍 Fetching session from database for PIN: %s", pin)
	session, err := models.GetSession(db.DB, pin)
	if err != nil {
		log.Printf("❌ Database error for PIN %s: %v", pin, err)
		http.Error(w, "PIN not registered or missing base directory", http.StatusForbidden)
		return
	}

	log.Printf("✅ Session found for PIN: %s", pin)
	log.Printf("📊 Session details - BaseDirectory: %v, UserConnected: %v, PCConnected: %v",
		session.BaseDirectory, session.UserConnected, session.PCConnected)

	if session.BaseDirectory == nil {
		log.Printf("❌ Base directory is nil for PIN: %s", pin)
		http.Error(w, "PIN not registered or missing base directory", http.StatusForbidden)
		return
	}

	if session.UserConnected {
		log.Printf("❌ User already connected for PIN: %s", pin)
		http.Error(w, "User already connected to this PIN", http.StatusForbidden)
		return
	}

	log.Printf("🔍 Attempting WebSocket upgrade for PIN: %s", pin)
	userConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ Failed to upgrade user socket for PIN %s: %v", pin, err)
		return
	}
	log.Printf("✅ WebSocket upgrade successful for PIN: %s", pin)

	// Mark user as connected in the database
	log.Printf("🔍 Marking user as connected in database for PIN: %s", pin)
	err = models.MarkUserConnected(db.DB, pin, true)
	if err != nil {
		log.Printf("❌ Failed to mark user as connected in database for PIN %s: %v", pin, err)
		userConn.Close()
		return
	}
	log.Printf("✅ User marked as connected in database for PIN: %s", pin)

	// Store user connection
	ActiveUserConnections[pin] = userConn
	log.Printf("✅ User connection stored for PIN: %s", pin)

	// Cleanup on disconnect
	defer func() {
		log.Printf("🔌 User disconnecting from PIN: %s", pin)
		models.MarkUserConnected(db.DB, pin, false)
		delete(ActiveUserConnections, pin)
		userConn.Close()
		log.Printf("✅ User cleanup completed for PIN: %s", pin)
	}()

	log.Printf("🚀 Starting message relay for PIN: %s", pin)
	// Start relaying messages
	go relayMessages(userConn, pcConn, "user→pc")
	relayMessages(pcConn, userConn, "pc→user") // this will block until PC disconnects or closes socket
}

// Helper function to get active PC connection PINs for debugging
func getActivePCConnectionPins() []string {
	pins := make([]string, 0, len(ActivePCConnections))
	for pin := range ActivePCConnections {
		pins = append(pins, pin)
	}
	return pins
}
