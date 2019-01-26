package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

// Initializes state and listens on the given host/port
func Listen(host string, port string) {
	sse := newSSEBroker()

	messageService := newMessageService(func(msg Message) {
		str, _ := json.Marshal(msg)
		sse.Notifier <- []byte(str)
	})

	router := mux.NewRouter()

	router.HandleFunc("/sse", sse.connHandler)
	router.HandleFunc("/api/messages", messageService.createHandler).Methods("POST")
	router.HandleFunc("/api/messages", messageService.indexHandler).Methods("GET")

	// Static files
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./pkg/waffle/ui/")))

	hostWithPort := fmt.Sprintf("%s:%s", host, port)

	log.Printf("[waffle/server] Listening on http://%s\n", hostWithPort)
	log.Fatal(http.ListenAndServe(hostWithPort, router))
}
