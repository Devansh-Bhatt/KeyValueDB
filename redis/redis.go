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
	Comms      chan commands.Client_Comm
	Store      *store.Db
	Repl_info  repl.ReplicationInfo
	SlavesConn *[]resp.Writer
	Clients    *[]net.Conn
}

type Redis_Master_Server struct {
	*Redis
	IsSlave     bool
	Slave_Comms chan commands.Client_Comm
}

type Redis_Slave_Server struct {
	*Redis
}

func New_Redis_Master_Server() *Redis_Master_Server {
	comms := make(chan commands.Client_Comm)
	slavecomms := make(chan commands.Client_Comm)
	return &Redis_Master_Server{
		Redis: &Redis{
			Store: store.NewDb(),
			Repl_info: repl.ReplicationInfo{
				Role:               "master",
				Master_replid:      util.Randomalphanumericgenerator(40),
				Master_repl_offset: 0,
				Connected_slaves:   0,
			},
			Comms:      comms,
			SlavesConn: &[]resp.Writer{},
			Clients:    &[]net.Conn{},
		},
		IsSlave:     false,
		Slave_Comms: slavecomms,
	}
}

func New_Redis_Slave_Server() *Redis_Slave_Server {
	comms := make(chan commands.Client_Comm)
	return &Redis_Slave_Server{
		Redis: &Redis{
			Store: store.NewDb(),
			Repl_info: repl.ReplicationInfo{
				Role:               "slave",
				Master_replid:      util.Randomalphanumericgenerator(40),
				Master_repl_offset: 0,
				Connected_slaves:   0,
			},
			Comms: comms,
		},
		// IsSlave: true,
	}
}

func (sRedis *Redis_Slave_Server) Handshake_Slave_Master(conn net.Conn, masterHost string, masterPort string) {
	fmt.Println("In Handshake_Slave_Master")
	respReader := resp.NewRespHandler(conn)
	// respWriter := resp.NewRespWriter(conn)
	sRedis.Sync_Master(conn, masterPort)
	// var buffer
	for {
		value, err := respReader.ParseAny()
		fmt.Println(value)
		// respWriter.Write(value)
		if err == io.EOF {
			continue
		}
		// sRedis.Comms <- commands.Client_Comm{
		// 	Comm: value,
		// 	RW:   respWriter,
		// }
		switch value.Typ {
		case resp.ArrayType:
			Reqargs := value.Array
			Comm := Reqargs[0].Bulk
			Comm_Args := Reqargs[1:]
			Metadata := &commands.MetaData{
				Db: sRedis.Store,
				Ri: &sRedis.Repl_info,
			}
			commands.Handlers[strings.ToLower(Comm)](Metadata, Comm_Args)
		case resp.RDBType:
			fmt.Println("here got an RDB")
			fmt.Println(value)
		}
	}

}

func (sRedis *Redis_Slave_Server) Sync_Master(conn net.Conn, masterPort string) {
	fmt.Println("In Sync_Master ")
	// buf := make([]byte, 1024)
	_, err := conn.Write([]byte(PING))
	// _, err = conn.Read(buf[:])
	fmt.Printf("Sync_Ping Done \n")
	if err != nil {
		fmt.Println("Error:", err)
	}

	_, err = conn.Write([]byte(fmt.Sprintf(REPLCONF_1, masterPort)))
	// _, err = conn.Read(buf[:])
	// fmt.Printf("%v", string(buf))
	_, err = fmt.Printf("Sync_Repl_Conf 1 Done \n")
	if err != nil {
		fmt.Println("Error Reading from connection")
	}

	_, err = conn.Write([]byte(REPLCONF_2))
	// _, err = conn.Read(buf[:])
	// fmt.Printf("%v", string(buf))
	_, err = fmt.Printf("Sync_Repl_Conf 2 Done \n")
	if err != nil {
		fmt.Println("Error Reading from connection")
	}

	_, err = conn.Write([]byte(PSYNC))
	// _, err = conn.Read(buf[:])
	// fmt.Printf("%v", string(buf))
	_, err = fmt.Printf("Sync_Pysnc Done \n")
	if err != nil {
		fmt.Println("Error Reading from connection")
	}
}

func (redis *Redis) Handle_Comms() {
	fmt.Println("In Handle_Comms")
	Metadata := &commands.MetaData{
		Db:     redis.Store,
		Ri:     &redis.Repl_info,
		Comm:   redis.Comms,
		Slaves: redis.SlavesConn,
	}
	for {
		select {
		case Value := <-redis.Comms:
			{
				Metadata.Client = Value.Client
				Metadata.RW = Value.RW
				fmt.Println(Value)
				switch Value.Comm.Typ {
				case resp.ArrayType:
					Reqargs := Value.Comm.Array
					Comm := Reqargs[0].Bulk
					fmt.Println(Comm)
					Comm_Args := Reqargs[1:]
					respValue := commands.Handlers[strings.ToLower(Comm)](Metadata, Comm_Args)
					fmt.Printf("The response : %v \n", respValue)
					err := Value.RW.Write(respValue)
					fmt.Println("The response has been sent")
					if err != nil {
						fmt.Println(err)
					}
				case resp.RDBType:
					fmt.Println("here in RDB")
					fmt.Println(Value.Comm.Bytes)

				case resp.BulkStringType:
					fmt.Println(Value.Comm)
				}
			}
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
	for {
		value, err := respHandler.ParseAny()
		if err == io.EOF {
			continue
		}
		fmt.Println(value)
		sRedis.Comms <- commands.Client_Comm{
			Comm: value,
			RW:   respWriter,
		}
	}
}

func (mRedis *Redis_Master_Server) Handle_Conn(conn net.Conn) {
	defer conn.Close()
	fmt.Println("Connected from handleConn")
	fmt.Println("Handling conn from Redis_Master_Server type")
	respHandler, respWriter := Init_Resp(conn)
	// go mRedis.Handle_Comms(respWriter, conn)
	for {
		value, err := respHandler.ParseAny()
		if err == io.EOF {
			continue
		}
		mRedis.Comms <- commands.Client_Comm{
			Comm:   value,
			RW:     respWriter,
			Client: &conn,
		}
		mRedis.Slave_Comms <- commands.Client_Comm{
			Comm:   value,
			RW:     respWriter,
			Client: &conn,
		}
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

func (mRedis *Redis_Master_Server) Start_Server(port int) {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		fmt.Printf("Failed to bind to port %d", port)
		os.Exit(1)
	}
	go mRedis.Handle_Comms()
	go mRedis.Handle_Slave_Conn()
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}
		*mRedis.Clients = append(*mRedis.Clients, conn)
		// ClientList := *mRedis.Clients
		// ClientList = append(ClientList, conn)
		// mRedis.Clients = &ClientList
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
	go sRedis.Handle_Comms()
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
		select {
		case Value := <-mRedis.Slave_Comms:
			{

				if Value.Comm.Array[0].Bulk == "set" {
					fmt.Println("Sending comm to slave")
					fmt.Printf("Size of slave list is : %v \n", len(*mRedis.SlavesConn))
					for _, slave := range *mRedis.SlavesConn {
						fmt.Printf("The value is : %v to %v \n", Value.Comm, slave)
						slave.Write(Value.Comm)

					}
				}
			}
		}
	}
}
