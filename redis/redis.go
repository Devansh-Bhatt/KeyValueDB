package redis

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/commands"
	"github.com/codecrafters-io/redis-starter-go/repl"
	"github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/store"
	"github.com/codecrafters-io/redis-starter-go/util"
)

const (
	REPLCONF_1 = "*3\r\n$8\r\nREPLCONF\r\n$14\r\nlistening-port\r\n$4\r\n%s\r\n"
	REPLCONF_2 = "*3\r\n$8\r\nREPLCONF\r\n$4\r\ncapa\r\n$6\r\npsync2\r\n"
	PSYNC      = "*3\r\n$5\r\nPSYNC\r\n$1\r\n?\r\n$2\r\n-1\r\n"
	PING       = "*1\r\n$4\r\nping\r\n"
)

type RedisInterface interface {
	Handle_Conn(conn net.Conn)
	// Start_Server(port int)
}
type Redis struct {
	Comms      *chan resp.Value
	Store      *store.Db
	Repl_info  repl.ReplicationInfo
	SlavesConn []net.Conn
}

type Redis_Master_Server struct {
	*Redis
	IsSlave bool
}

type Redis_Slave_Server struct {
	*Redis
	IsSlave bool
}

func New_Redis_Master_Server() *Redis_Master_Server {
	comms := make(chan resp.Value, 1)
	return &Redis_Master_Server{
		Redis: &Redis{
			Store: store.NewDb(),
			Repl_info: repl.ReplicationInfo{
				Role:               "master",
				Master_replid:      util.Randomalphanumericgenerator(40),
				Master_repl_offset: 0,
			},
			Comms:      &comms,
			SlavesConn: []net.Conn{},
		},
		IsSlave: false,
	}
}

func New_Redis_Slave_Server() *Redis_Slave_Server {
	comms := make(chan resp.Value, 1)
	return &Redis_Slave_Server{
		Redis: &Redis{
			Store: store.NewDb(),
			Repl_info: repl.ReplicationInfo{
				Role:               "slave",
				Master_replid:      util.Randomalphanumericgenerator(40),
				Master_repl_offset: 0,
			},
			Comms: &comms,
		},
		IsSlave: true,
	}
}

func (sRedis *Redis_Slave_Server) Handshake_Slave_Master(conn net.Conn, masterHost string, masterPort string) {
	fmt.Println("In Handshake_Slave_Master")
	respReader := resp.NewRespHandler(conn)
	go func() {
		for {
			value, err := respReader.ParseAny()
			if err == io.EOF {
				continue
			}
			switch value.Typ {
			case resp.ArrayType:
				Reqargs := value.Array
				Comm := Reqargs[0].Bulk
				Comm_Args := Reqargs[1:]
				Metadata := &commands.MetaData{
					Db: sRedis.Store,
					Ri: sRedis.Repl_info,
				}
				commands.Handlers[strings.ToLower(Comm)](Metadata, Comm_Args)
			}
		}
	}()
	sRedis.Sync_Master(conn, masterPort)
}

func (sRedis *Redis_Slave_Server) Sync_Master(conn net.Conn, masterPort string) {
	fmt.Println("In Sync_Master ")
	buf := make([]byte, 1024)
	_, err := conn.Write([]byte(PING))
	_, err = conn.Read(buf[:])
	if err != nil {
		fmt.Println("Error:", err)
	}

	_, err = conn.Write([]byte(fmt.Sprintf(REPLCONF_1, masterPort)))
	_, err = conn.Read(buf[:])
	if err != nil {
		fmt.Println("Error Reading from connection")
	}

	_, err = conn.Write([]byte(REPLCONF_2))
	_, err = conn.Read(buf[:])
	if err != nil {
		fmt.Println("Error Reading from connection")
	}

	_, err = conn.Write([]byte(PSYNC))
	_, err = conn.Read(buf[:])
	if err != nil {
		fmt.Println("Error Reading from connection")
	}
}

func (redis *Redis) Handle_Comms(respWriter *resp.Writer, conn net.Conn) {
	fmt.Println("In Handle_Comms")
	// for val := range redis.Comms {
	// 	fmt.Println(val)
	// }
	for {
		fmt.Println("Waiting for comm")
		Value := <-*redis.Comms
		fmt.Println("Got for comm")
		// fmt.Println(Value)
		switch Value.Typ {
		case resp.ArrayType:
			fmt.Println("In ArrayType")
			Reqargs := Value.Array
			Comm := Reqargs[0].Bulk
			Comm_Args := Reqargs[1:]
			Metadata := &commands.MetaData{
				Db:     redis.Store,
				Ri:     redis.Repl_info,
				Comm:   *redis.Comms,
				Slaves: redis.SlavesConn,
				Client: conn,
			}
			respValue := commands.Handlers[strings.ToLower(Comm)](Metadata, Comm_Args)
			respWriter.Write(respValue)
			// if strings.ToLower(Comm) == "psync" {
			// 	// I know that it is a slave server and not a normal client and the Handshake is successful
			// 	respWriter.Write(commands.SendEmptyRDb(Metadata))
			// }
			// case resp.StringType:
			// 	switch strings.ToLower(Value.Str) {
			// 	case "ping":
			// 		respValue := resp.Value{
			// 			Typ: resp.StringType,
			// 			Str: "PONG",
			// 		}

			// 		respWriter.Write(respValue)
			// 	}
		}
	}

}

func Init_Resp(conn net.Conn) (*resp.Resp, *resp.Writer) {
	respHandler := resp.NewRespHandler(conn)
	respWriter := resp.NewRespWriter(conn)
	return respHandler, respWriter
}

func (sRedis *Redis_Slave_Server) Handle_Conn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connected from handleConn")
	fmt.Println("Handling conn from Redis Type")
	respHandler, respWriter := Init_Resp(conn)
	go sRedis.Handle_Comms(respWriter, conn)
	for {
		value, err := respHandler.ParseAny()
		if err == io.EOF {
			continue
		}
		fmt.Println(value)
		*sRedis.Comms <- value
	}
}

func (mRedis *Redis_Master_Server) Handle_Conn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connected from handleConn")
	fmt.Println("Handling conn from Redis_Master_Server type")
	respHandler, respWriter := Init_Resp(conn)
	go mRedis.Handle_Comms(respWriter, conn)
	for {
		value, err := respHandler.ParseAny()
		if err == io.EOF {
			continue
		}
		fmt.Println("Sending to routine")
		*mRedis.Comms <- value
		fmt.Printf("size of comms channel : %d", len(*mRedis.Comms))
		fmt.Println("Value sent to routine")
	}
}
func Init_Server(isSlave bool) (*Redis_Master_Server, *Redis_Slave_Server) {
	if isSlave {
		redisServer := New_Redis_Slave_Server()
		return nil, redisServer
	}
	redisServer := New_Redis_Master_Server()
	return redisServer, nil
}

// func (Redis *Redis) Start_Server(port int) {}
func (mRedis *Redis_Master_Server) Start_Server(port int) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d", port)
		os.Exit(1)
	}
	go mRedis.Handle_Slave_Conn()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go mRedis.Handle_Conn(conn)
	}
}
func (sRedis *Redis_Slave_Server) Start_Server(port int, conn net.Conn, masterHost string, masterPort string) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d", port)
		os.Exit(1)
	}
	go sRedis.Handshake_Slave_Master(conn, masterHost, masterPort)
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		go sRedis.Handle_Conn(conn)
	}
}

func (mRedis *Redis_Master_Server) Handle_Slave_Conn() {
	for {
		Value := <-*mRedis.Comms
		for _, slave := range mRedis.SlavesConn {
			slave.Write(Value.Bytes)
		}
	}
}
