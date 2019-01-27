package main

import (
	"os"

	"github.com/pkg/profile"

	"github.com/thejchap/waffle/pkg/waffle"
)

const defaultHost = "0.0.0.0"
const defaultPort = "3000"

func main() {
	env := os.Getenv("APP_ENV")

	if env != "production" {
		defer profile.Start(profile.MemProfile).Stop()
	}

	port := os.Getenv("PORT")

	if port == "" {
		port = defaultPort
	}

	server.Listen(defaultHost, port)
}
