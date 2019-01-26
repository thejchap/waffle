package main

import (
	"os"

	"github.com/thejchap/waffle/pkg/waffle"
)

const DefaultHost = "0.0.0.0"
const DefaultPort = "3000"

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = DefaultPort
	}

	server.Listen(DefaultHost, port)
}
