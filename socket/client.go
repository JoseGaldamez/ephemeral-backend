package socket

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

// readPump procesa los mensajes entrantes del cliente
func (client *Client) readPump() {
	defer func() {
		client.hub.unregister <- client // Notifica al hub que el cliente se ha desconectado
		client.conn.Close()
	}()

	client.conn.SetReadLimit(maxMessageSize)
	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error { client.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

	for {
		_, message, err := client.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Error de lectura de WebSocket: %v", err)
			}
			break
		}

		// Deserializa el mensaje del cliente
		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("Error al deserializar mensaje del cliente: %v", err)
			continue
		}

		// Asigna el ID de la sala y el usuario al mensaje (el cliente puede no enviarlo correctamente)
		msg.RoomID = client.roomID
		msg.UserID = client.userID
		msg.RoomName = client.roomName // Asigna el nombre de la sala, opcional
		msg.Type = "chat" // Por defecto, es un mensaje de chat del usuario

		// Envía el mensaje al canal de broadcast global del hub para enrutamiento
		client.hub.broadcast <- msg
	}
}

// writePump envía mensajes al cliente
func (client *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()
	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// El hub cerró el canal.
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}

			w.Write(message)

			// Añade mensajes en cola
			n := len(client.send)
			for i := 0; i < n; i++ {
				w.Write(<-client.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// Envía un ping para mantener la conexión viva
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}