package main

import (
	"fmt"
	"log"

	"github.com/appu900/lokey/server"
)

const pink = "\033[38;5;213m"
const reset = "\033[0m"

func main() {
	fmt.Println(pink + `
██╗      ██████╗ ██╗  ██╗███████╗██╗   ██╗
██║     ██╔═══██╗██║ ██╔╝██╔════╝╚██╗ ██╔╝
██║     ██║   ██║█████╔╝ █████╗   ╚████╔╝ 
██║     ██║   ██║██╔═██╗ ██╔══╝    ╚██╔╝  
███████╗╚██████╔╝██║  ██╗███████╗   ██║   
╚══════╝ ╚═════╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   
` + reset)

	log.Println("starting lokey")

	server.StartTcpServer()
}