package parser

import "fmt"

func ParseEcho(buf []byte) []byte {
	arglength := buf[1] - '0'
	msg := string(buf[4 : 4+arglength])
	// response := "$" + string(arglength) + "\r\n" + string(msg) + "\r\n"

	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", arglength, msg))
}
