package server

import (
	"net"
	"log"
)

func readCommand(conn net.Conn) (string, error) {
	var buffer []byte = make([]byte, 1024)
	n, err := conn.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:n]), nil
}

func respondToClient(conn net.Conn, cmd string) error {
	_, err := conn.Write([]byte(cmd))
	if err != nil {
		return err
	}
	return nil
}
func StartTcpServer() {
	log.Println("starting the tcp server")
	var connected_clients int = 0
	listner, err := net.Listen("tcp", ":4444")
	if err != nil {
		panic(err)
	}
	for {
		connection, err := listner.Accept()
		if err != nil {
			panic(err)
		}
		connected_clients++
		log.Println("New client connected. Total clients: ", connected_clients)
		log.Println("connection recived from remote address", connection.RemoteAddr())
		for {
			command, err := readCommand(connection)
			if err != nil {
				log.Println("Error reading command: ", err)
				connection.Close()
				connected_clients--
				log.Println("Client disconnected. Total clients: ", connected_clients)
				break
			}
			log.Println("command recived", command)
			log.Println("Received command: ", command)
			if err := respondToClient(connection, command); err != nil {
				log.Println("Error responding to client: ", err)
			}
		}
	}
}
