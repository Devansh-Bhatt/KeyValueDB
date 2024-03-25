package parser

import (
	"fmt"
)

var storage map[string]string = make(map[string]string)

func ParseSet(buf []byte) []byte {
	keylength := buf[1] - '0'

	valuelength := buf[7+keylength] - '0'

	key := buf[4 : 4+keylength]
	value := buf[9+keylength : 9+keylength+valuelength]

	storage[string(key)] = string(value)

	return []byte(fmt.Sprintf("+OK\r\n"))
}

func ParseGet(buf []byte) []byte {
	keylength := buf[1] - '0'

	key := buf[4 : 4+keylength]

	if value := storage[string(key)]; value == "" {
		return []byte("$-1\r\n")
	} else {
		return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
	}
}
