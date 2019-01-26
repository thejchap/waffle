package main

import (
	"os"

	"github.com/thejchap/waffle/pkg/waffle"
)

const defaultHost = "0.0.0.0"
const defaultPort = "3000"

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}

	server.Listen(defaultHost, port)
}
