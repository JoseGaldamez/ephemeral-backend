package socket

import (
	"encoding/json"
	"log"
)

func NewHub() *Hub {
	return &Hub{
		rooms:       make(map[string]*Room),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		broadcast:   make(chan Message),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			// Un cliente quiere unirse a una sala.
			// Creamos la sala si no existe.
			h.mu.Lock() // Bloquea para modificar el mapa de salas
			room, ok := h.rooms[client.roomID]
			if !ok {
				room = NewRoom(client.roomID, client.roomName)
				h.rooms[client.roomID] = room
				go room.Run() // Inicia la goroutine para la nueva sala
				log.Printf("Sala '%s' creada.\n", room.Name)
			}
			h.mu.Unlock() // Desbloquea
			
			room.register <- client // Envía el cliente al canal de registro de la sala
			log.Printf("Cliente %s unido a la sala '%s'\n", client.userID, room.Name)


		case client := <-h.unregister:
			// Un cliente quiere salir de una sala o se ha desconectado
			h.mu.RLock() // Bloquea en modo lectura (varias goroutines pueden leer a la vez)
			room, ok := h.rooms[client.roomID]
			h.mu.RUnlock() // Desbloquea

			if ok {
				room.unregister <- client // Envía el cliente al canal de desregistro de la sala
				log.Printf("Cliente %s ha dejado la sala '%s'\n", client.userID, room.Name)
			}
			// Puedes añadir lógica para cerrar la sala si se queda sin clientes


		case message := <-h.broadcast:
			// Un mensaje debe ser enviado a una sala específica
			h.mu.RLock() // Bloquea en modo lectura
			room, ok := h.rooms[message.RoomID]
			h.mu.RUnlock() // Desbloquea

			if ok {
				// Serializa el mensaje a JSON y envíalo al canal de broadcast de la sala
				jsonMessage, err := json.Marshal(message)
				if err != nil {
					log.Printf("Error al serializar el mensaje: %v", err)
					continue
				}

				room.broadcast <- jsonMessage
			} else {
				log.Printf("Error: Sala '%s' no encontrada para el mensaje de broadcast.\n", message.RoomID)
			}
		}
	}
}