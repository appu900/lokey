package main

import (
	"log"

	"github.com/appu900/lokey/server"
)

func main() {
	log.Println("starting lokey")
	server.StartTcpServer()
}
