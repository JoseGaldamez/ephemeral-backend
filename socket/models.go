package socket

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Message struct {
	Type     string `json:"type"`
	RoomID   string `json:"idRoom"`
	RoomName string `json:"roomName"`
	UserID string `json:"idClient"`
	Content  string `json:"content"`
}

type Client struct {
	hub *Hub
	conn *websocket.Conn
	send chan []byte
	roomID string
	roomName string // Nombre de la sala, opcional
	userID string
}

type Hub struct {
	rooms map[string]*Room
	mu sync.RWMutex // Mutex para proteger el acceso concurrente a las salas
	register   chan *Client
	unregister chan *Client
	broadcast   chan Message
}

type Room struct {
	ID        string
	Name      string
	clients   map[*Client]bool
	broadcast chan []byte
	register  chan *Client
	unregister chan *Client
}