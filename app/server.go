package main

import (
	"fmt"
	// Uncomment this block to pass the first stage
	"net"
	"os"

	"github.com/codecrafters-io/redis-starter-go/app/parser"
)

func handleconn(conn net.Conn) {
	defer conn.Close()
	// pong := "+PONG\r\n"
	// conn.Write([]byte(pong))
	for {
		buf := make([]byte, 124)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println("Could not read the client messagex")
		}
		response := parser.MainParser(buf)
		conn.Write(response)
	}

}

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

	// Uncomment this block to pass the first stage
	//
	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleconn(conn)
	}

}
