package main

import (
	"log"
	"net/http"

	"github.com/JoseGaldamez/ephemeral-backend/server"
)

func main() {
	http.HandleFunc("/", server.HandleMain )

	log.Println("Starting server on port 8080")
	error := http.ListenAndServe(":8080", nil)
	if error != nil {
		log.Fatal("Error al iniciar el servidor: ", error)
	} 
}
