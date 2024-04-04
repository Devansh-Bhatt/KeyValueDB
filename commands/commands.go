package commands

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/redis"
	"github.com/codecrafters-io/redis-starter-go/resp"
	. "github.com/codecrafters-io/redis-starter-go/resp"
)

var Handlers = map[string]func(*redis.Redis, []Value) Value{
	"echo":     Echo,
	"set":      Set,
	"get":      Get,
	"ping":     Ping,
	"info":     Info,
	"replconf": ReplConf,
	"psync":    Psync,
}

// const (
// 	set = "Set",
// 	echo = "Echo",
// 	get = "Get"

// )

func Ping(redis *redis.Redis, args []Value) Value {
	return Value{
		Typ: StringType,
		Str: "PONG",
	}
}

func Echo(redis *redis.Redis, args []Value) Value {
	return Value{
		Typ:  BulkStringType,
		Bulk: args[0].Bulk,
	}
}

func Set(redis *redis.Redis, args []Value) Value {
	key := args[0].Bulk
	val := args[1].Bulk

	switch len(args) {
	case 2:
		redis.Store.Set(key, []byte(val), -1)

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
		redis.Store.Set(key, []byte(val), conv)
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

func Get(redis *redis.Redis, args []Value) resp.Value {
	fmt.Println("reached Get")
	val, err := redis.Store.Get(args[0].Bulk)

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

func Info(redis *redis.Redis, args []Value) resp.Value {
	SubComm := strings.ToLower(args[0].Bulk)
	switch SubComm {
	case "replication":
		return redis.GetInfo()
	default:
		return Value{
			Typ: NULLType,
		}
	}

}

func ReplConf(redis *redis.Redis, args []Value) Value {
	return Value{
		Typ: StringType,
		Str: "OK",
	}
}

func Psync(redis *redis.Redis, args []Value) Value {
	return Value{
		Typ: StringType,
		Str: fmt.Sprintf("FULLRESYNC %s 0", redis.Repl_info.Master_replid),
	}
}
