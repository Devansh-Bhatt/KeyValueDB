package parser

import (
	"strings"
)

const (
	Integer      = ':'
	Strings      = '+'
	Bulk_Strings = '$'
	Arrays       = '*'
)

func MainParser(buf []byte) []byte {
	first := buf[0]
	response := make([]byte, 124)
	switch first {
	case Integer:
		response = parseInteger(buf)
	case Strings:
		response = parseString(buf)
	case Bulk_Strings:
		response = parseBulk(buf)
	case Arrays:
		response = parseArrays(buf)
	}
	return response

}

func parseArrays(buf []byte) []byte {
	// length := int(buf[1] - '0')

	commlength := int(buf[5] - '0')

	command := buf[8 : 8+commlength]

	if strings.ToLower(string(command)) == "echo" {
		return ParseEcho(buf[10+commlength:])
	} else if strings.ToLower(string(command)) == "set" {
		return ParseSet(buf[10+commlength:])
	} else if strings.ToLower(string(command)) == "get" {
		return ParseGet(buf[10+commlength:])
	} else {
		return []byte("+PONG\r\n")
	}
}

func parseBulk(buf []byte) []byte {
	return buf
}

func parseInteger(buf []byte) []byte {
	return buf
}

func parseString(buf []byte) []byte {
	return buf
}
