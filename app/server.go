package main

import (
	"flag"
	"fmt"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/redis"
)

var (
	port       int
	replicaof  string
	masterHost string
	masterport string
)

func main() {
	fmt.Println("Logs from your program will appear here!")
	flag.IntVar(&port, "port", 6380, "Start Server on : ")
	flag.StringVar(&replicaof, "replicaof", "", "Host Ip and Port")
	flag.Parse()
	if len(strings.TrimSpace(replicaof)) != 0 {
		repl_parts := strings.Split(replicaof, " ")
		masterHost = repl_parts[0]
		masterport = repl_parts[1]
		// RedisServer := redis.New_Redis_Slave_Server()
		conn, err := net.Dial("tcp", fmt.Sprintf("%s:%s", masterHost, masterport))
		defer conn.Close()
		if err != nil {
			fmt.Println("Could not connect to Master")
		}
		_, RedisServer := redis.Init_Server(true)
		RedisServer.Start_Server(port, conn, masterHost, masterport)
	} else {

		// RedisServer := redis.New_Redis_Master_Server()
		RedisServer, _ := redis.Init_Server(false)
		RedisServer.Start_Server(port)

	}

}
