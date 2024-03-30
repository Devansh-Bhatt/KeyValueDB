package commands

import (
	"fmt"
	"strconv"

	"github.com/codecrafters-io/redis-starter-go/resp"
	. "github.com/codecrafters-io/redis-starter-go/resp"
	. "github.com/codecrafters-io/redis-starter-go/store"
)

var Handlers = map[string]func(*Db, []Value) Value{
	"echo": Echo,
	"set":  Set,
	"get":  Get,
	"ping": Ping,
}

// const (
// 	set = "Set",
// 	echo = "Echo",
// 	get = "Get"

// )

func Ping(db *Db, args []Value) Value {
	return Value{
		Typ: StringType,
		Str: "PONG",
	}
}

func Echo(db *Db, args []Value) Value {
	return Value{
		Typ:  BulkStringType,
		Bulk: args[0].Bulk,
	}
}

func Set(db *Db, args []Value) Value {
	key := args[0].Bulk
	val := args[1].Bulk

	switch len(args) {
	case 2:
		db.Set(key, []byte(val), -1)

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
		db.Set(key, []byte(val), conv)
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

func Get(db *Db, args []Value) resp.Value {
	fmt.Println("reached Get")
	val, err := db.Get(args[0].Bulk)

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
