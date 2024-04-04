package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/commands"
	"github.com/codecrafters-io/redis-starter-go/redis"
	"github.com/codecrafters-io/redis-starter-go/resp"
)

var (
	port        int
	replicaof   string
	RedisServer *redis.Redis
	isSlave     = false
	masterHost  string
	masterport  string
)

func handleconn(conn net.Conn, redis *redis.Redis) {
	defer conn.Close()

	fmt.Println("Connected")
	respHandler := resp.NewRespHandler(conn)
	respWriter := resp.NewRespWriter(conn)
	for {
		value, err := respHandler.ParseAny()
		if err == io.EOF {
			continue
		}
		fmt.Println(value)
		switch value.Typ {
		case resp.ArrayType:

			Reqargs := value.Array
			Comm := Reqargs[0].Bulk
			Comm_Args := Reqargs[1:]
			respValue := commands.Handlers[strings.ToLower(Comm)](redis, Comm_Args)
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
	flag.IntVar(&port, "port", 6379, "Start Server on : ")
	flag.StringVar(&replicaof, "replicaof", "", "Host Ip and Port")
	flag.Parse()
	if len(strings.TrimSpace(replicaof)) != 0 {
		isSlave = true
		masterHost = replicaof
		masterport = flag.Args()[0]
		RedisServer = redis.NewRedisServer(isSlave)
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", masterHost, masterport))
		defer conn.Close()
		if err != nil {
			fmt.Println("Could not connect to Master")
		}
		buf := make([]byte, 1024)
		_, err = conn.Write([]byte("*1\r\n$4\r\nping\r\n"))
		_, err = conn.Read(buf[:])
		if err != nil {
			fmt.Println("Error Reading from connection")
		}

		_, err = conn.Write([]byte(fmt.Sprintf("*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n%s\r\n", strconv.Itoa(port))))
		_, err = conn.Read(buf[:])
		if err != nil {
			fmt.Println("Error Reading from connection")
		}

		_, err = conn.Write([]byte("*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"))
		_, err = conn.Read(buf[:])
		if err != nil {
			fmt.Println("Error Reading from connection")
		}

		_, err = conn.Write([]byte("*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"))
		_, err = conn.Read(buf[:])
		if err != nil {
			fmt.Println("Error Reading from connection")
		}
	} else {
		RedisServer = redis.NewRedisServer(isSlave)
	}
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
		go handleconn(conn, RedisServer)
	}

}

// func main() {
// 	fmt.Println("Logs from your program will appear here!")
// 	var port int
// 	var replicaof string
// 	var Redis *redis.Redis
// 	flag.IntVar(&port, "port", 6379, "Start Server on : ")
// 	flag.StringVar(&replicaof, "replicaof", "", "Host Ip and Port")
// 	flag.Parse()
// 	if len(strings.TrimSpace(replicaof)) != 0 {
// 		Redis = redis.NewRedisSlave()
// 	} else {
// 		Redis = redis.NewRedisMaster()
// 	}
// 	l, err := net.Listen("tcp", ":4000")
// 	if err != nil {
// 		fmt.Printf("Failed to bind to port %d\n", port)
// 		fmt.Printf("%s\n", err)
// 		os.Exit(1)
// 	}
// 	go test.Test()
// 	for {
// 		conn, err := l.Accept()
// 		if err != nil {
// 			fmt.Println("Error accepting connection: ", err.Error())
// 			os.Exit(1)
// 		}
// 		go handleconn(conn, Redis)
// 	}
// }
