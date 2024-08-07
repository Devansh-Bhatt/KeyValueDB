package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/repl"
	"github.com/codecrafters-io/redis-starter-go/resp"
	. "github.com/codecrafters-io/redis-starter-go/resp"
	"github.com/codecrafters-io/redis-starter-go/store"
)

type MetaData struct {
	Db *store.Db
	Ri repl.ReplicationInfo
}

var Handlers = map[string]func(*MetaData, []Value) Value{
	"echo":     Echo,
	"set":      Set,
	"get":      Get,
	"ping":     Ping,
	"info":     Info,
	"replconf": ReplConf,
	"psync":    Psync,
}

func Ping(Md *MetaData, args []Value) Value {
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
	fmt.Println("reached set")
	key := args[0].Bulk
	val := args[1].Bulk

	switch len(args) {
	case 2:
		Md.Db.Set(key, []byte(val), -1)
		fmt.Println("Key has been set")
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
		fmt.Println("Key with TTL has been set")
		return Value{
			Typ: StringType,
			Str: "OK",
		}
	default:
		return Value{
			Typ: NULLType,
			// Err: "Could not set the Value",
		}
	}
}

func Get(Md *MetaData, args []Value) resp.Value {
	fmt.Println("reached Get")
	val, err := Md.Db.Get(args[0].Bulk)
	fmt.Println(string(val))
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
	return Value{
		Typ: StringType,
		Str: "OK",
	}
}

func Psync(Md *MetaData, args []Value) Value {
	return Value{
		Typ: StringType,
		Str: fmt.Sprintf("FULLRESYNC %s 0", Md.Ri.Master_replid),
	}
}

func SendEmptyRDb(Md *MetaData) Value {
	return Md.Ri.FullResync()
}
