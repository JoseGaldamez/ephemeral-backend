package main

import (
	"log"
	"net/http"

	"github.com/JoseGaldamez/ephemeral-backend/socket"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html") // Sirve el archivo HTML principal
	})

	hub := socket.NewHub() // Crea un nuevo hub para manejar las conexiones
	go hub.Run() // Inicia el hub en una goroutine

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		socket.ServeWs(hub, w, r)
	})

	log.Println("Starting server on port 8080")
	error := http.ListenAndServe(":8080", nil)
	if error != nil {
		log.Fatal("Error al iniciar el servidor: ", error)
	} 
}
