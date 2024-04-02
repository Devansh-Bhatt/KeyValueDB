package main

import (
	"flag"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/commands"
	"github.com/codecrafters-io/redis-starter-go/redis"
	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/test"
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
	var port int
	var replicaof string
	var Redis *redis.Redis
	var RedisSlave *redis.RedisSlave
	flag.IntVar(&port, "port", 6379, "Start Server on : ")
	flag.StringVar(&replicaof, "replicaof", "", "Host Ip and Port")
	flag.Parse()
	if len(strings.TrimSpace(replicaof)) != 0 {
		Master := replicaof
		MasterPort := flag.Args()[0]
		RedisSlave = redis.NewRedisSlave()

		conn, err := RedisSlave.ConnectMaster(Master, MasterPort)
		// respwriter := resp.NewRespWriter(conn)
		if err != nil {
			fmt.Printf(err.Error())
		}
		_, err = conn.Write([]byte("*1\r\n$4\r\nping\r\n"))

	} else {
		Redis = redis.NewRedisMaster()
	}
	fmt.Println(replicaof)
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d", port)
		os.Exit(1)
	}
	go test.Test()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go handleconn(conn, Redis)
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
