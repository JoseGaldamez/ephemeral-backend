package socket

import (
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {		
		return true
	},
}

// ServeWs maneja las conexiones WebSocket para los clientes
func ServeWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	// Obtén el roomID y userID de la URL o cabeceras (para ejemplo, de la URL)
	roomID := r.URL.Query().Get("roomID")
	if roomID == "" {
		http.Error(w, "roomID es requerido", http.StatusBadRequest)
		return
	}
	
	// Genera un userID único para cada conexión
	userID := uuid.New().String()[:8] // Solo los primeros 8 caracteres

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Error al actualizar la conexión a WebSocket: %v", err)
		return
	}

	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan []byte, 256), // Buffer para mensajes salientes
		roomID:   roomID,
		roomName: "Sala " + roomID, // Nombre de la sala opcional, puedes personalizarlo
		userID:   userID,
	}

	client.hub.register <- client // Registra el cliente en el hub

	// Permite que múltiples goroutines puedan leer y escribir en la misma conexión
	// pero solo una escritura a la vez.
	go client.writePump() // Goroutine para escribir mensajes al cliente
	go client.readPump()  // Goroutine para leer mensajes del cliente
}
