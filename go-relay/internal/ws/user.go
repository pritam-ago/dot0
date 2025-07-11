package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func relayMessages(from, to *websocket.Conn, tag string) {
	for {
		msgType, msg, err := from.ReadMessage()
		if err != nil {
			log.Println("🔌 Disconnected [" + tag + "]:", err)
			from.Close()
			to.Close()
			break
		}

		err = to.WriteMessage(msgType, msg)
		if err != nil {
			log.Println("❌ Write error [" + tag + "]:", err)
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

	userConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("❌ Failed to upgrade user socket:", err)
		return
	}

	log.Println("👤 User connected to PIN:", pin)

	go relayMessages(userConn, pcConn, "user→pc")
	go relayMessages(pcConn, userConn, "pc→user")
}
