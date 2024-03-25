package parser

func ParseEcho(buf []byte) []byte {
	arglength := buf[1] - '0'
	msg := buf[4 : 4+arglength]
	response := "$" + string(arglength) + "\r\n" + string(msg) + "\r\n"

	return []byte(response)
}
