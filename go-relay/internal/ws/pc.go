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

	// Listen for messages
	for {
		_, raw, err := conn.ReadMessage()
		if err != nil {
			log.Println("💀 PC disconnected:", err)
			delete(ActivePCConnections, pin)
			_ = models.UpdatePCConnectionStatus(db.DB, pin, false)
			break
		}

		var msg Message
		if err := json.Unmarshal(raw, &msg); err != nil {
			log.Println("❌ Invalid message from PC:", err)
			continue
		}

		// Handle different message types
		switch msg.Type {
		case "register_base_dir":
			var payload struct {
				Path string `json:"path"`
			}
			b, _ := json.Marshal(msg.Data)
			_ = json.Unmarshal(b, &payload)

			err := models.UpdateBaseDirectory(db.DB, pin, payload.Path)
			if err != nil {
				log.Println("❌ Failed to store base dir:", err)
			} else {
				log.Println("📂 Base directory registered for", pin, "→", payload.Path)
			}

		case "list_files":
			// Relay to user if connected
			if userConn, exists := ActiveUserConnections[pin]; exists {
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Println("❌ Failed to relay list_files to user:", err)
				}
			}

		case "file_content":
			// Relay file content to user
			if userConn, exists := ActiveUserConnections[pin]; exists {
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Println("❌ Failed to relay file_content to user:", err)
				}
			}

		case "upload_file":
			// Relay upload request to user
			if userConn, exists := ActiveUserConnections[pin]; exists {
				log.Println("📤 Relaying upload request to user for PIN:", pin)
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Println("❌ Failed to relay upload_file to user:", err)
				}
			} else {
				log.Println("⚠️ No user connected for upload request, PIN:", pin)
			}

		case "upload_response":
			// Relay upload response from PC back to user
			if userConn, exists := ActiveUserConnections[pin]; exists {
				log.Println("📥 Relaying upload response to user for PIN:", pin)
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Println("❌ Failed to relay upload_response to user:", err)
				}
			} else {
				log.Println("⚠️ No user connected for upload response, PIN:", pin)
			}

		default:
			// Relay any other messages to user if connected
			if userConn, exists := ActiveUserConnections[pin]; exists {
				err := userConn.WriteMessage(websocket.TextMessage, raw)
				if err != nil {
					log.Println("❌ Failed to relay message to user:", err)
				}
			}
		}
	}
}
