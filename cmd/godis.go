package main

import (
	"bufio"
	"fmt"
	"net"

	parser "github.com/leachj/godis"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	for {
		result, err := parser.ParseRESP(reader)
		if err != nil {
			fmt.Println("Connection closed")
			return
		}
		fmt.Printf("Received message: %v \n", result)

		switch value := result.(type) {

		case []interface{}:
			fmt.Printf("Received array %s \n", value[0])
			switch first := value[0].(type) {
			case []byte:
				fmt.Printf("Received bytes %s \n", first)
				if string(first) == "PING" {
					enc, _ := parser.GenerateRESP("PONG")
					conn.Write([]byte(enc))
				} else if string(first) == "ECHO" {
					enc, _ := parser.GenerateRESP(value[1])
					conn.Write([]byte(enc))
				} else {
					conn.Write([]byte("-ERR unknown command\r\n"))
				}
			default:
				conn.Write([]byte("-ERR unknown command\r\n"))

			}
		default:
			conn.Write([]byte("-ERR unknown command\r\n"))
		}

	}
}

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
		return
	}
	defer listener.Close()
	fmt.Println("Server is listening on port 6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err)
			continue
		}
		go handleConnection(conn)
	}
}
