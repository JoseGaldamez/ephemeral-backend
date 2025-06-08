package server

import (
	"fmt"
	"net/http"
)


func HandleMain(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Received request for:", r.URL.Path)
	fmt.Fprintln(w, "Go Server is running!")
}