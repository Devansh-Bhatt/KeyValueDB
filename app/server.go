package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/commands"
	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/store"
)

func handleconn(conn net.Conn) {
	defer conn.Close()

	fmt.Println("Connected")
	respHandler := resp.NewRespHandler(conn)
	respWriter := resp.NewRespWriter(conn)
	clientDb := store.NewDb()
	for {
		value, err := respHandler.ParseAny()
		fmt.Println(value)
		if err != nil {
			fmt.Println("From Here", err)
		}
		switch value.Typ {
		case resp.ArrayType:

			Reqargs := value.Array
			Comm := Reqargs[0].Bulk
			// fmt.Println(Comm)
			Comm_Args := Reqargs[1:]
			respValue := commands.Handlers[strings.ToLower(Comm)](clientDb, Comm_Args)
			respWriter.Write(respValue)
		case resp.StringType:
			switch strings.ToLower(value.Str) {
			case "ping":
				respValue := resp.Value{
					Typ: resp.StringType,
					Str: "PONG",
				}

				respWriter.Write(respValue)
			}

		}
	}

}

func main() {
	fmt.Println("Logs from your program will appear here!")
	var port int

	flag.IntVar(&port, "port", 6379, "Start Server on : ")
	flag.Parse()

	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d", port)
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

// func maintest() {
// 	fmt.Println("Logs from your program will appear here!")

// 	l, err := net.Listen("tcp", ":4000")
// 	if err != nil {
// 		fmt.Println("Failed to bind to port 6379")
// 		os.Exit(1)
// 	}
// 	go test.Test()

// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection: ", err.Error())
// 			os.Exit(1)
// 		}
// 		go handleconn(conn)
// 	}

// }
