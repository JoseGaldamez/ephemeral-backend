package socket

import (
	"encoding/json"
	"fmt"
	"log"
)

func NewRoom(id string, name string) *Room {
	return &Room{
		ID:         id,
		Name:       name,
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (room *Room) Run() {
	for {
		select {
		case client := <-room.register:
			room.clients[client] = true
			log.Printf("Cliente %s se ha unido a la sala %s", client.userID, room.Name)
			// Opcional: Notificar a los otros clientes de la sala que alguien se unió
			joinMsg := Message{
				Type:    "system",
				RoomID:  room.ID,
				RoomName: room.Name,
				Content: fmt.Sprintf("%s se ha unido a la sala.", client.userID),
			}
			jsonJoinMsg, _ := json.Marshal(joinMsg)
			for singleClient := range room.clients {
				if singleClient != client { // No enviarlo al propio cliente que se une
					select {
					case singleClient.send <- jsonJoinMsg:
					default:
						close(singleClient.send)
						delete(room.clients, singleClient)
					}
				}
			}

		case client := <-room.unregister:
			if _, ok := room.clients[client]; ok {
				delete(room.clients, client)
				close(client.send) // Cierra el canal de envío del cliente
				log.Printf("Cliente %s ha salido de la sala %s", client.userID, room.Name)
				// Opcional: Notificar a los otros clientes de la sala que alguien se fue
				leaveMsg := Message{
					Type:    "system",
					RoomID:  room.ID,
					Content: fmt.Sprintf("%s ha dejado la sala.", client.userID),
				}
				jsonLeaveMsg, _ := json.Marshal(leaveMsg)
				for c := range room.clients {
					select {
					case c.send <- jsonLeaveMsg:
					default:
						close(c.send)
						delete(room.clients, c)
					}
				}
			}

		case message := <-room.broadcast:
			for client := range room.clients {
				select {
				case client.send <- message:
				default:
					// Si el canal de envío está bloqueado, el cliente no está leyendo o está desconectado
					close(client.send)
					delete(room.clients, client) // Remueve el cliente inactivo
				}
			}
		}
	}
}