package commands

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/repl"
	"github.com/codecrafters-io/redis-starter-go/resp"
	. "github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type Client_Comm struct {
	Comm   resp.Value
	RW     *resp.Writer
	Client *net.Conn
}

type MetaData struct {
	Db     *store.Db
	Ri     *repl.ReplicationInfo
	Comm   chan Client_Comm
	Slaves *[]resp.Writer
	Client *net.Conn
	RW     *resp.Writer
}

var Handlers = map[string]func(*MetaData, []Value) Value{
	"echo":     Echo,
	"set":      Set,
	"get":      Get,
	"ping":     Ping,
	"info":     Info,
	"replconf": ReplConf,
	"psync":    Psync,
	"emptyrdb": SendEmptyRDb,
	"addslave": Add_Slave,
}

func Ping(Md *MetaData, args []Value) Value {
	// fmt.Println("In Ping")
	return Value{
		Typ: StringType,
		Str: "PONG",
	}
}

func Echo(Md *MetaData, args []Value) Value {
	return Value{
		Typ:  BulkStringType,
		Bulk: args[0].Bulk,
	}
}

func Set(Md *MetaData, args []Value) Value {
	// fmt.Println("Reached Set")
	key := args[0].Bulk
	val := args[1].Bulk

	switch len(args) {
	case 2:
		Md.Db.Set(key, []byte(val), -1)
		fmt.Printf("Set : %v %v\n", key, val)
		return Value{
			Typ: StringType,
			Str: "OK",
		}
	case 4:
		expiry := args[3].Bulk

		conv, err := strconv.ParseInt(expiry, 10, 64)

		if err != nil {
			return Value{
				Typ: ErrorType,
				Err: "Wrong arguments",
			}
		}
		Md.Db.Set(key, []byte(val), conv)
		// fmt.Println("Key with TTL has been set")
		return Value{
			Typ: StringType,
			Str: "OK",
		}
	default:
		// fmt.Println("Here ....... ")
		return Value{
			Typ: NULLType,
			// Err: "Could not set the Value",
		}
	}
}

func Get(Md *MetaData, args []Value) resp.Value {
	// fmt.Println("Reached Get")
	val, err := Md.Db.Get(args[0].Bulk)
	fmt.Printf("Get %v: %v \n", args[0].Bulk, string(val))
	if err != nil {
		return Value{
			Typ: NULLType,
		}
	}
	return Value{
		Typ:  BulkStringType,
		Bulk: string(val),
	}
}

func Info(Md *MetaData, args []Value) resp.Value {
	// fmt.Println("In Get Info")
	SubComm := strings.ToLower(args[0].Bulk)
	switch SubComm {
	case "replication":
		return Md.Ri.GetInfo()
	default:
		return Value{
			Typ: NULLType,
		}
	}

}

func ReplConf(Md *MetaData, args []Value) Value {
	// fmt.Println("In Repl_Conf")
	// fmt.Println(args[0].Bulk)
	switch args[0].Bulk {
	case "listening-port":
		Add_Slave := Value{
			Typ: ArrayType,
			Array: []Value{
				{
					Bulk: "AddSlave",
				},
			},
		}
		// fmt.Println("Putting in the channel the command to add the slave")
		go func() {
			Md.Comm <- Client_Comm{
				Comm:   Add_Slave,
				RW:     Md.RW,
				Client: Md.Client,
			}
		}()
		// fmt.Println("added slave")
		return Value{
			Typ: StringType,
			Str: "OK",
		}
	case "capa":
		return Value{
			Typ: StringType,
			Str: "OK",
		}
	default:
		return Value{
			Typ: ErrorType,
			Str: "Wrong Command",
		}

	}

}

func Add_Slave(Md *MetaData, args []Value) Value {
	*Md.Slaves = append(*Md.Slaves, *Md.RW)
	Md.Ri.Connected_slaves++
	// fmt.Println(Md.Slaves)
	return Value{
		Typ: StringType,
		Str: "OK",
	}
}

func Psync(Md *MetaData, args []Value) Value {
	// fmt.Println("In Psync")
	FullResyncComm := Value{
		Typ: ArrayType,
		Array: []Value{
			{
				Bulk: "EmptyRDB",
			},
		},
	}
	go func() {
		Md.Comm <- Client_Comm{
			Comm:   FullResyncComm,
			Client: Md.Client,
			RW:     Md.RW,
		}
	}()
	return Value{
		Typ: StringType,
		Str: fmt.Sprintf("FULLRESYNC %s 0", Md.Ri.Master_replid),
	}
}

func SendEmptyRDb(Md *MetaData, args []Value) Value {
	// fmt.Println("In Send Empty RDB")
	return Md.Ri.FullResync()
}
