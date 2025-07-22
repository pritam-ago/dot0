package ws

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/pritam-ago/go-relay/internal/db"
	"github.com/pritam-ago/go-relay/internal/models"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type FileInfo struct {
	Name        string `json:"name"`
	Path        string `json:"path"`
	IsDirectory bool   `json:"is_directory"`
	Size        *int64 `json:"size,omitempty"`
	Modified    string `json:"modified,omitempty"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var ActivePCConnections = make(map[string]*websocket.Conn)
var ActiveUserConnections = make(map[string]*websocket.Conn)

func HandlePCConnect(w http.ResponseWriter, r *http.Request) {
	// Add CORS headers for WebSocket connections
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	pin := r.URL.Path[len("/connect-pc/"):]
	log.Printf("🔍 PC connection attempt for PIN: %s", pin)

	if pin == "" {
		log.Printf("❌ Missing PIN in PC request path: %s", r.URL.Path)
		http.Error(w, "Missing PIN", http.StatusBadRequest)
		return
	}

	log.Printf("🔍 Attempting WebSocket upgrade for PC with PIN: %s", pin)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("❌ WebSocket upgrade failed for PC with PIN %s: %v", pin, err)
		return
	}

	log.Printf("✅ PC WebSocket upgrade successful for PIN: %s", pin)
	log.Printf("✅ PC connected with PIN: %s", pin)
	ActivePCConnections[pin] = conn
	log.Printf("📊 Active PC connections after adding: %v", getActivePCConnectionPins())

	// Update DB: PC is now connected
	log.Printf("🔍 Updating PC connection status in database for PIN: %s", pin)
	err = models.UpdatePCConnectionStatus(db.DB, pin, true)
	if err != nil {
		log.Printf("⚠️ Failed to update DB for PC connection PIN %s: %v", pin, err)
	} else {
		log.Printf("✅ PC connection status updated in database for PIN: %s", pin)
	}

	// Listen for messages
	log.Printf("🎧 Starting message listener for PC with PIN: %s", pin)
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Printf("💀 PC disconnected for PIN %s: %v", pin, err)
			delete(ActivePCConnections, pin)
			log.Printf("📊 Active PC connections after removal: %v", getActivePCConnectionPins())
			_ = models.UpdatePCConnectionStatus(db.DB, pin, false)
			break
		}

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Printf("❌ Invalid message from PC with PIN %s: %v", pin, err)
			continue
		}

		log.Printf("📨 Received message from PC with PIN %s: %s", pin, msg.Type)

		// Handle different message types
		switch msg.Type {
		case "register_base_dir":
			var payload struct {
				Path string `json:"path"`
			}
			b, _ := json.Marshal(msg.Data)
			_ = json.Unmarshal(b, &payload)

			log.Printf("🔍 Registering base directory for PIN %s: %s", pin, payload.Path)
			err := models.UpdateBaseDirectory(db.DB, pin, payload.Path)
			if err != nil {
				log.Printf("❌ Failed to store base dir for PIN %s: %v", pin, err)
			} else {
				log.Printf("📂 Base directory registered for PIN %s → %s", pin, payload.Path)
			}

		case "list_files":
			// Relay to user if connected
			if userConn, exists := ActiveUserConnections[pin]; exists {
				log.Printf("📤 Relaying list_files to user for PIN: %s", pin)
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Printf("❌ Failed to relay list_files to user for PIN %s: %v", pin, err)
				}
			} else {
				log.Printf("⚠️ No user connected for list_files, PIN: %s", pin)
			}

		case "file_content":
			// Relay file content to user
			if userConn, exists := ActiveUserConnections[pin]; exists {
				log.Printf("📤 Relaying file_content to user for PIN: %s", pin)
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Printf("❌ Failed to relay file_content to user for PIN %s: %v", pin, err)
				}
			} else {
				log.Printf("⚠️ No user connected for file_content, PIN: %s", pin)
			}

		case "upload_file":
			// This message is coming from PC to user (response), relay it
			if userConn, exists := ActiveUserConnections[pin]; exists {
				log.Printf("📤 Relaying upload response to user for PIN: %s", pin)
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Printf("❌ Failed to relay upload response to user for PIN %s: %v", pin, err)
				}
			} else {
				log.Printf("⚠️ No user connected for upload response, PIN: %s", pin)
			}

		case "upload_response":
			// Relay upload response from PC back to user
			if userConn, exists := ActiveUserConnections[pin]; exists {
				log.Printf("📥 Relaying upload response to user for PIN: %s", pin)
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Printf("❌ Failed to relay upload_response to user for PIN %s: %v", pin, err)
				}
			} else {
				log.Printf("⚠️ No user connected for upload response, PIN: %s", pin)
			}

		default:
			// Relay any other messages to user if connected
			if userConn, exists := ActiveUserConnections[pin]; exists {
				log.Printf("📤 Relaying %s message to user for PIN: %s", msg.Type, pin)
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Printf("❌ Failed to relay %s message to user for PIN %s: %v", msg.Type, pin, err)
				}
			} else {
				log.Printf("⚠️ No user connected for %s message, PIN: %s", msg.Type, pin)
			}
		}
	}
}

// Add a function to handle messages from user to PC
func relayUserToPC(userConn, pcConn *websocket.Conn, pin string) {
	for {
		msgType, msg, err := userConn.ReadMessage()
		if err != nil {
			log.Printf("🔌 User disconnected for PIN %s: %v", pin, err)
			break
		}

		log.Printf("📨 Received message from user for PIN %s: %s", pin, string(msg))

		// Relay the message to PC
		err = pcConn.WriteMessage(msgType, msg)
		if err != nil {
			log.Printf("❌ Failed to relay message to PC for PIN %s: %v", pin, err)
			break
		}

		log.Printf("✅ Message relayed to PC for PIN: %s", pin)
	}
}
